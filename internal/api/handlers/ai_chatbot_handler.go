//TODO: remove after testing AI Chatbot done

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type AIChatbotHandler struct {
	aiChatbotService *services.AIChatbotService
}

func NewAIChatbotHandler(aiChatbotService *services.AIChatbotService) *AIChatbotHandler {
	return &AIChatbotHandler{aiChatbotService: aiChatbotService}
}

func (h *AIChatbotHandler) AddMessage(c *gin.Context) {
	var req dto.AddMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.aiChatbotService.AddAndRunMessage(c.Request.Context(), req.ChannelID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
