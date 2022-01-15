package handlers

import "github.com/gin-gonic/gin"

type handler struct {
}

func New() *handler {
	return &handler{}
}

func (h *handler) InitRouter() *gin.Engine {
	router := gin.New()

	return router
}
