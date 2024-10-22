package slack_handlers

import (
	"fmt"

	"github.com/slack-go/slack"
)

func (s *SlackHandler) HandleBlockAction(payload slack.InteractionCallback) (string, error) {
	for _, action := range payload.ActionCallback.BlockActions {
		switch action.ActionID {
		case "submit":
			return "", s.handleCandidateSheetSubmission(payload)
			// ... handle other action IDs as needed ...
		}
	}
	return "", nil
}

func (s *SlackHandler) handleCandidateSheetSubmission(payload slack.InteractionCallback) error {
	// Extract the submitted link
	submittedLink := payload.BlockActionState.Values["candidate_file"]["candidate_file_input"].Value

	// Process the submitted link (e.g., validate, store, etc.)
	// For now, we'll just send a confirmation message
	responseMsg := fmt.Sprintf("Received candidate sheet link: %s", submittedLink)

	responseBlocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", responseMsg, false, false),
			nil,
			nil,
		),
	}

	_, _, err := s.slackClient.PostMessage(payload.Channel.ID, slack.MsgOptionBlocks(responseBlocks...))
	if err != nil {
		return fmt.Errorf("failed to send confirmation message: %w", err)
	}
	newEmployeeSkillFile, err := s.ggSheetService.HandleFileCandidateOffer(submittedLink)
	if err != nil {
		return fmt.Errorf("failed to handle file candidate offer: %w", err)
	}
	responseMsg2 := fmt.Sprintf("File skill: %s", newEmployeeSkillFile.SpreadsheetUrl)

	responseBlocks2 := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", responseMsg2, false, false),
			nil,
			nil,
		),
	}
	_, _, err = s.slackClient.PostMessage(payload.Channel.ID, slack.MsgOptionBlocks(responseBlocks2...))
	if err != nil {
		return fmt.Errorf("failed to send confirmation message: %w", err)
	}
	return nil
}
