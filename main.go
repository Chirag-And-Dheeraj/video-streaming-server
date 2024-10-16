package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
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
		log.Println("POST: " + path)
		controllers.UploadVideo(w, r, db)
	} else if method == "GET" {
		log.Println("GET: " + path)

		if path == "/video/" {
			log.Println("Get all video details")
			controllers.GetVideos(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/?$", path); err == nil && matched {
			log.Println("Get video details")
			controllers.GetVideo(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/?$", path); err == nil && matched {
			log.Print("Manifest Request:")
			log.Println(regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/?$", path))
			controllers.ManifestFileHandler(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/[a-zA-B0-9_-]+.ts/?$", r.URL.Path); err == nil && matched {
			log.Print("Segment Request:")
			log.Println(regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/[a-zA-B0-9_-]+.ts/?$", path))
			controllers.TSFileHandler(w, r, db)
		} else {
			response := fmt.Sprintf("Error: handler for %s not found", html.EscapeString(r.URL.Path))
			http.Error(w, response, http.StatusNotFound)
		}
	} else if method == "DELETE" {
		if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/?$", path); err == nil && matched {
			log.Println("DELETE: " + path)
			controllers.DeleteHandler(w, r, db)
		}
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		response := fmt.Sprintf("Error: handler for %s not found", html.EscapeString(r.URL.Path))
		http.Error(w, response, http.StatusNotFound)
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

func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := map[string]interface{}{
			"FILE_SIZE_LIMIT": os.Getenv("FILE_SIZE_LIMIT"),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func setUpRoutes(db *sql.DB) {
	log.Println("Setting up routes...")
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/upload", uploadPageHandler)
	http.HandleFunc("/list", listPageHandler)
	http.HandleFunc("/watch", watchPageHandler)
	http.HandleFunc("/config", configHandler)
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
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	log.Println("Server is listening on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), nil))
}
