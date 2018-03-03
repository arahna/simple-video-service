package daemon

import (
	"github.com/arahna/simple-video-service/internal/pkg/model"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type VideoDaemon struct {
	repo         model.VideoRepository
	videoService VideoService
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
	wp := newWorkerPool(d.videoService, d.repo)
	for i := 0; i < 3; i++ {
		go func(i int) {
			wg.Add(1)
			wp.worker(taskChan, i)
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
	if video, err := d.repo.FindOneWithStatus(model.VideoCreated); video != nil && err == nil {
		if err2 := d.updateStatus(video, model.VideoProcessing, nil); err2 != nil {
			log.WithFields(log.Fields{"video id": video.Uid,"err": err2}).Error("failed to update status to processing")
			return nil
		}
		return &task{*video}
	}
	return nil
}

func (d *VideoDaemon) updateStatus(video *model.Video, status model.VideoStatus, err error) error {
	if err != nil {
		return err
	}
	video.Status = model.VideoReady
	if err = d.repo.Save(video); err != nil {
		return err
	}
	return nil
}