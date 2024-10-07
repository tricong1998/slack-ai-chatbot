package services

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
)

type SlackService struct {
	slackConfig *config.SlackConfig
	slackClient *slack.Client
}

type ISlackService interface {
	SendMessage(ctx context.Context, channelID *string, message string) error
}

func NewSlackService(slackConfig *config.SlackConfig, slackClient *slack.Client) *SlackService {
	return &SlackService{
		slackConfig: slackConfig,
		slackClient: slackClient,
	}
}

func (s *SlackService) SendMessage(ctx context.Context, channelID *string, message string) error {
	attachment := slack.Attachment{
		Pretext: message,
	}

	fmt.Println("message", message)
	fmt.Println("s", s)
	fmt.Println("slackConfig", s.slackConfig)
	channel := channelID
	if channel == nil {
		channel = &s.slackConfig.Channel
	}
	_, _, err := s.slackClient.PostMessage(*channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return err
	}

	return nil
}
