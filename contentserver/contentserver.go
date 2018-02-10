package contentserver

import (
	"fmt"
	"os"
	"path/filepath"
	"io"
)

const thumbnailUrlPattern = "/content/%s/screen.jpg"
const videoUrlPattern = "/content/%s/%s"
const dirPath = "content"


func GetThumbnailUrl(uid string) string {
	return fmt.Sprintf(thumbnailUrlPattern, uid)
}

func GetVideoUrl(uid, fileName string) string {
	return fmt.Sprintf(videoUrlPattern, uid, fileName)
}

func SaveFile(src io.Reader, uid, fileName string) error {
	file, err := createFile(uid, fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, src)
	return err
}

func createFile(uid, fileName string) (*os.File, error) {
	filePath := filepath.Join(dirPath, uid)
	if err := os.Mkdir(filePath, os.ModeDir); err != nil {
		return nil, err
	}
	filePath = filepath.Join(filePath, fileName)
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}