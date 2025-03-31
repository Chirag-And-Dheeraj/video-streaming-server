package controllers

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"video-streaming-server/config"
	. "video-streaming-server/types"
	"video-streaming-server/utils"
)

// @desc Create new video resource
// @route POST /video
func UploadVideo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fileName := r.Header.Get("file-name")
	isFirstChunk := r.Header.Get("first-chunk")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))
	user, err := utils.GetUserFromRequest(r)

	if err != nil {
		log.Println(err)
		http.Error(w, "User not logged in.", http.StatusUnauthorized)
		return
	}
	sizeLimit, _ := strconv.Atoi(config.AppConfig.FileSizeLimit)

	if fileSize > sizeLimit {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File size greater than 15 MB is not acceptable"))
		return
	}

	serverFileName := fileName + ".mp4"

	if isFirstChunk == "true" {
		title := r.Header.Get("title")
		description := r.Header.Get("description")
		log.Println("Started receiving chunks for: " + fileName)
		log.Println("Size of the file received:", fileSize)
		log.Println("Title: ", title)
		log.Println("Description: ", description)
		log.Println("Creating a database record...")
		log.Println("User is: ", user.ID)

		insertStatement, err := db.Prepare(`
			INSERT INTO 
				videos
					(
						video_id,
						title,
						description,
						upload_initiate_time,
						upload_status,
						delete_flag,
						user_id
					) 
				VALUES 
					($1,$2,$3,$4,$5,$6,$7)
		`)

		if err != nil {
			log.Fatal(err)
		}

		_, err = insertStatement.Exec(fileName, title, description, time.Now(), 0, 0, user.ID)

		if err != nil {
			log.Println(err)
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		} else {
			log.Println("Database record created.")
		}
	}

	d, _ := io.ReadAll(r.Body)

	var tmpFile *os.File

	if isFirstChunk == "true" {
		tmpFile, err = os.Create("./video/" + serverFileName)
		if err != nil {
			log.Println("Error creating a temp file on the server:", err)
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}
	} else {
		tmpFile, err = os.OpenFile("./video/"+serverFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Println("Error opening the temp file on the server for appending chunks:", err)
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}
	}

	_, err = tmpFile.Write(d)

	if err != nil {
		log.Println("Error appending chunks to file:", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	fileInfo, err := tmpFile.Stat()

	if err != nil {
		log.Println("Error getting file info:", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	if fileInfo.Size() == int64(fileSize) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Video received completely and is now being processed."))
		go utils.PostUploadProcessFile(serverFileName, fileName, tmpFile, db)

	} else {
		w.WriteHeader(http.StatusPartialContent)
		w.Write([]byte("Receiving chunks of the video."))
	}
}

// @desc Get All Videos
// @route GET /video
func GetVideos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Querying the database now for a list of videos...")
	user, err := utils.GetUserFromRequest(r)

	if err != nil {
		log.Println(err)
		http.Error(w, "User not logged in.", http.StatusUnauthorized)
		return
	}

	getUserVideosQuery, err := db.Prepare(`
		SELECT
			video_id,
			title,
			description,
			thumbnail
		FROM
			videos
		WHERE
			upload_status=1
		AND 
			delete_flag=0
		AND
			user_id=$1
		ORDER BY
			upload_initiate_time DESC;
	`)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := getUserVideosQuery.Query(user.ID)

	if err != nil {
		log.Println("Error running select query for all video records.")
		http.Error(w, "Error retrieving videos", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	log.Println("Query executed.")

	records := make([]ListVideosResponseItem, 0)

	for rows.Next() {
		var id string
		var title string
		var description string
		var thumbnail sql.NullString

		err := rows.Scan(&id, &title, &description, &thumbnail)

		if err != nil {
			log.Println("Error scanning rows")
			log.Println(err)
			http.Error(w, "Error retreiving records", http.StatusInternalServerError)
			return
		}

		thumbValue := "../static/logo/android-chrome-192x192.png"
		if thumbnail.Valid && thumbnail.String != "" {
			thumbValue = thumbnail.String
		}

		record := ListVideosResponseItem{
			ID:          id,
			Title:       title,
			Description: description,
			Thumbnail:   thumbValue,
		}

		records = append(records, record)
	}

	recordsJSON, err := json.Marshal(records)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error retreiving records", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(recordsJSON))
	}
}

// @desc Get a Video
// @route GET /video/[id]
func GetVideo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]
	log.Println("Details of " + videoId + " requested.")

	user, err := utils.GetUserFromRequest(r)

	if err != nil {
		log.Println(err)
		http.Error(w, "User not logged in.", http.StatusUnauthorized)
		return
	}

	detailsQuery, err := db.Prepare(`
		SELECT
			title, description
		FROM
			videos
		WHERE
			video_id=$1
		AND
			user_id=$2;
	`)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error retrieving record", http.StatusInternalServerError)
	}

	defer detailsQuery.Close()

	var title, description string
	err = detailsQuery.QueryRow(videoId, user.ID).Scan(&title, &description)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "Video record not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Error retrieving record", http.StatusInternalServerError)
		return
	}

	videoDetails := &Video{
		ID:          videoId,
		Title:       title,
		Description: description,
	}
	videoDetailsJSON, err := json.Marshal(videoDetails)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error retrieving video", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(videoDetailsJSON))
	}
}

