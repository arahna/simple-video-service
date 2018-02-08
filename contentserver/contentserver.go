package contentserver

import "fmt"

const thumbnailUrlPattern = "/content/%v/screen.jpg"
const videoUrlPattern = "/content/%v/index.mp4"


func GetThumbnailUrl(id string) string {
	return fmt.Sprintf(thumbnailUrlPattern, id)
}

func GetVideoUrl(id string) string {
	return fmt.Sprintf(videoUrlPattern, id)
}