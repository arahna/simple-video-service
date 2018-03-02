package main

import (
	"math/rand"
	"time"
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/arahna/simple-video-service/internal/pkg/killsignal"
	"github.com/arahna/simple-video-service/internal/app/daemon"
	"github.com/arahna/simple-video-service/internal/pkg/database"
	"github.com/arahna/simple-video-service/internal/pkg/model"
	"github.com/arahna/simple-video-service/internal/app/daemon/ffmpeg"
)

const logFileName = "log/daemon.log"

func main() {
	rand.Seed(time.Now().Unix())
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}

	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := model.NewVideoRepository(db)

	videoService := ffmpeg.New()

	kc := killsignal.NewChan()
	stopChan := make(chan struct{})
	d := daemon.New(repo, videoService)
	wg := d.RunWorkerPool(stopChan)
	killsignal.Wait(kc)
	stopChan <- struct{}{}
	wg.Wait()
}