// @desc Get Manifest File
// @route GET /video/[id]/stream
func ManifestFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := strings.Split(r.URL.Path[1:], "/")[1]

	file, err := utils.GetManifestFile(w, videoId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error retrieving video", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/x-mpegURL")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	}
}

// @desc Get TS File
// @route GET /video/[id]/stream/[id].ts
func TSFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoName := strings.Split(r.URL.Path[1:], "/")[3]

	videoComps := strings.Split(videoName, ".")
	hashChecksum := sha1.New()
	hashChecksum.Write([]byte(videoComps[0]))
	fileId := fmt.Sprintf("%x", hashChecksum.Sum(nil))[:36]

	log.Println("Video chunk requested: " + fileId)

	getSegmentFile := "https://cloud.appwrite.io/v1/storage/buckets/" + config.AppConfig.AppwriteBucketID + "/files/" + fileId + "/view"

	request, err := http.NewRequest("GET", getSegmentFile, nil)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error Streaming Video", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Appwrite-Response-Format", config.AppConfig.AppwriteResponseFormat)
	request.Header.Set("X-Appwrite-Project", config.AppConfig.AppwriteProjectID)
	request.Header.Set("X-Appwrite-Key", config.AppConfig.AppwriteKey)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error Streaming Video", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		http.Error(w, "Video chunk not found", http.StatusNotFound)
		return
	}

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error Streaming Video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/MP2T")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

// @desc Update Video Details
// @route UPDATE
func UpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]

	if videoId == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	var reqBody UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	query, err := db.Prepare(`
		UPDATE
			videos
		SET
			title=$1,
			description=$2
		WHERE
			video_id=$3;
	`)

	if err != nil {
		log.Fatal(err)
	}

	result, err := query.Exec(reqBody.Title, reqBody.Description, videoId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error Streaming Video", http.StatusInternalServerError)
		return
	} else {
		rows, _ := result.RowsAffected()
		if  rows == 0 {
			http.Error(w, "No record found with given ID", http.StatusNotFound)
			return
		}
	}

	log.Println("Database record updated.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @desc Delete the video
// @route DELETE /video/[id]
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]

	log.Println("Updating delete_flag in database record...")
	updateStatement, err := db.Prepare(`
		UPDATE
			videos
		SET
			delete_flag=$1
			WHERE
			video_id=$2;
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = updateStatement.Exec(1, videoId)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database record updated.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	go utils.DeleteVideo(w, r, db, videoId)
}
