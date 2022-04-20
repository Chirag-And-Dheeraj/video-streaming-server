package models

type Video struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	VideoURL    string `json:"video_url"`
	Description string `json:"description"`
}