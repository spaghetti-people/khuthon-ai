package api

import (
	"net/http"

	"github.com/spaghetti-people/khuthon-ai/internal/api/handler"
	"github.com/spaghetti-people/khuthon-ai/internal/app"

	"github.com/gin-gonic/gin"
)

func NewRouter(container *app.Container) *gin.Engine {
	r := gin.Default()

	// health check page
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 식물 채팅
	r.POST("/plant/chat", handler.NewPlantChatHandler(container.PlantChat).Handle)

	return r
}
