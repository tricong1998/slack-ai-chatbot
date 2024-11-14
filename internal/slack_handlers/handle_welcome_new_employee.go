package slack_handlers

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/util"
)

func (s *SlackHandler) handleGreetingNewEmployeeSubmission(payload slack.InteractionCallback) error {
	submittedSkillFile := payload.BlockActionState.Values["skill_file"]["skill_file_input"].Value
	submittedPersonalEmail := payload.BlockActionState.Values["personal_email"]["personal_email_input"].Value
	if !util.IsValidGoogleSheetLink(submittedSkillFile) {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid skill file link")
		return err
	}
	if !util.IsValidEmail(submittedPersonalEmail) {
		err := s.slackService.SendMessage(context.Background(), &payload.Channel.ID, "Invalid personal email")
		return err
	}
	err := s.uiPathJobService.CreateGreetingJob(dto.UIPathGreetingNewEmployee{
		SkillFile:     submittedSkillFile,
		PersonalEmail: submittedPersonalEmail,
	}, payload.Channel.ID)
	return err
}

func (s *SlackHandler) handleGreetingNewEmployeeEvent(channelID string) error {
	return s.slackService.SendWelcomeNewEmployeeForm(context.Background(), channelID)
}
