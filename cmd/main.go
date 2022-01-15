package main

import (
	"context"
	"github.com/anonimpopov/WTFTest/config"
	"github.com/anonimpopov/WTFTest/internal/handlers"
	"github.com/anonimpopov/WTFTest/internal/server"
	"github.com/anonimpopov/WTFTest/pkg/mongo"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	logrus.Info("Server starting")
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("cant load config: %v", err)
	}

	mongoClient, err := mongo.GetMongoClient(cfg.Mongo.URL)
	if err != nil {
		logrus.Fatalf("cant connect to mongodb: %v", err)
	}

	router := handlers.New()

	srv, shutdownChan := server.New(cfg.Server.Port, router.InitRouter())

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Errorf("ListenAndServe error: %v", err)
		}
	}()

	<-shutdownChan
	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	srv.Shutdown(shutdownContext)
	if err := mongoClient.Disconnect(shutdownContext); err != nil {
		logrus.Errorf("error occured on mongodb connection close: %v", err)
	}
}
