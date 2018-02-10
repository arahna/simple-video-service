package handlers

import (
	"net/http"
	"io"
	"os"
	log "github.com/sirupsen/logrus"
)

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		http.Error(w, "Invalid file format", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile(header.Filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(file, fileReader); err != nil {
		log.WithField("err", err).Error("Failed to save file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, "The video uploaded"); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}
