package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"database/sql"
	"context"
	"encoding/json"
)

const contextDbKey = "db"

func Router(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", video).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideo).Methods(http.MethodPost)

	return dbMiddleware(logMiddleware(r), db)
}

func dbMiddleware(h http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextDbKey, db)))
	})
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

func getDatabase(r *http.Request) *sql.DB {
	db, ok := r.Context().Value(contextDbKey).(*sql.DB)
	if !ok {
		log.WithField("context", r.Context()).Error("Can't get database from context")
		return nil
	}
	return db
}