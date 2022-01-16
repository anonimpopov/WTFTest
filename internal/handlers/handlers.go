package handlers

import (
	"github.com/anonimpopov/WTFTest/internal/models"
	"github.com/anonimpopov/WTFTest/internal/service/metric"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	router.GET("/counter.gif", h.saveMetric)
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
		if err := h.metricsService.SaveMetric(c, models.Action{Country: country, Type: actionType}); err != nil {
			logrus.Errorf("error during save metric country: %v, action: %v, err: %v", country, actionType, err)
		}
	}()
}
