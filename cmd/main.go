package main

import (
	"context"
	"github.com/anonimpopov/WTFTest/config"
	"github.com/anonimpopov/WTFTest/internal/handlers"
	"github.com/anonimpopov/WTFTest/internal/repository/secondRealisation"
	"github.com/anonimpopov/WTFTest/internal/repository/thirdRealisation"
	"github.com/anonimpopov/WTFTest/internal/server"
	"github.com/anonimpopov/WTFTest/internal/service/metric"
	"github.com/anonimpopov/WTFTest/internal/service/metricBatch"
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
	db := mongoClient.Database(cfg.Mongo.Database)

	//metricsRepo := firstRealisation.New(db.Collection("pixi1"))
	metricsRepo := secondRealisation.New(db.Collection("pixi2"))
	metricsBatchRepo := thirdRealisation.New(db.Collection("pixi3"))
	stopMetricBatchChan := metricsBatchRepo.Init()

	if err := metricsRepo.Init(); err != nil {
		logrus.Fatalf("error during init repo: %v", err)
	}

	metricsService := metric.New(metricsRepo)
	meticsBatchService := metricBatch.New(metricsBatchRepo)
	router := handlers.New(metricsService, meticsBatchService)

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
	stopMetricBatchChan <- true
}
