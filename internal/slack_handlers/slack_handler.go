package slack_handlers

import (
	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type SlackHandler struct {
	slackClient      *slack.Client
	aiChatbotService *services.AIChatbotService
}

func NewSlackHandler(slackClient *slack.Client, aiChatbotService *services.AIChatbotService) *SlackHandler {
	return &SlackHandler{slackClient: slackClient, aiChatbotService: aiChatbotService}
}
