package main

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"video-streaming-server/controllers"
	"video-streaming-server/database"
	"video-streaming-server/utils"
)

/*
	/video/ - Get All Videos
	/video/[id] - Get A Video
	/video/[id]/stream - Get The Manifest File For The Video
	/video/[id]/stream/[filename] - Get The Segment of Video
*/

func videoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	path := r.URL.Path
	method := r.Method

	if method == "POST" {
		controllers.UploadVideo(w, r, db)
	} else if method == "GET" {
		log.Println("GET: " + path)
	
		if path == "/video/" {
			controllers.GetVideos(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9]+/?$", path); err == nil && matched {
			log.Println("Get video details")
			controllers.GetVideo(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9]+/stream/?$", path); err == nil && matched {
			log.Print("Manifest Request:")
			log.Println(regexp.MatchString("^/video/[a-zA-B0-9]+/stream/?$", path))
			controllers.GetManifestFile(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9]+/stream/[a-zA-B0-9_]+.ts/?$", r.URL.Path); err == nil && matched {
			log.Print("Segment Request:")
			log.Println(regexp.MatchString("^/video/[a-zA-B0-9]+/stream/[a-zA-B0-9_]+.ts/?$", path))
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Routes set.")
}

func initServer() {
	log.Println("Initializing server...")
	utils.LoadEnvVars()
	db := database.Connect()
	setUpRoutes(db)
	utils.ResumeUploadIfAny(db)
}

func main() {
	initServer()
	log.Println("Server is running on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
