package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	// "github.com/Chirag-And-Dheeraj/video-streaming-server/models"
)

func loadEnvVars() {
	fmt.Println("Setting environment variables...")
	
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

	fmt.Println("Environment variables set.")	
}

func breakFile(videoPath string, fileName string) {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	if err := os.Mkdir(fmt.Sprintf("segments/%s", fileName), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-y" , "-i" , videoPath, "-codec", "copy", "-map", "0","-f", "segment", "-segment_time", "10", "-segment_format", "mpegts", "-segment_list", "C:\\Users\\Dell\\Desktop\\Documents\\Projects\\Video-Streaming\\video-streaming-server\\segments\\" + fileName + "\\" + fileName + ".m3u8", "-segment_list_type", "m3u8", "C:\\Users\\Dell\\Desktop\\Documents\\Projects\\Video-Streaming\\video-streaming-server\\segments\\"  + fileName + "\\" + fileName + "_" + "segment_no_%d.ts")

	output, err := cmd.CombinedOutput()

	// err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", output)
		log.Fatal(err)
	} else {
		fmt.Printf("Tod di file, tukdo tukdo mein.")
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("Video upload endpoint hit...")
	fileName := r.Header.Get("file-name")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))
	fmt.Println("Name of the file received:", fileName)
	fmt.Println("Size of the file received:", fileSize)
	d, _ := ioutil.ReadAll(r.Body)
	tmpFile, _ := os.OpenFile("./video/"+fileName, os.O_APPEND|os.O_CREATE, 0644)
	tmpFile.Write(d)
	fmt.Fprintf(w, "Received chunk!")
	defer tmpFile.Close()

	fileInfo, _ := tmpFile.Stat()
	fmt.Println(fileInfo.Size())
	fmt.Println("Extra:", int64(fileSize) - int64(fileInfo.Size()))
	if fileInfo.Size() == int64(fileSize) {
		fmt.Fprintf(w, "\nFile received completely!!")
		fmt.Println("Todne ka prayaas chalu hain....")
		breakFile(("./video/"+fileName), fileName)

		// initializeUpload := fmt.Sprintf("https://drive.deta.sh/v1/c0unaxfn/video-streaming-server/uploads?name=%s", fileName)
		// log.Println(initializeUpload)

		// request, err := http.NewRequest("POST", initializeUpload, nil)
		// request.Header.Add("X-Api-Key", "c0unaxfn_dRDfc2XqobqBNYTeZcZf4uLZAZTrkoRb")
		// log.Println(request)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// client := &http.Client{}

    	// response, err := client.Do(request)

		// if err != nil {
		// 	log.Fatal(err)
		// }
    	// defer response.Body.Close()

		// jsonBody, err := ioutil.ReadAll(response.Body)

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// log.Printf(string(jsonBody))

		// body := make(map[string]interface{})
		// er := json.Unmarshal(jsonBody, &body)

		// if er != nil {
		// 	log.Fatal(err)
		// }

		// log.Println(body["upload_id"])

		// for i := 0; i < ; i++ {
		// 	chunkUpload := fmt.Sprintf("https://drive.deta.sh/v1/c0unaxfn/video-streaming-server/uploads/%s/parts?name=%s&part=%s", body["upload_id"], fileName, i)

		// 	req, err := http.NewRequest("POST", chunkUpload, )
		// 	req.Header.Add("X-Api-Key", "c0unaxfn_dRDfc2XqobqBNYTeZcZf4uLZAZTrkoRb")
		// }
		
		files, err := ioutil.ReadDir(fmt.Sprintf("segments/%s", fileName))
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fileBytes, err := ioutil.ReadFile(fmt.Sprintf("segments/%s/%s", fileName, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			postBody := bytes.NewBuffer(fileBytes)
			uploadChunk := fmt.Sprintf("https://drive.deta.sh/v1/{id}/video-streaming-server/files?name=%s/%s", fileName, file.Name())

			request, err := http.NewRequest("POST", uploadChunk, postBody)
			request.Header.Add("X-Api-Key", os.Getenv("PROJECT_KEY"))

			client := &http.Client{}

    		response, err := client.Do(request)

			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			jsonBody, err := ioutil.ReadAll(response.Body)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf(string(jsonBody))
		}
	}

	fmt.Println("---------------------------------------------------------------------")
}

var validPath = regexp.MustCompile("^/(upload)/([a-zA-Z0-9]+)$")

func setUpRoutes() {
	fmt.Println("Setting up routes...")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/video", videoHandler)
	fmt.Println("Routes set.")
}

func initServer() {
	fmt.Println("Initializing server...")
	loadEnvVars()
	setUpRoutes()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := "./client/index.html"
	http.ServeFile(w, r, p)
}

func main() {
	initServer()
	fmt.Println("Server is running on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}