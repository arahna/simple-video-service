package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/arahna/simple-video-service/internal/pkg/model"
)

type statusResponse struct {
	Status model.VideoStatus `json:"status"`
}

func status(repository model.VideoRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["ID"]

		if !ok || id == "" {
			http.Error(w, "Invalid parameter value", http.StatusBadRequest)
			return
		}

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
}
