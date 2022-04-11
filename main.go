package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	// "github.com/Chirag-And-Dheeraj/video-streaming-server/models"
)

func breakFile(fileName string) {
	// ffmpeg -i input.mp4 -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls index.m3u8
	// exec.Command("cd video").Run()
	// exec.Command("dir").Run()

	cmdString := fmt.Sprintf("ffmpeg -i %s -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls %s.m3u8", "video/" + fileName, "video/" + fileName)
	err := exec.Command(cmdString).Run()

	if err != nil {
		panic(err)
	}
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Video upload endpoint hit...")
	fileName := r.Header.Get("file-name")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))
	fmt.Println("Name of the file received:", fileName)	
	fmt.Println("Size of the file received:", fileSize)
	d, _ := ioutil.ReadAll(r.Body)
	tmpFile, _ := os.OpenFile("./video/"+fileName, os.O_APPEND|os.O_CREATE, 0644)
	tmpFile.Write(d)
	fmt.Fprintf(w, "Received chunk!")

	fileInfo, _ := tmpFile.Stat()
	fmt.Println(fileInfo.Size())
	fmt.Println("Extra:", int64(fileSize) - int64(fileInfo.Size()))
	if fileInfo.Size() == int64(fileSize) {
		// breakFile(fileName)
		fmt.Fprintf(w, "\nFile received completely!!")
	}
	fmt.Println("---------------------------------------------------------------------")
}

var validPath = regexp.MustCompile("^/(upload)/([a-zA-Z0-9]+)$")

func setUpRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/video", videoHandler)
}

func initServer() {
	fmt.Println("Initializing server...")
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