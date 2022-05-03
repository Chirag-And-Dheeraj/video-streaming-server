package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"video-streaming-server/database"
)

type video struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	File_name   string `json:"file_name"`
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

func breakFile(videoPath string, fileName string) bool {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	if err := os.Mkdir(fmt.Sprintf("segments/%s", fileName), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-codec", "copy", "-map", "0", "-f", "segment", "-segment_time", "10", "-segment_format", "mpegts", "-segment_list", os.Getenv("ROOT_PATH")+"\\segments\\"+fileName+"\\"+fileName+".m3u8", "-segment_list_type", "m3u8", os.Getenv("ROOT_PATH")+"\\segments\\"+fileName+"\\"+fileName+"_"+"segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", output)
		log.Fatal(err)
		return false
	} else {
		return true
	}
}

func uploadToDeta(fileName string) {
	files, err := ioutil.ReadDir(fmt.Sprintf("segments/%s", fileName))

	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		err = os.Remove("segments/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Println("Now uploading chunks of " + fileName + " to Deta Drive...")
	var count int = -1
	for idx, file := range files {
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("segments/%s/%s", fileName, file.Name()))

		if err != nil {
			log.Fatal(err)
		}

		postBody := bytes.NewBuffer(fileBytes)
		param := url.QueryEscape(fileName + "/" + file.Name())
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
			err = os.Remove("segments/" + fileName + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if count == len(files)-1 {
		err = os.Remove("segments/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
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

			_, err = insertStatement.Exec(fileName, title, description, time.Now(), 0)

			if err != nil {
				log.Fatal(err)
			} else {
				log.Print("Database record created.")
			}

		}

		d, _ := ioutil.ReadAll(r.Body)
		tmpFile, _ := os.OpenFile("./video/"+fileName, os.O_APPEND|os.O_CREATE, 0644)
		tmpFile.Write(d)
		defer tmpFile.Close()

		fileInfo, _ := tmpFile.Stat()

		if fileInfo.Size() == int64(fileSize) {
			fmt.Fprintf(w, "\nFile received completely!!")
			log.Println("Received all chunks for: " + fileName)
			log.Println("Breaking the video into .ts files.")

			breakResult := breakFile(("./video/" + fileName), fileName)

			if breakResult {
				log.Println("Successfully broken " + fileName + " into .ts files.")
			} else {
				log.Println("Error breaking " + fileName + " into .ts files.")
			}

			uploadToDeta(fileName)

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

			_, err = updateStatement.Exec(1, time.Now(), fileName)

			if err != nil {
				log.Fatal(err)
			} else {
				log.Print("Database record updated.")
				log.Println("Finished uploading", fileName, " :)")
			}
		}
	} else if r.Method == "GET" && len(r.URL.Query()) == 0 {
		log.Println("Get request on the video endpoint :)")
		log.Println("Querying the database now for a list of videos...")
		rows, err := db.Query(`
			SELECT 
				video_id,
				title,
				description,
				file_name
			FROM
				videos;
		`)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		log.Println("Query executed.")
		log.Println("Now printing results...")

		records := make([]video, 0)

		for rows.Next() {
			var id int
			var title string
			var description string
			var file_name string
			err := rows.Scan(&id, &title, &description, &file_name)

			if err != nil {
				log.Fatal(err)
			}

			record := video{
				ID:          id,
				Title:       title,
				Description: description,
				File_name:   file_name,
			}

			records = append(records, record)
		}
		
		recordsJSON, err := json.Marshal(records)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Records in JSON", string(recordsJSON))
			fmt.Fprintf(w, string(recordsJSON))
		}
	} else if r.Method == "GET" && len(r.URL.Query()) != 0 {
		log.Println("Bring me videoooooo....")
		
		params := r.URL.Query()

		if params.Has("v") {
			v := params.Get("v")

			stmt, err := db.Prepare(`
				SELECT 
					manifest_url 
				FROM 
					videos 
				WHERE 
					video_id=?
			`)

			var manifest_url string
			err = stmt.QueryRow(v).Scan(&manifest_url)

			if err != nil {
				log.Fatal(err)
			}

			log.Println(manifest_url)
		}
	}
}

func resumeUploadIfAny() {
	folders, err := ioutil.ReadDir("segments")

	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		uploadToDeta(folder.Name())
	}
}

var validPath = regexp.MustCompile("^/(upload)/([a-zA-Z0-9]+)$")

func homePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		return
	} else if r.Method == "GET" {
		log.Println("Get request on the home page endpoint :)")
		p := "./client/index.html"
		http.ServeFile(w, r, p)
	}
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("Get request on the upload page endpoint :)")
		p := "./client/upload.html"
		http.ServeFile(w, r, p)
	}
}

func viewPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("Get request on the view videos endpoint :)")
		p := "./client/view.html"
		http.ServeFile(w, r, p)
	}
}

func setUpRoutes(db *sql.DB) {
	log.Println("Setting up routes...")
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/upload", uploadPageHandler)
	http.HandleFunc("/view", viewPageHandler)
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		videoHandler(w, r, db)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Routes set.")
}

func initServer() {
	log.Println("Initializing server...")
	loadEnvVars()
	db := database.Connect()
	setUpRoutes(db)
}

func main() {
	initServer()
	log.Println("Server is running on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
	resumeUploadIfAny()
}
