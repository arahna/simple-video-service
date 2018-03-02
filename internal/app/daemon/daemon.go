package daemon

import (
	"github.com/arahna/simple-video-service/internal/pkg/model"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"github.com/arahna/simple-video-service/internal/pkg/contentserver"
	"math"
)

type VideoDaemon struct {
	Repo         model.VideoRepository
	VideoService VideoService
}

type task struct {
	video model.Video
}

func New(repo model.VideoRepository, service VideoService) *VideoDaemon {
	return &VideoDaemon{repo, service}
}

func (d *VideoDaemon) RunWorkerPool(stopChan chan struct{}) *sync.WaitGroup {
	var wg sync.WaitGroup
	taskChan := d.runTaskProvider(stopChan)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			d.worker(taskChan, i)
			wg.Done()
		}(i)
	}
	return &wg
}

func (d *VideoDaemon) runTaskProvider(stopChan chan struct{}) <-chan *task {
	resultChan := make(chan *task)
	stopTaskChan := make(chan struct{})
	taskChan := d.taskProvider(stopTaskChan)
	stop := func() {
		stopTaskChan <- struct{}{}
		close(resultChan)
	}

	go func() {
		for {
			select {
			case <-stopChan:
				stop()
				return
			case task := <-taskChan:
				select {
				case <-stopChan:
					stop()
					return
				case resultChan <- task:
				}

			}
		}
	}()

	return resultChan
}

func (d *VideoDaemon) taskProvider(stopChan chan struct{}) <-chan *task {
	taskChan := make(chan *task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(taskChan)
				return
			default:
				if task := d.generateTask(); task != nil {
					log.WithField("video id", task.video.Uid).Info("got the new video to process")
					taskChan <- task
				} else {
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return taskChan
}

func (d *VideoDaemon) generateTask() *task {
	if video, err := d.Repo.FindOneWithStatus(model.VideoCreated); video != nil && err == nil {
		if err2 := d.updateStatus(video, model.VideoProcessing, nil); err2 != nil {
			log.WithFields(log.Fields{"video id": video.Uid,"err": err2}).Error("failed to update status to processing")
			return nil
		}
		return &task{*video}
	}
	return nil
}

func (d *VideoDaemon) worker(taskChan <-chan *task, name int) {
	log.Printf("start worker %v\n", name)
	for task := range taskChan {
		video := &task.video
		log.WithFields(log.Fields{"video id": video.Uid, "worker": name}).Info("start processing")
		if err := d.processTask(task); err != nil {
			log.WithFields(log.Fields{"video id": video.Uid, "name": name, "err": err}).Error("processing video error")
		} else {
			log.WithFields(log.Fields{"video id": video.Uid, "name": name}).Info("processing video success")
		}
	}
}

func (d *VideoDaemon) processTask(t *task) error {
	video := t.video
	err := d.updateDuration(&video)
	err = d.createThumbnail(&video, err)
	err = d.updateStatus(&video, model.VideoReady, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *VideoDaemon) updateDuration(video *model.Video) error {
	duration, err := d.VideoService.Duration(contentserver.GetVideoPath(video.Uid, video.FileName))
	if err != nil {
		video.Status = model.VideoError
		if err2 := d.Repo.Save(video); err2 != nil {
			return err2
		}
		return err
	}
	video.Duration = uint(duration)
	return nil
}

func (d *VideoDaemon) createThumbnail(video *model.Video, err error) error {
	if err != nil {
		return err
	}
	videoPath := contentserver.GetVideoPath(video.Uid, video.FileName)
	thumbPath := contentserver.GetThumbnailPath(video.Uid)
	err = d.VideoService.CreateVideoThumbnail(videoPath, thumbPath, uint64(math.Min(1.00, float64(video.Duration))))
	if err != nil {
		return err
	}
	return nil
}

func (d *VideoDaemon) updateStatus(video *model.Video, status model.VideoStatus, err error) error {
	if err != nil {
		return err
	}
	video.Status = model.VideoReady
	if err = d.Repo.Save(video); err != nil {
		return err
	}
	return nil
}