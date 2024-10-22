package slack_handlers

import "context"

func (s *SlackHandler) handleCandidateSheetEvent(channelID string) error {
	return s.slackService.SendCandidateFileForm(context.Background(), channelID)
}
