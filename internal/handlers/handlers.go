package handlers

import (
	"github.com/anonimpopov/WTFTest/internal/models"
	"github.com/anonimpopov/WTFTest/internal/service/metric"
	"github.com/anonimpopov/WTFTest/internal/service/metricBatch"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	metricsService      *metric.Service
	metricsBatchService *metricBatch.Service
}

func New(ms *metric.Service, mbs *metricBatch.Service) *handler {
	return &handler{ms, mbs}
}

func (h *handler) InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/counter.gif", h.saveMetric)
	router.GET("/metrics", h.getMetrics)
	return router
}

func (h *handler) saveMetric(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "image/gif")
	_, err := c.Writer.Write([]byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x21, 0xF9, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x00})
	if err != nil {
		logrus.Errorf("error during writint responce: %v", err)
	}

	go func() {
		country, exists := c.GetQuery("country")
		if !exists {
			country = "unset"
		}

		actionType, exists := c.GetQuery("action")
		if !exists {
			actionType = "unset"
		}
		//if err := h.metricsService.SaveMetric(c, models.Action{Country: country, Type: actionType}); err != nil {
		//	logrus.Errorf("error during save metric country: %v, action: %v, err: %v", country, actionType, err)
		//}

		if err := h.metricsBatchService.SaveMetric(models.Action{Country: country, Type: actionType}); err != nil {
			logrus.Errorf("error during save metric country: %v, action: %v, err: %v", country, actionType, err)
		}
	}()
}

func (h *handler) getMetrics(c *gin.Context) {
	tmpFrom, _ := c.GetQuery("from")
	tmpTo, _ := c.GetQuery("to")

	from, err := strconv.ParseInt(tmpFrom, 10, 64)
	if err != nil {
		from = time.Now().Add(-24 * time.Hour).Unix()
	}
	to, err := strconv.ParseInt(tmpTo, 10, 64)
	if err != nil {
		to = time.Now().Unix()
	}
	metrics, err := h.metricsBatchService.GetMetrics(c, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/json")
	_, err = c.Writer.Write(metrics)
	if err != nil {
		logrus.Errorf("error during writint responce: %v", err)
	}
}
