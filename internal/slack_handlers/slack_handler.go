package slack_handlers

import (
	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type SlackHandler struct {
	slackClient      *slack.Client
	slackService     *services.SlackService
	aiChatbotService *services.AIChatbotService
	ggSheetService   *services.GSheetService
	uiPathJobService *services.UIPathJobService
}

func NewSlackHandler(slackClient *slack.Client, slackService *services.SlackService, aiChatbotService *services.AIChatbotService, ggSheetService *services.GSheetService, uiPathJobService *services.UIPathJobService) *SlackHandler {
	return &SlackHandler{slackClient: slackClient, slackService: slackService, aiChatbotService: aiChatbotService, ggSheetService: ggSheetService, uiPathJobService: uiPathJobService}
}
