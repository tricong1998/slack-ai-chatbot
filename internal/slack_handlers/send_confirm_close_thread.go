package slack_handlers

import (
	"context"
)

func (s *SlackHandler) SendConfirmCloseThread(channelID string, threadID string) error {
	return s.slackService.SendConfirmCloseThread(context.Background(), channelID, threadID)
}
