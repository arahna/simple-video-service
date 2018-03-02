package daemon

type VideoService interface {
	Duration(filePath string) (float64, error)
	CreateVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset uint64) error
}
