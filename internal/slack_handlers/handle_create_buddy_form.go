package slack_handlers

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/util"
)

func (s *SlackHandler) handleCreateBuddyFormFileEvent(payload slack.InteractionCallback) error {
	submittedTransformationInputFile := payload.BlockActionState.Values["transformation_input_file"]["transformation_input_file_input"].Value
	submittedTransformationOutputFile := payload.BlockActionState.Values["transformation_output_file"]["transformation_output_file_input"].Value
	if !util.IsValidGoogleSheetLink(submittedTransformationInputFile) {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid transformation input file link")
		return err
	}
	if !util.IsValidGoogleSheetLink(submittedTransformationOutputFile) {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid transformation output file link")
		return err
	}
	err := s.uiPathJobService.CreateFillBuddyJob(dto.UIPathFillBuddyInput{
		InputSheet:  submittedTransformationInputFile,
		OutputSheet: submittedTransformationOutputFile,
	}, payload.Channel.ID)
	return err
}

func (s *SlackHandler) handleCreateBuddyFormEvent(channelID string) error {
	return s.slackService.SendCreateBuddyForm(context.Background(), channelID)
}
