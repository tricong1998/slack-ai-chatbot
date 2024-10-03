package slack_handlers

import (
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	slackClient *slack.Client
}

func NewSlackHandler(slackClient *slack.Client) *SlackHandler {
	return &SlackHandler{slackClient: slackClient}
}
