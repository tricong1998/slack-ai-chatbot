package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type SlackHandler struct {
	slackService services.ISlackService
}

func NewSlackHandler(
	slackService services.ISlackService,
) *SlackHandler {
	return &SlackHandler{
		slackService: slackService,
	}
}

func (h *SlackHandler) SendMessage(ctx *gin.Context) {
	var input dto.SendMessageDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.slackService.SendMessage(ctx, input.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message sent to Slack"})
}
