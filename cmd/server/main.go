package main

import (
	"net/http"
	"github.com/arahna/simple-video-service/internal/app/server/handlers"
	"github.com/arahna/simple-video-service/internal/pkg/database"
	log "github.com/sirupsen/logrus"
	"os"
	"context"
	"os/signal"
	"syscall"
	"database/sql"
)

const logFileName = "log/my.log"

func main() {
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

	killSignalChan := getKillSignalChan()
	srv := startServer(":8000", db)

	waitForKillSignal(killSignalChan)

	srv.Shutdown(context.Background())
}

func startServer(serverUrl string, db *sql.DB) *http.Server {
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")

	router := handlers.Router(db)
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