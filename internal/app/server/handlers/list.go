package handlers

import (
	"net/http"
	"github.com/arahna/simple-video-service/internal/pkg/model"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
	"strconv"
	"strings"
)

const searchParam = "searchString"
const skipParam = "skip"
const limitParam = "limit"

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

	q := r.URL.Query()
	search := strings.TrimSpace(q.Get(searchParam))
	skip, err := strconv.Atoi(q.Get(skipParam))
	if err != nil {
		skip = 0
	}
	limit, err := strconv.Atoi(q.Get(limitParam))
	if err != nil {
		limit = 0
	}

	repository := model.NewVideoRepository(db)
	videos, err := repository.FindWithStatus(model.VideoReady, search, skip, limit)
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

func toVideoListItem(video *model.Video) videoListItem {
	return videoListItem{
		video.Uid,
		video.Title,
		video.Duration,
		contentserver.GetThumbnailUrl(video.Uid),
	}
}
