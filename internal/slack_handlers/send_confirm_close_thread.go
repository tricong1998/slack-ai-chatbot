package slack_handlers

import (
	"context"
	"fmt"
)

func (s *SlackHandler) SendConfirmCloseThread(channelID string, threadID string) error {
	fmt.Println("SendConfirmCloseThread")
	fmt.Println("channelID", channelID)
	fmt.Println("threadID", threadID)
	return s.slackService.SendConfirmCloseThread(context.Background(), channelID, threadID)
}
