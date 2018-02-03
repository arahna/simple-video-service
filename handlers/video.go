package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/sirupsen/logrus"
)

type Video struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Url       string `json:"url"`
}

func video(w http.ResponseWriter, _ *http.Request) {
	item := Video{
		"d290f1ee-6c54-4b01-90e6-d701748f0851",
		"Black Retrospetive Woman",
		15,
		"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
		"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/index.mp4"}
	b, err := json.Marshal(item)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, string(b)); err != nil {
		logrus.WithField("err", err).Error("write response error")
	}
}
