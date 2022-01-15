package main

import (
	"context"
	"github.com/anonimpopov/WTFTest/internal/handlers"
	"github.com/anonimpopov/WTFTest/internal/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("Server starting")
	router := handlers.New()

	srv, shutdownChan := server.New(":8080", router.InitRouter())

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Errorf("ListenAndServe error: %v", err)
		}
	}()

	<-shutdownChan
	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

}
