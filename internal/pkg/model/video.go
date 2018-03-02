package model

import (
	"github.com/google/uuid"
)

type VideoStatus uint8

const (
	VideoCreated    VideoStatus = 1 + iota
	VideoProcessing
	VideoReady
	VideoDeleted
	VideoError
)

type Video struct {
	Id       uint
	Uid      string
	Status   VideoStatus
	Title    string
	FileName string
	Duration uint
}

func NewVideo(title, fileName string, status VideoStatus) *Video {
	uid := uuid.New()
	return &Video{
		Uid:      uid.String(),
		Title:    title,
		FileName: fileName,
		Duration: 0,
		Status:   status,
	}
}
