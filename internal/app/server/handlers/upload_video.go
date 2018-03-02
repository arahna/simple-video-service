package handlers

import (
	"net/http"
	"github.com/arahna/simple-video-service/internal/pkg/model"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
)

func uploadVideo(repository model.VideoRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReader, header, err := r.FormFile("file[]")
		if err != nil {
			http.Error(w, "Invalid parameters", http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != "video/mp4" {
			http.Error(w, "Invalid file format", http.StatusBadRequest)
			return
		}

		fileName := header.Filename // TODO: clean file name
		video := model.NewVideo(header.Filename, fileName, model.VideoCreated)

		err = contentserver.SaveFile(fileReader, video.Uid, video.FileName)
		if err != nil {
			video.Status = model.VideoError
			repository.Save(video)
			writeInternalServerError(w, err, "Failed to save file to content server")
			return
		}

		if err := repository.Save(video); err != nil {
			writeInternalServerError(w, err, "Failed to save file to DB")
			return
		}

		writeSuccessResponse(w, "The video uploaded")
	}
}