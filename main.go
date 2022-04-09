package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	// "github.com/Chirag-And-Dheeraj/video-streaming-server/models"
)

func videoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Video upload endpoint hit...")
	fileName := r.Header.Get("file-name")
	fmt.Println("Name of the file received: " + fileName)	
	d, _ := ioutil.ReadAll(r.Body)
	tmpFile, _ := os.OpenFile("./video/"+fileName, os.O_APPEND|os.O_CREATE, 0644)
	defer tmpFile.Close()
	tmpFile.Write(d)
	fmt.Fprintf(w, "Received chunk!")
	w.WriteHeader(200)
	return
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
	fmt.Println("Server is running on port", 8000)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}