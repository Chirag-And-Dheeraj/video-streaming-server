package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"time"
	"video-streaming-server/config"
	"video-streaming-server/controllers"
	"video-streaming-server/database"
	"video-streaming-server/middleware"
	"video-streaming-server/repositories"
	"video-streaming-server/services"
	"video-streaming-server/shared"
	"video-streaming-server/sse"
	"video-streaming-server/types"
	"video-streaming-server/utils"
)

/*
	/video/ - Get All Videos
	/video/[id] - Get A Video
	/video/[id]/stream - Get The Manifest File For The Video
	/video/[id]/stream/[filename] - Get The Segment of Video
*/

func getRequestLogger(r *http.Request, logger *slog.Logger) (*slog.Logger, error) {
	user, err := utils.GetUserFromRequest(r)
	if err != nil {
		logger.Error("Failed to get user from request", "error", err)
		return nil, err
	}

	requestGroup := slog.Group(
		"request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user", user.ID),
	)

	return logger.With(requestGroup), nil
}

func videoHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger) {
	videoLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	videoLogger.Info("Received request for video endpoint.")

	path := r.URL.Path
	method := r.Method

	if method == http.MethodPost {
		controllers.UploadVideo(w, r, db)
	} else if method == http.MethodGet {
		if path == "/video/" {
			videoLogger.Info("Get all video details")
			controllers.GetVideos(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/?$", path); err == nil && matched {
			videoLogger.Info("Get video details")
			controllers.GetVideo(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/?$", path); err == nil && matched {
			videoLogger.Info("Manifest Request")
			controllers.ManifestFileHandler(w, r, db)
		} else if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/stream/[a-zA-B0-9_-]+.ts/?$", r.URL.Path); err == nil && matched {
			videoLogger.Info("Segment Request")
			controllers.TSFileHandler(w, r, db)
		} else {
			response := fmt.Sprintf("Error: handler for %s not found", html.EscapeString(r.URL.Path))
			http.Error(w, response, http.StatusNotFound)
		}
	} else if method == http.MethodDelete {
		videoLogger.Info("Delete video details")
		if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/?$", path); err == nil && matched {
			controllers.DeleteHandler(w, r, db)
		}
	} else if method == http.MethodPatch {
		videoLogger.Info("Update video details")
		if matched, err := regexp.MatchString("^/video/[a-zA-B0-9-]+/?$", path); err == nil && matched {
			controllers.UpdateHandler(w, r, db)
		}
	}
}

type HomePageData struct {
	IsLoggedIn bool
	Username   string
}

func homePageHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	homePageLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	homePageLogger.Info("Received request for home page.")

	if r.URL.Path != "/" {
		response := fmt.Sprintf("Error: handler for %s not found", html.EscapeString(r.URL.Path))
		http.Error(w, response, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
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

func registerPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger) {
	registerLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	registerLogger.Info("Received request for register page.")

	if r.Method == http.MethodPost {
		userRepository := repositories.NewUserRepository(db)
		userService := services.NewUserService(userRepository)
		controllers.RegisterUser(w, r, userService)
	} else if r.Method == http.MethodGet {
		p := "./client/register.html"
		http.ServeFile(w, r, p)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger) {
	loginLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	loginLogger.Info("Received request for login page.")

	if r.Method == http.MethodPost {
		userRepository := repositories.NewUserRepository(db)
		userService := services.NewUserService(userRepository)
		controllers.LoginUser(w, r, userService)
	} else if r.Method == http.MethodGet {
		p := "./client/login.html"
		http.ServeFile(w, r, p)
	}
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	uploadLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	uploadLogger.Info("Received request for upload page.")

	if r.Method == http.MethodGet {
		log.Println("GET: " + r.URL.Path)
		p := "./client/upload.html"
		http.ServeFile(w, r, p)
	}
}

func listPageHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	listLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	listLogger.Info("Received request for list page.")

	if r.Method == http.MethodGet {
		log.Println("GET: " + r.URL.Path)
		p := "./client/list.html"
		http.ServeFile(w, r, p)
	}
}

func watchPageHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	watchLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	watchLogger.Info("Received request for watch page.")

	if r.Method == http.MethodGet {
		log.Println("GET: " + r.URL.Path)
		p := "./client/watch.html"
		http.ServeFile(w, r, p)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	configLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	configLogger.Info("Received request for config endpoint.")

	if r.Method != http.MethodGet {
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	logoutLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logoutLogger.Info("Received request for logout.")

	if r.Method == http.MethodPost {
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

func serverSentEventsHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	sseLogger, err := getRequestLogger(r, logger)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sseLogger.Info("Received request for server-sent events.")
	// TODO: talk to Jaden about security review
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := utils.GetUserFromRequest(r)
	if err != nil {
		log.Printf("failed to extract user from request: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	path, err := utils.GetRefererPathFromRequest(r)
	if err != nil {
		log.Printf("failed to extract referer path from request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID := types.UserID(user.ID)
	sessionID := types.SessionID(sse.InitializeSSEConnection(userID, path))

	// defer sachMeDeferWaleFunctionsReturnKePehleExecuteHoteHainTestingFunction()
	defer sse.RemoveSSEConnection(userID, sessionID)
	log.Printf("new sessionID is: %s", sessionID)
	utils.PrettyPrintMap(shared.GlobalUserSSEConnectionsMap, "GlobalUserSSEConnectionsMap")

	log.Printf("user %s connected to the event stream", user.Username)

	disconnected := r.Context().Done()
	rc := http.NewResponseController(w)

	// the intent behind this flush is for an initial connection heartbeat
	rc.Flush()

	channel := shared.GlobalUserSSEConnectionsMap[userID].Sessions[sessionID].EventChannel

	for {
		select {
		case <-disconnected:
			log.Printf("user %s session %s disconnected", userID, sessionID)
			return
		case rawMessage := <-channel:

			message := rawMessage
			event := message.Event
			rawData := message.Data
			eventData, err := json.Marshal(rawData)
			if err != nil {
				log.Printf("failed marshalling data for event %s to user %s on session %s", event, userID, sessionID)
				continue
			}
			log.Printf("sending event %s data %s to user %s on session %s", event, eventData, userID, sessionID)
			if n, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, eventData); err != nil {
				log.Printf("Unable to write: %s", err.Error())
			} else {
				log.Printf("Wrote %d bytes", n)
			}
			rc.Flush()
		}
	}
}

func setUpRoutes(db *sql.DB, logger *slog.Logger) {
	slog.Info("Setting up routes...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homePageHandler(w, r, logger)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerPageHandler(w, r, db, logger)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginPageHandler(w, r, db, logger)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, logger)
	})
	http.HandleFunc("/upload", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		uploadPageHandler(w, r, logger)
	}))
	http.HandleFunc("/list", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		listPageHandler(w, r, logger)
	}))
	http.HandleFunc("/watch", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		watchPageHandler(w, r, logger)
	}))
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		configHandler(w, r, logger)
	})
	http.HandleFunc("/video/", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		videoHandler(w, r, db, logger)
	}))
	http.HandleFunc("/server-events/", middleware.AuthRequired(func(w http.ResponseWriter, r *http.Request) {
		serverSentEventsHandler(w, r, logger)
	}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	slog.Info("Routes set.")
}

func initServer(logger *slog.Logger) {
	slog.Info("Initializing server...")
	utils.LoadEnvVars()
	database.DB = database.GetDBConn()
	setUpRoutes(database.DB, logger)
	// utils.ResumeUploadIfAny(database.DB)
}

func main() {
	// file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// // TODO: I want to see if I dont close the file, what and how memory leaks
	// defer file.Close()

	// w := io.MultiWriter(os.Stdout, file)

	handlerOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	if err := config.LoadConfig(".env"); err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	initServer(logger)
	slog.Info(
		"Server is listening on",
		slog.String("address", config.AppConfig.Addr),
		slog.String("port", config.AppConfig.Port),
	)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.AppConfig.Addr, config.AppConfig.Port), nil))
}
