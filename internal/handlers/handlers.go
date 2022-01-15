package handlers

import (
	"github.com/anonimpopov/WTFTest/internal/service/metric"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	metricsService *metric.Service
}

func New(ms *metric.Service) *handler {
	return &handler{ms}
}

func (h *handler) InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/counter.gif", metricHandler)
	return router
}

func metricHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
