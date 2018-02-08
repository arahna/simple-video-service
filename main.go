package main

import (
	"net/http"
	"github.com/arahna/simple-video-service/handlers"
	log "github.com/sirupsen/logrus"
	"os"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	killSignalChan := getKillSignalChan()
	srv := startServer(":8000")

	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}

func startServer(serverUrl string) *http.Server {
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")

	router := handlers.Router()
	srv := &http.Server{Addr: serverUrl, Handler: router}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <- chan os.Signal) {
	killSignal := <- killSignalChan

	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}