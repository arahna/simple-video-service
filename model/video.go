package model

import (
	"github.com/google/uuid"
)

type VideoStatus uint8

const (
	Error     VideoStatus = 1 + iota
	Uploading
	Ready
)

type Video struct {
	Id       int
	Uid      string
	Status   VideoStatus
	Title    string
	FileName string
	Duration uint
}

func NewVideo(title, fileName string, duration uint, status VideoStatus) *Video {
	uid := uuid.New()
	return &Video{
		Uid:      uid.String(),
		Title:    title,
		FileName: fileName,
		Duration: duration,
		Status:   status,
	}
}
