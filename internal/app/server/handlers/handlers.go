package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"encoding/json"
	"github.com/arahna/simple-video-service/internal/pkg/model"
)

func Router(repo model.VideoRepository) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", list(repo)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", video(repo)).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideo(repo)).Methods(http.MethodPost)
	s.HandleFunc("/video/{ID}/status", status(repo)).Methods(http.MethodGet)

	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"url": r.URL,
			"remoteAddress": r.RemoteAddr,
			"userAgent": r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func writeJsonResponse(w http.ResponseWriter, data interface {}) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"data": data,
			}).Error("Failed to json encode")
		http.Error(w, "", http.StatusInternalServerError)
		return false
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return writeSuccessResponse(w, string(jsonData))
}

func writeSuccessResponse(w http.ResponseWriter, response string) bool {
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, response); err != nil {
		log.WithField("err", err).Error("write response error")
		return false
	}
	return true
}

func writeInternalServerError(w http.ResponseWriter, err error, message string) {
	http.Error(w, "", http.StatusInternalServerError)

	if err != nil {
		log.WithField("err", err).Error(message)
	} else if message != "" {
		log.Error(message)
	}
}