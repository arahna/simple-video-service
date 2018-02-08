package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/arahna/simple-video-service/videodb"
	"github.com/arahna/simple-video-service/contentserver"
)

type videoItem struct {
	videoListItem
	Url string `json:"url"`
}

func video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["ID"]

	if !ok || id == "" {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	video, found := videodb.Find(id)

	if !found {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	item := toVideoItem(video)
	jsonResponse, err2 := json.Marshal(item)

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, jsonResponse)
}

func toVideoItem(video videodb.Video) videoItem {
	return videoItem{
		toVideoListItem(video),
		contentserver.GetVideoUrl(video.Id),
	}
}
