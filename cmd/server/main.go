package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"os"
	"context"
	"database/sql"
	"github.com/arahna/simple-video-service/internal/app/server/handlers"
	"github.com/arahna/simple-video-service/internal/pkg/database"
	"github.com/arahna/simple-video-service/internal/pkg/killsignal"
	"github.com/arahna/simple-video-service/configs"
	"path/filepath"
)

const logFileName = "server.log"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(filepath.Join(configs.LogDir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}

	db, err := database.InitDatabase(configs.DatabaseSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	kc := killsignal.NewChan()
	srv := startServer(":8000", db)
	killsignal.Wait(kc)
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