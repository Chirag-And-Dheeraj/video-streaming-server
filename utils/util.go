package utils

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func LoadEnvVars() {
	log.Println("Setting environment variables...")

	envFile, err := os.Open(".env")

	if err != nil {
		log.Fatal(err)
	}

	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)

	for scanner.Scan() {
		lineSplit := strings.Split(scanner.Text(), "=")
		os.Setenv(lineSplit[0], lineSplit[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Environment variables set.")
}

func breakFile(videoPath string, fileName string) bool {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	log.Println("Inside BreakFile function.")

	if err := os.Mkdir(fmt.Sprintf("segments/%s", fileName), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	log.Println("Created directory inside segments folder.")

	log.Println("Video path: " + videoPath)

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-codec", "copy", "-map", "0", "-f", "segment", "-segment_time", "7", "-segment_format", "mpegts", "-segment_list", os.Getenv("ROOT_PATH")+"/segments/"+fileName+"/"+fileName+".m3u8", "-segment_list_type", "m3u8", os.Getenv("ROOT_PATH")+"/segments/"+fileName+"/"+fileName+"_"+"segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("%s\n", output)
		log.Fatal(err)
		return false
	} else {
		return true
	}
}

func ResumeUploadIfAny(db *sql.DB) {
	folders, err := os.ReadDir("segments")

	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		uploadToAppwrite(folder.Name(), db)
	}
}

func uploadToAppwrite(folderName string, db *sql.DB) {
	files, err := os.ReadDir(fmt.Sprintf("segments/%s", folderName))

	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		err = os.Remove("segments/" + folderName)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Println("Now uploading chunks of " + folderName + " to Appwrite Storage...")
	var count int = -1
	for idx, file := range files {
		fileToUpload, err := os.ReadFile(fmt.Sprintf("segments/%s/%s", folderName, file.Name()))

		if err != nil {
			log.Fatal(err)
		}

		uploadRequestURL := "https://cloud.appwrite.io/v1/storage/buckets/" + os.Getenv("BUCKET_ID") + "/files"

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		fileId := "nil"
		fileComps := strings.Split(file.Name(), ".")
		if fileComps[1] == "m3u8" {
			fileId = fileComps[0]
		} else {
			fileId = GetFileId(fileComps[0])
		}

		err = writer.WriteField("fileId", fileId)

		if err != nil {
			log.Fatal(err)
		}

		part, err := writer.CreateFormFile("file", file.Name())

		if err != nil {
			log.Fatal(err)
		}

		_, err = part.Write(fileToUpload)

		if err != nil {
			log.Fatal(err)
		}

		err = writer.Close()

		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest("POST", uploadRequestURL, &requestBody)
		if err != nil {
			log.Printf("Error creating request")
			log.Fatal(err)
		}

		request.Header.Set("Content-Type", writer.FormDataContentType())
		request.Header.Set("X-Appwrite-Response-Format", os.Getenv("APPWRITE_RESPONSE_FORMAT"))
		request.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
		request.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_KEY"))

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != 201 {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal(string(body))
		} else {
			count = idx
			err = os.Remove("segments/" + folderName + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Updating upload status in database record...")
	updateStatement, err := db.Prepare(`
	UPDATE
		videos
	SET
		upload_status=$1,
		upload_end_time=$2
		WHERE
		video_id=$3;
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = updateStatement.Exec(1, time.Now(), folderName)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database record updated.")
		log.Println("Finished uploading", folderName, " :)")
	}

	if count == len(files)-1 {
		err = os.Remove("segments/" + folderName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

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

func PostUploadProcessFile(serverFileName string, fileName string, tmpFile *os.File, db *sql.DB) {
	log.Println("Received all chunks for: " + serverFileName)
	log.Println("Breaking the video into .ts files.")
	breakResult := breakFile(("./video/" + serverFileName), fileName)
	if breakResult {
		log.Println("Successfully broken " + fileName + " into .ts files.")
		log.Println("Deleting the original file from server.")
		closeVideoFile(tmpFile)
		uploadToAppwrite(fileName, db)
		log.Println("Successfully uploaded chunks of", fileName, "to Appwrite Storage")
	} else {
		log.Println("Error breaking " + fileName + " into .ts files.")
	}
}

func GetManifestFile(w http.ResponseWriter, videoId string) ([]byte, error) {
	log.Println("Video ID: " + videoId)

	getManifestFile := "https://cloud.appwrite.io/v1/storage/buckets/" + os.Getenv("BUCKET_ID") + "/files/" + videoId + "/view"

	request, err := http.NewRequest("GET", getManifestFile, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Appwrite-Response-Format", os.Getenv("APPWRITE_RESPONSE_FORMAT"))
	request.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
	request.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_KEY"))

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 404 {
		http.Error(w, "Video record not found", http.StatusNotFound)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func GetFileId(fileName string) string {
	hashChecksum := sha1.New()
	hashChecksum.Write([]byte(fileName))
	fileId := fmt.Sprintf("%x", hashChecksum.Sum(nil))[:36]

	return fileId
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, db *sql.DB, videoId string) {
	fileBytes, err := GetManifestFile(w, videoId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error Deleting Video", http.StatusInternalServerError)
		return
	}

	file := string(fileBytes)
	lines := strings.Split(file, "\n")

	deleteUrl := "https://cloud.appwrite.io/v1/storage/buckets/" + os.Getenv("BUCKET_ID") + "/files/"

	for i := 0; i < len(lines); i++ {
		if strings.HasSuffix(lines[i], ".ts") {
			fileName := strings.Split(lines[i], ".")[0]
			fileId := GetFileId(fileName)

			request, err := http.NewRequest("DELETE", deleteUrl + fileId, nil)

			if err != nil {
				log.Println(err)
				http.Error(w, "Error Deleting Chunk", http.StatusInternalServerError)
				return
			}

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Appwrite-Response-Format", os.Getenv("APPWRITE_RESPONSE_FORMAT"))
			request.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
			request.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_KEY"))

			client := &http.Client{}

			response, err := client.Do(request)
			if err != nil {
				log.Println(err)
				http.Error(w, "Error Deleting Chunk", http.StatusInternalServerError)
				return
			}
			defer response.Body.Close()
		}
	}

	log.Println("Deleted all .ts files...")

	request, err := http.NewRequest("DELETE", deleteUrl + videoId, nil)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error Deleting Video", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Appwrite-Response-Format", os.Getenv("APPWRITE_RESPONSE_FORMAT"))
	request.Header.Set("X-Appwrite-Project", os.Getenv("APPWRITE_PROJECT_ID"))
	request.Header.Set("X-Appwrite-Key", os.Getenv("APPWRITE_KEY"))

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error Deleting Video", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	log.Println("Deleted .m3u8 file...")

	query, err := db.Prepare(`DELETE FROM videos WHERE video_id=$1`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(videoId)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting record", http.StatusInternalServerError)
		return
	}

	log.Println("Deleted database record...")
}
