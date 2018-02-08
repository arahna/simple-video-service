package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/arahna/simple-video-service/videodb"
	"github.com/arahna/simple-video-service/contentserver"
)

type videoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func list(w http.ResponseWriter, _ *http.Request) {
	videos := videodb.GetAll()
	var items []videoListItem
	for _, video := range videos {
		items = append(items, toVideoListItem(video))
	}
	jsonResponse, err := json.Marshal(items)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, jsonResponse)
}

func toVideoListItem(video videodb.Video) videoListItem {
	return videoListItem{
		video.Id,
		video.Name,
		video.Duration,
		contentserver.GetThumbnailUrl(video.Id),
	}
}
