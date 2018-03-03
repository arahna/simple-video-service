package daemon

import (
	log "github.com/sirupsen/logrus"
	"github.com/arahna/simple-video-service/internal/pkg/model"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
	"math"
)

type workerPool struct {
	videoService VideoService
	repo         model.VideoRepository
}

func newWorkerPool(service VideoService, repo model.VideoRepository) *workerPool {
	return &workerPool{service, repo}
}

func (w *workerPool) worker(taskChan <-chan *task, name int) {
	log.Printf("start workerPool %v\n", name)
	for task := range taskChan {
		video := &task.video
		log.WithFields(log.Fields{"video id": video.Uid, "workerPool": name}).Info("start processing")
		if err := w.processTask(task); err != nil {
			log.WithFields(log.Fields{"video id": video.Uid, "name": name, "err": err}).Error("processing video error")
		} else {
			log.WithFields(log.Fields{"video id": video.Uid, "name": name}).Info("processing video success")
		}
	}
}

func (w *workerPool) processTask(t *task) error {
	video := t.video
	err := w.updateDuration(&video)
	err = w.createThumbnail(&video, err)
	err = w.updateStatus(&video, model.VideoReady, err)
	if err != nil {
		return err
	}
	return nil
}

func (w *workerPool) updateDuration(video *model.Video) error {
	duration, err := w.videoService.Duration(contentserver.GetVideoPath(video.Uid, video.FileName))
	if err != nil {
		video.Status = model.VideoError
		if err2 := w.repo.Save(video); err2 != nil {
			return err2
		}
		return err
	}
	video.Duration = uint(duration)
	return nil
}

func (w *workerPool) createThumbnail(video *model.Video, err error) error {
	if err != nil {
		return err
	}
	videoPath := contentserver.GetVideoPath(video.Uid, video.FileName)
	thumbPath := contentserver.GetThumbnailPath(video.Uid)
	err = w.videoService.CreateVideoThumbnail(videoPath, thumbPath, uint64(math.Min(1.00, float64(video.Duration))))
	if err != nil {
		return err
	}
	return nil
}

func (w *workerPool) updateStatus(video *model.Video, status model.VideoStatus, err error) error {
	if err != nil {
		return err
	}
	video.Status = model.VideoReady
	if err = w.repo.Save(video); err != nil {
		return err
	}
	return nil
}