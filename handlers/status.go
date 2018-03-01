package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/arahna/simple-video-service/model"
)

type statusResponse struct {
	Status model.VideoStatus `json:"status"`
}

func status(w http.ResponseWriter, r *http.Request) {
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
		writeInternalServerError(w, err, "Failed to find the video")
		return
	}

	if video == nil {
		http.Error(w, "The video not found", http.StatusNotFound)
		return
	}

	response := statusResponse{Status: video.Status}
	writeJsonResponse(w, response)
}
