package ffmpeg

import (
	"github.com/arahna/simple-video-service/internal/app/daemon"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ffmpeg struct {
}

func New() daemon.VideoService {
	return &ffmpeg{}
}

func (ffmpeg *ffmpeg) Duration(videoPath string) (float64, error) {
	result, err := exec.Command(`ffprobe`, `-v`, `error`, `-show_entries`, `format=duration`, `-of`, `default=noprint_wrappers=1:nokey=1`, videoPath).Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func (ffmpeg *ffmpeg) CreateVideoThumbnail(videoPath string, thumbPath string, thumbSecondsOffset uint64) error {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-ss`, ffmpegTimeFromSeconds(thumbSecondsOffset), `-vframes`, `1`, thumbPath).Run()
}

func ffmpegTimeFromSeconds(seconds uint64) string {
	return time.Unix(int64(seconds), 0).UTC().Format(`15:04:05.000000`)
}

