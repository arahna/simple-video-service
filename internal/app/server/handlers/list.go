package handlers

import (
	"net/http"
	"github.com/arahna/simple-video-service/internal/pkg/model"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
)

type videoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  uint   `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func list(w http.ResponseWriter, r *http.Request) {
	db := getDatabase(r)
	if db == nil {
		writeInternalServerError(w, nil, "")
		return
	}
	repository := model.NewVideoRepository(db)
	videos, err := repository.GetReady()
	if err != nil {
		writeInternalServerError(w, err, "Failed to get video list")
		return
	}
	items := make([]videoListItem, len(videos))
	for i, video := range videos {
		items[i] = toVideoListItem(video)
	}
	writeJsonResponse(w, items)
}

func toVideoListItem(video model.Video) videoListItem {
	return videoListItem{
		video.Uid,
		video.Title,
		video.Duration,
		contentserver.GetThumbnailUrl(video.Uid),
	}
}
