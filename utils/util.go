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

func BreakFile(videoPath string, fileName string) bool {
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
		UploadToAppwrite(folder.Name(), db)
	}
}

func UploadToAppwrite(folderName string, db *sql.DB) {
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
			hashChecksum := sha1.New()
			hashChecksum.Write([]byte(fileComps[0]))
			fileId = fmt.Sprintf("%x", hashChecksum.Sum(nil))[:36]
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
		log.Println(response.StatusCode)
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
