package utils

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	if err := os.Mkdir(fmt.Sprintf("segments/%s", fileName), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-codec", "copy", "-map", "0", "-f", "segment", "-segment_time", "7", "-segment_format", "mpegts", "-segment_list", os.Getenv("ROOT_PATH")+"\\segments\\"+fileName+"\\"+fileName+".m3u8", "-segment_list_type", "m3u8", os.Getenv("ROOT_PATH")+"\\segments\\"+fileName+"\\"+fileName+"_"+"segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", output)
		log.Fatal(err)
		return false
	} else {
		return true
	}
}

func ResumeUploadIfAny(db *sql.DB) {
	folders, err := ioutil.ReadDir("segments")

	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		UploadToDeta(folder.Name(), db)
	}
}

func UploadToDeta(folderName string, db *sql.DB) {
	files, err := ioutil.ReadDir(fmt.Sprintf("segments/%s", folderName))

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

	log.Println("Now uploading chunks of " + folderName + " to Deta Drive...")
	var count int = -1
	for idx, file := range files {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("segments/%s/%s", folderName, file.Name()))

		if err != nil {
			log.Fatal(err)
		}

		postBody := bytes.NewBuffer(fileBytes)
		param := url.QueryEscape(file.Name())
		uploadRequestURL := "https://drive.deta.sh/v1/" + os.Getenv("PROJECT_ID") + "/video-streaming-server/files?name=" + param

		request, err := http.NewRequest("POST", uploadRequestURL, postBody)
		request.Header.Add("X-Api-Key", os.Getenv("PROJECT_KEY"))

		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != 201 {
			log.Fatal(response.StatusCode)
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
		upload_status=?,
		upload_end_time=?
		WHERE
		video_id=?;
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

