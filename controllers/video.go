package controllers

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"video-streaming-server/utils"
	. "video-streaming-server/structs"
)

func closeVideoFile(tmpFile *os.File) {
	err := tmpFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(tmpFile.Name())

	if err != nil {
		log.Fatal(err)
	}
}

func UploadVideo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fileName := r.Header.Get("file-name")
	isFirstChunk := r.Header.Get("first-chunk")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))

	hashChecksum := sha1.New()

	hashChecksum.Write([]byte(fileName))

	fileNameHash := fmt.Sprintf("%x", hashChecksum.Sum(nil))

	serverFileName := fileNameHash + ".mp4"

	if isFirstChunk == "true" {
		title := r.Header.Get("title")
		description := r.Header.Get("description")
		log.Println("Started receiving chunks for: " + fileName)
		log.Println("Size of the file received:", fileSize)
		log.Println("Title: ", title)
		log.Println("Description: ", description)
		log.Println("FileName hash: " + fileNameHash)
		log.Println("Creating a database record...")

		insertStatement, err := db.Prepare(`
			INSERT INTO 
				videos
					(
						video_id,
						title,
						description,
						upload_initiate_time,
						upload_status
					) 
				VALUES 
					(?,?,?,?,?);
		`)

		if err != nil {
			log.Fatal(err)
		}

		_, err = insertStatement.Exec(fileNameHash, title, description, time.Now(), 0)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Print("Database record created.")
		}
	}
	d, _ := ioutil.ReadAll(r.Body)

	tmpFile, _ := os.OpenFile("./video/"+serverFileName, os.O_APPEND|os.O_CREATE, 0644)
	tmpFile.Write(d)

	fileInfo, _ := tmpFile.Stat()

	if fileInfo.Size() == int64(fileSize) {
		defer closeVideoFile(tmpFile)
		fmt.Fprintf(w, "\nFile received completely!!")
		log.Println("Received all chunks for: " + serverFileName)
		log.Println("Breaking the video into .ts files.")

		breakResult := utils.BreakFile(("./video/" + serverFileName), strings.Split(serverFileName, ".")[0])

		if breakResult {
			log.Println("Successfully broken " + fileName + " into .ts files.")
		} else {
			log.Println("Error breaking " + fileName + " into .ts files.")
		}
		utils.UploadToDeta(strings.Split(serverFileName, ".")[0])

		log.Println("Successfully uploaded chunks of", fileName, "to Deta Drive")
		log.Println("Updating upload status in database record...")
		updateStatement, err := db.Prepare(`
		UPDATE
			videos 
		SET 
			upload_status=?,
			upload_end_time=?
			WHERE
			video_id=?;
		`)

		if err != nil {
			log.Fatal(err)
		}

		_, err = updateStatement.Exec(1, time.Now(), fileNameHash)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Database record updated.")
			log.Println("Finished uploading", fileName, " :)")
		}
	}
}

func GetVideos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Querying the database now for a list of videos...")
	rows, err := db.Query(`
		SELECT 
			video_id,
			title,
			description
		FROM
			videos;
	`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	log.Println("Query executed.")

	records := make([]Video, 0)

	for rows.Next() {
		var id string
		var title string
		var description string

		err := rows.Scan(&id, &title, &description)

		if err != nil {
			log.Fatal(err)
		}

		record := Video{
			ID:          id,
			Title:       title,
			Description: description,
		}

		records = append(records, record)
	}

	recordsJSON, err := json.Marshal(records)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintf(w, string(recordsJSON))
	}
}

func GetManifestFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	video_id := strings.Split(r.URL.Path[1:], "/")[1]

	log.Println("Video ID: " + video_id)

	getManifestFile := "https://drive.deta.sh/v1/" + os.Getenv("PROJECT_ID") + "/video-streaming-server/files/download?name=" + video_id + ".m3u8"

	request, err := http.NewRequest("GET", getManifestFile, nil)
	request.Header.Add("X-Api-Key", os.Getenv("PROJECT_KEY"))

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.Write(bodyBytes)
}

func GetTSFiles(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoChunkFileName := strings.Split(r.URL.Path[1:], "/")[2]

	log.Println("Video chunk requested: " + videoChunkFileName)

	getSegmentFile := "https://drive.deta.sh/v1/" + os.Getenv("PROJECT_ID") + "/video-streaming-server/files/download?name=" + videoChunkFileName

	request, err := http.NewRequest("GET", getSegmentFile, nil)
	request.Header.Add("X-Api-Key", os.Getenv("PROJECT_KEY"))

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "video/MP2T")
	w.Write(bodyBytes)
}