package services

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
)

type SlackService struct {
	slackConfig config.SlackConfig
	slackClient *slack.Client
}

type ISlackService interface {
	SendMessage(ctx context.Context, message string) error
}

func NewSlackService(slackConfig config.SlackConfig, slackClient *slack.Client) *SlackService {
	return &SlackService{
		slackConfig: slackConfig,
		slackClient: slackClient,
	}
}

func (s *SlackService) SendMessage(ctx context.Context, message string) error {
	attachment := slack.Attachment{
		Text: message,
	}

	channelID := s.slackConfig.Channel
	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return err
	}

	return nil
}
