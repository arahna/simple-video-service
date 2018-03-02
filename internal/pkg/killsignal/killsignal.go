package killsignal

import (
	"os"
	"os/signal"
	"syscall"
	log "github.com/sirupsen/logrus"
)

func NewChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func Wait(killSignalChan <- chan os.Signal) {
	killSignal := <- killSignalChan

	// todo: refactoring logging
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
