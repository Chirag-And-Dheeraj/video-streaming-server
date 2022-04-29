package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"video-streaming-server/database"
)

func loadEnvVars() {
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

func breakFile (videoPath string, fileName string) bool {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	if err := os.Mkdir(fmt.Sprintf("segments/%s", fileName), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-y" , "-i" , videoPath, "-codec", "copy", "-map", "0","-f", "segment", "-segment_time", "10", "-segment_format", "mpegts", "-segment_list", "D:\\ideas\\video-streaming-server\\segments\\" + fileName + "\\" + fileName + ".m3u8", "-segment_list_type", "m3u8", "D:\\ideas\\video-streaming-server\\segments\\"  + fileName + "\\" + fileName + "_" + "segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	// err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", output)
		log.Fatal(err)
		return false
	} else {
		return true
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// log.Println(r.Method)
	// log.Println("Video upload endpoint hit...")
	fileName := r.Header.Get("file-name")
	isFirstChunk := r.Header.Get("first-chunk")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))

	if isFirstChunk == "true" {
		title := r.Header.Get("title")
		description := r.Header.Get("description")
		log.Println("Started receiving chunks for: " + fileName)
		log.Println("Size of the file received:", fileSize)
		log.Println("Title: ", title)
		log.Println("Description: ", description)
		log.Println("Creating a database record...")

		insertStatement, err := db.Prepare(`INSERT INTO videos
		(
			file_name,
			title,
			description,
			upload_initiate_time,
			upload_status
		) VALUES (?,?,?,?,?);`)

		if err != nil {
			log.Fatal(err)
		}

		result, err := insertStatement.Exec(fileName, title, description, time.Now(),0)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(result)
			log.Print("Database record created.")
		}

	}
	
	d, _ := ioutil.ReadAll(r.Body)
	tmpFile, _ := os.OpenFile("./video/"+fileName, os.O_APPEND|os.O_CREATE, 0644)
	tmpFile.Write(d)
	
	// fmt.Fprintf(w, "Received chunk!")
	defer tmpFile.Close()

	fileInfo, _ := tmpFile.Stat()
	
	// log.Println(fileInfo.Size())
	// log.Println("Extra:", int64(fileSize) - int64(fileInfo.Size()))
	
	if fileInfo.Size() == int64(fileSize) {
		fmt.Fprintf(w, "\nFile received completely!!")
		log.Println("Received all chunks for: " + fileName)
		log.Println("Breaking the video into .ts files.")

		breakResult := breakFile(("./video/"+fileName), fileName)

		if breakResult {
			log.Println("Successfully broken " + fileName + " into .ts files.")
		} else {
			log.Println("Error breaking " + fileName + " into .ts files.")
		}
		
		files, err := ioutil.ReadDir(fmt.Sprintf("segments/%s", fileName))
		
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Now uploading chunks of " + fileName + " to Deta Drive...")

		for _ , file := range files {
			fileBytes, err := ioutil.ReadFile(fmt.Sprintf("segments/%s/%s", fileName, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			postBody := bytes.NewBuffer(fileBytes)
			uploadChunk := fmt.Sprintf("https://drive.deta.sh/v1/"+ os.Getenv("PROJECT_ID") + "/video-streaming-server/files?name=%s/%s", fileName, file.Name())

			request, err := http.NewRequest("POST", uploadChunk, postBody)
			request.Header.Add("X-Api-Key", os.Getenv("PROJECT_KEY"))

			client := &http.Client{}

    		response, err := client.Do(request)

			if err != nil {
				log.Fatal(err)
			}
			// log.Println("Chunk number", i, "uploaded successfully.")
			defer response.Body.Close()
		}
		log.Println("Successfully uploaded chunks of", fileName, "to Deta Drive")
		log.Println("Updating upload status in database record...")
		updateStatement, err := db.Prepare(`
		UPDATE
			videos 
		SET 
			upload_status=?,
			upload_end_time=?
		WHERE
			file_name=?;
		`)

		if err != nil {
			log.Fatal(err)
		}

		result, err := updateStatement.Exec(1, time.Now(), fileName)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(result)
			log.Print("Database record updated.")
		}

	}
	// log.Println("---------------------------------------------------------------------")
}

var validPath = regexp.MustCompile("^/(upload)/([a-zA-Z0-9]+)$")

func setUpRoutes(db *sql.DB) {
	log.Println("Setting up routes...")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/video", func (w http.ResponseWriter, r *http.Request) {
		videoHandler(w, r, db)
	})
	log.Println("Routes set.")
}

func initServer() {
	log.Println("Initializing server...")
	loadEnvVars()
	db := database.Connect()
	setUpRoutes(db)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := "./client/index.html"
	http.ServeFile(w, r, p)
}

func main() {
	initServer()
	log.Println("Server is running on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}