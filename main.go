package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"video-streaming-server/controllers"
	"video-streaming-server/database"
	"video-streaming-server/utils"
	. "video-streaming-server/structs"
)

func videoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		controllers.UploadVideo(w, r, db)
	} else if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		log.Print("Manifest Request:")
		log.Println(regexp.MatchString("^/video/[a-zA-B0-9]+/?$", r.URL.Path))
		log.Print("Segment Request:")
		log.Println(regexp.MatchString("^/video/[a-zA-B0-9]+/[a-zA-B0-9]+/?$", r.URL.Path))
	
		if r.URL.Path == "/video/" {
			controllers.GetVideos(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9]+/?$", r.URL.Path); err == nil && matched {
			controllers.GetManifestFile(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9]+/[a-zA-B0-9_]+.ts/?$", r.URL.Path); err == nil && matched {
			controllers.GetTSFiles(w, r, db)
		}
	}
}

var validPath = regexp.MustCompile("^/(upload)/([a-zA-Z0-9]+)$")

func homePageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		return
	} else if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/index.html"
		http.ServeFile(w, r, p)
	}
}

func videoDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		video_id := r.URL.Path[len("/video/details/"):]
		log.Println("Details of " + video_id + " requested.")
		detailsQuery, err := db.Prepare(`SELECT
			title, description
		FROM
			videos
		WHERE
			video_id=?`)
		if err != nil {
			log.Fatal(err)
		}
		defer detailsQuery.Close()
		var  title, description string
		err = detailsQuery.QueryRow(video_id).Scan(&title, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Video ID: " + video_id)
		log.Println("Title: " + title)
		log.Println("Description: " + description)
		videoDetails := &Video{
			ID : video_id,
			Title : title,
			Description : description,
		}
		videoDetailsJSON, err := json.Marshal(videoDetails)
		
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, string(videoDetailsJSON))		
	}
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/upload.html"
		http.ServeFile(w, r, p)
	}
}

func listPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/list.html"
		http.ServeFile(w, r, p)
	}
}

func watchPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/watch.html"
		http.ServeFile(w, r, p)
	}
}

func setUpRoutes(db *sql.DB) {
	log.Println("Setting up routes...")
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/upload", uploadPageHandler)
	http.HandleFunc("/list", listPageHandler)
	http.HandleFunc("/watch", watchPageHandler)
	http.HandleFunc("/video/", func(w http.ResponseWriter, r *http.Request) {
		videoHandler(w, r, db)
	})
	http.HandleFunc("/video/details/", func(w http.ResponseWriter, r *http.Request) {
		videoDetailsHandler(w, r, db)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Routes set.")
}

func initServer() {
	log.Println("Initializing server...")
	utils.LoadEnvVars()
	db := database.Connect()
	setUpRoutes(db)
	utils.ResumeUploadIfAny()
}

func main() {
	initServer()
	log.Println("Server is running on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
