package controllers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"video-streaming-server/config"
	"video-streaming-server/shared/logger"
	. "video-streaming-server/types"
	"video-streaming-server/utils"
)

// @desc Create new video resource
// @route POST /video
func UploadVideo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fileName := r.Header.Get("file-name")
	isFirstChunk := r.Header.Get("first-chunk")
	fileSize, _ := strconv.Atoi(r.Header.Get("file-size"))
	user, err := utils.GetUserFromRequest(r)
	title := r.Header.Get("title")

	userID := UserID(user.ID)

	if err != nil {
		logger.Log.Warn("failed to get user from request", "error", err)
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	sizeLimit, _ := strconv.Atoi(config.AppConfig.FileSizeLimit)

	if fileSize > sizeLimit {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File size greater than 15 MB is not acceptable"))
		return
	}

	serverFileName := fileName + ".mp4"

	if isFirstChunk == "true" {
		description := r.Header.Get("description")

		insertStatement, err := db.Prepare(`
			INSERT INTO 
				videos
					(
						video_id,
						title,
						description,
						upload_initiate_time,
						status,
						delete_flag,
						user_id
					) 
				VALUES 
					($1,$2,$3,$4,$5,$6,$7)
		`)

		if err != nil {
			logger.Log.Error("failed to prepare insert statement", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		_, err = insertStatement.Exec(fileName, title, description, time.Now(), 0, 0, user.ID)

		if err != nil {
			logger.Log.Error("failed to execute insert statement", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	d, _ := io.ReadAll(r.Body)

	var tmpFile *os.File

	if isFirstChunk == "true" {
		tmpFile, err = os.Create("./video/" + serverFileName)
		if err != nil {
			logger.Log.Error("failed to create file", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Error processing file")
			return
		}
	} else {
		tmpFile, err = os.OpenFile("./video/"+serverFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logger.Log.Error("failed to open file for appending", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	_, err = tmpFile.Write(d)

	if err != nil {
		logger.Log.Error("failed to write to file", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	fileInfo, err := tmpFile.Stat()

	if err != nil {
		logger.Log.Error("failed to get file info", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if fileInfo.Size() == int64(fileSize) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Video received completely and is now being processed."))

		err := utils.UpdateVideoStatus(db, fileName, UploadedOnServer)
		if err != nil {
			logger.Log.Error("failed to update video status", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		go utils.PostUploadProcessFile(serverFileName, fileName, title, tmpFile, db, userID)

	} else {
		w.WriteHeader(http.StatusPartialContent)
		w.Write([]byte("Receiving chunks of the video."))
	}
}

// @desc Get All Videos
// @route GET /video
func GetVideos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user, err := utils.GetUserFromRequest(r)

	if err != nil {
		logger.Log.Warn("failed to get user from request", "error", err)
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	getUserVideosQuery, err := db.Prepare(`
		SELECT
			video_id,
			title,
			description,
			thumbnail,
			status
		FROM
			videos
		WHERE
			delete_flag=0
		AND
			status <> 0
		AND
			user_id=$1
		ORDER BY
			upload_initiate_time DESC;
	`)

	if err != nil {
		logger.Log.Error("failed to prepare query", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rows, err := getUserVideosQuery.Query(user.ID)

	if err != nil {
		logger.Log.Error("failed to execute query", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	defer rows.Close()

	records := make([]VideoResponseType, 0)

	for rows.Next() {
		var id string
		var title string
		var description string
		var thumbnail sql.NullString
		var status VideoStatus

		err := rows.Scan(&id, &title, &description, &thumbnail, &status)

		if err != nil {
			logger.Log.Error("failed to scan row", "error", err)
			utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		thumbValue := "../static/logo/android-chrome-192x192.png"
		if thumbnail.Valid && thumbnail.String != "" {
			thumbValue = thumbnail.String
		}

		record := VideoResponseType{
			ID:          id,
			Title:       title,
			Description: description,
			Thumbnail:   thumbValue,
			Status:      status,
		}

		records = append(records, record)
	}

	recordsJSON, err := json.Marshal(records)

	if err != nil {
		logger.Log.Error("failed to marshal records", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(recordsJSON))
}

// @desc Get a Video
// @route GET /video/[id]
func GetVideo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]

	user, err := utils.GetUserFromRequest(r)

	if err != nil {
		logger.Log.Error("failed to get user from request", "error", err)
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	detailsQuery, err := db.Prepare(`
		SELECT
			title, description
		FROM
			videos
		WHERE
			video_id=$1
		AND
			user_id=$2;
	`)

	if err != nil {
		logger.Log.Error("failed to prepare query", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	defer detailsQuery.Close()

	var title, description string
	err = detailsQuery.QueryRow(videoId, user.ID).Scan(&title, &description)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("video not found", "videoId", videoId, "userId", user.ID)
			utils.SendError(w, http.StatusNotFound, "Video not found")
			return
		}
		logger.Log.Error("failed to query video details", "error", err, "videoId", videoId)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	videoDetails := &Video{
		ID:          videoId,
		Title:       title,
		Description: description,
	}
	videoDetailsJSON, err := json.Marshal(videoDetails)

	if err != nil {
		logger.Log.Error("failed to marshal response", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(videoDetailsJSON))
}

// @desc Get Manifest File
// @route GET /video/[id]/stream
func ManifestFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := strings.Split(r.URL.Path[1:], "/")[1]

	file, err := utils.GetManifestFile(w, videoId)

	if err != nil {
		logger.Log.Error("failed to retrieve manifest file", "videoId", videoId, "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Error retrieving video")
	} else {
		w.Header().Set("Content-Type", "application/x-mpegURL")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	}
}

// @desc Get TS File
// @route GET /video/[id]/stream/[id].ts
func TSFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoName := strings.Split(r.URL.Path[1:], "/")[3]
	videoComps := strings.Split(videoName, ".")
	segment := videoComps[0]
	fileId := utils.GetFileId(segment)

	getSegmentFile := "https://cloud.appwrite.io/v1/storage/buckets/" + config.AppConfig.AppwriteBucketID + "/files/" + fileId + "/view"

	request, err := http.NewRequest(http.MethodGet, getSegmentFile, nil)

	if err != nil {
		logger.Log.Error("failed to get user from request", "segment", segment, "error", err)
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Appwrite-Response-Format", config.AppConfig.AppwriteResponseFormat)
	request.Header.Set("X-Appwrite-Project", config.AppConfig.AppwriteProjectID)
	request.Header.Set("X-Appwrite-Key", config.AppConfig.AppwriteKey)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		logger.Log.Error("failed to fetch segment file", "segment", segment, "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		logger.Log.Error("segment file not found", "segment", segment)
		utils.SendError(w, http.StatusNotFound, "Segment file not found")
		return
	}

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		logger.Log.Error("failed to read segment file", "segment", segment, "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Error reading segment file")
		return
	}

	w.Header().Set("Content-Type", "video/MP2T")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

// @desc Update Video Details
// @route UPDATE
func UpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]
	user, err := utils.GetUserFromRequest(r)
	if err != nil {
		logger.Log.Error("failed to decode request body", "error", err)
		utils.SendError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	if videoId == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	var reqBody UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.Log.Warn("failed to get user from request", "error", err)
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query, err := db.Prepare(`
		UPDATE
			videos
		SET
			title=$1,
			description=$2
		WHERE
			video_id=$3
		AND
			user_id=$4;
	`)

	if err != nil {
		logger.Log.Error("failed to prepare update query", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	result, err := query.Exec(reqBody.Title, reqBody.Description, videoId, user.ID)

	if err != nil {
		logger.Log.Error("failed to execute update query", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Log.Error("failed to get rows affected", "error", err)
		return
	}

	if rowsAffected == 0 {
		logger.Log.Info("video not found for update", "videoId", videoId, "userId", user.ID)
		utils.SendError(w, http.StatusNotFound, "Video not found")
		return
	}

	logger.Log.Info("video details updated", "videoId", videoId, "userId", user.ID)
	w.WriteHeader(http.StatusOK)
}

// @desc Delete the video
// @route DELETE /video/[id]
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	videoId := r.URL.Path[len("/video/"):]

	updateStatement, err := db.Prepare(`
		UPDATE
			videos
		SET
			delete_flag=$1
			WHERE
			video_id=$2;
	`)

	if err != nil {
		logger.Log.Error("failed to prepare update statement", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_, err = updateStatement.Exec(1, videoId)

	if err != nil {
		logger.Log.Error("failed to execute update statement", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	go utils.DeleteVideo(w, r, db, videoId)
}
