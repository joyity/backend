package main

import (
	"github.com/joyity/backend/server"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	srv := server.New(log.WithField("component", "server").Logger)
	log.Info("starting server")
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("listen and serve")
	}
	log.Info("stopped server")
}
