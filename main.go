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

func breakFile(videoPath string, fileName string) {
	// ffmpeg -y -i DearZindagi.mkv -codec copy -map 0 -f segment -segment_time 7 -segment_format mpegts -segment_list DearZindagi_index.m3u8 -segment_list_type m3u8 ./segment_no_%d.ts

	cmd := exec.Command("ffmpeg", "-y" , "-i" , videoPath, "-codec", "copy", "-map", "0","-f", "segment", "-segment_time", "10", "-segment_format", "mpegts", "-segment_list", "D:/segments/" + fileName + ".m3u8", "-segment_list_type", "m3u8", "D:/segments/" + fileName + "_" + "segment_no_%d.ts")

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
		fmt.Fprintf(w, "\nFile received completely!!")
		fmt.Println("Todne ka prayaas chalu hain....")
		breakFile(("./video/"+fileName), fileName)
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