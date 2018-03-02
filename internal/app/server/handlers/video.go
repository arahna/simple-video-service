package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
	"github.com/arahna/simple-video-service/internal/pkg/model"
)

type videoItem struct {
	videoListItem
	Url string `json:"url"`
}

func video(w http.ResponseWriter, r *http.Request) {
	db := getDatabase(r)
	if db == nil {
		writeInternalServerError(w, nil, "")
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["ID"]

	if !ok || id == "" {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	repository := model.NewVideoRepository(db)
	video, err := repository.Find(id)
	if err != nil {
		writeInternalServerError(w, err, "Failed to find video")
		return
	}

	if video == nil {
		http.Error(w, "The video not found", http.StatusNotFound)
		return
	}

	item := toVideoItem(video)
	writeJsonResponse(w, item)
}

func toVideoItem(video *model.Video) videoItem {
	return videoItem{
		toVideoListItem(video),
		contentserver.GetVideoUrl(video.Uid, video.FileName),
	}
}
