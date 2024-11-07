package slack_handlers

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/util"
)

func (s *SlackHandler) handleCreateIntegrateTrainingSubmission(payload slack.InteractionCallback) error {
	fmt.Println("handleCreateIntegrateTrainingSubmission----")
	submittedSheetURL := payload.BlockActionState.Values["sheet_url"]["sheet_url_input"].Value
	submittedSheetName := payload.BlockActionState.Values["sheet_name"]["sheet_name_input"].Value
	if !util.IsValidGoogleSheetLink(submittedSheetURL) {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid skill file link")
		return err
	}
	if submittedSheetName == "" {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Sheet name is required")
		return err
	}
	err := s.uiPathJobService.CreateIntegrateTrainingJob(dto.UIPathCreateIntegrateTrainingInput{
		SheetURL:  submittedSheetURL,
		SheetName: submittedSheetName,
	}, payload.Channel.ID)
	fmt.Println("err-=-=-=-=-=-=", err)
	return err
}

func (s *SlackHandler) handleIntegrateTrainingEvent(channelID string) error {
	return s.slackService.SendIntegrateTrainingForm(context.Background(), channelID)
}
