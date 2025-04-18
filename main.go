package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
	"video-streaming-server/config"
	"video-streaming-server/controllers"
	"video-streaming-server/database"
	"video-streaming-server/middleware"
	"video-streaming-server/repositories"
	"video-streaming-server/services"
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

type HomePageData struct {
	IsLoggedIn bool
	Username   string
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		response := fmt.Sprintf("Error: handler for %s not found", html.EscapeString(r.URL.Path))
		http.Error(w, response, http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles("./client/index.html"))

	isLoggedIn := false
	username := ""
	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie != nil {
		isLoggedIn = true
		claims, _ := utils.DecodeJWT(cookie.Value)
		username = claims["username"].(string)
	}

	tmpl.Execute(w, HomePageData{
		IsLoggedIn: isLoggedIn,
		Username:   username,
	})
}

func registerPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	path := r.URL.Path
	method := r.Method

	if method == "POST" {
		log.Println("POST: " + path)
		userRepository := repositories.NewUserRepository(db)
		userService := services.NewUserService(userRepository)
		controllers.RegisterUser(w, r, userService)
	} else if method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/register.html"
		http.ServeFile(w, r, p)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	path := r.URL.Path
	method := r.Method

	if method == "POST" {
		log.Println("POST: " + path)
		userRepository := repositories.NewUserRepository(db)
		userService := services.NewUserService(userRepository)
		controllers.LoginUser(w, r, userService)
	} else if method == "GET" {
		log.Println("GET: " + r.URL.Path)
		p := "./client/login.html"
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
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileSizeLimit := config.AppConfig.FileSizeLimit
	if fileSizeLimit == "" {
		http.Error(w, "File size limit not configured", http.StatusInternalServerError)
		return
	}

	supportedFileTypes := []config.FileType{
		{FileType: "video/mp4", FileExtension: ".mp4"},
		{FileType: "video/x-matroska", FileExtension: ".mkv"},
		{FileType: "video/quicktime", FileExtension: ".mov"},
	}

	response := config.ConfigResponse{
		FileSizeLimit:      fileSizeLimit,
		SupportedFileTypes: supportedFileTypes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Expires:  time.Now(),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func serverSentEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("SSE API IS CONNECTED OKAY.")

	clientGone := r.Context().Done()
	rc := http.NewResponseController(w)

	if n, err := fmt.Fprintf(w, "event:cpu\ndata:User:bhai\n\n"); err != nil {
		log.Printf("Unable to write: %s", err.Error())
	} else {
		log.Printf("Wrote %d bytes", n)
	}

	rc.Flush()

	for {
		select {
		case <-clientGone:
			log.Println("Client is ded, I am sed 😞")
			return
		case msg := <-utils.EventsChannel:
			log.Println("received message bhai " + msg)
			if n, err := fmt.Fprintf(w, "event:cpu\ndata:User:%s\n\n", msg); err != nil {
				log.Printf("Unable to write: %s", err.Error())
			} else {
				log.Printf("Wrote %d bytes", n)
			}
			rc.Flush()
		}
	}
}

func setUpRoutes(db *sql.DB) {
	log.Println("Setting up routes...")
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerPageHandler(w, r, db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPageHandler(w, r, db)
	})
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", middleware.AuthRequired(uploadPageHandler))
	http.HandleFunc("/list", middleware.AuthRequired(listPageHandler))
	http.HandleFunc("/watch", middleware.AuthRequired(watchPageHandler))
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/video/", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		videoHandler(w, r, db)
	}))
	http.HandleFunc("/server-events/", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		serverSentEventsHandler(w, r)
	}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Routes set.")
}

func initServer() {
	log.Println("Initializing server...")
	utils.LoadEnvVars()
	database.DB = database.GetDBConn()
	setUpRoutes(database.DB)
	utils.ResumeUploadIfAny(database.DB)
	utils.EventsChannel = make(chan string)
	fmt.Println(utils.EventsChannel)
}

func main() {
	if err := config.LoadConfig(".env"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	initServer()
	log.Println("Server is listening on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.AppConfig.Addr, config.AppConfig.Port), nil))
}
