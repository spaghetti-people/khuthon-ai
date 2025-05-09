package handler

import (
	"net/http"

	"github.com/spaghetti-people/khuthon-ai/internal/model"
	"github.com/spaghetti-people/khuthon-ai/internal/service"
	"github.com/spaghetti-people/khuthon-ai/pkg/response"

	"github.com/gin-gonic/gin"
)

type PlantChatHandler struct {
	svc service.PlantChatService
}

func NewPlantChatHandler(svc service.PlantChatService) *PlantChatHandler {
	return &PlantChatHandler{svc: svc}
}

func (h *PlantChatHandler) Handle(c *gin.Context) {
	var req model.PlantChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Chat(c.Request.Context(), &req)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
