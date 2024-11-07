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
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	GetSigningSecret() string
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

func (s *SlackService) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return s.slackClient.PostMessage(channelID, options...)
}

func (s *SlackService) GetSigningSecret() string {
	return s.slackConfig.SigningSecret
}

func (s *SlackService) SendCandidateFileForm(ctx context.Context, channelID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the candidate file link (google sheet)", false, false),
			nil,
			nil,
		),
		slack.NewInputBlock(
			"candidate_file",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Candidate File",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the candidate file link (google sheet)",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "candidate_file_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the candidate file link (google sheet)"},
			},
		),
		slack.NewActionBlock(
			"submit_candidate_file",
			slack.NewButtonBlockElement(
				"submit_candidate_file",
				"submit_candidate_file",
				slack.NewTextBlockObject("plain_text", "Submit", false, false),
			),
		),
	}

	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (s *SlackService) SendConfirmCloseThread(ctx context.Context, channelID string, threadID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn",
				"We will close the thread after 1 minute, if you want to continue the conversation, please click the button below",
				false,
				false,
			),
			nil,
			nil,
		),
		slack.NewActionBlock(
			"confirm_continue_thread",
			slack.NewButtonBlockElement(
				"continue_thread",
				"continue_thread",
				slack.NewTextBlockObject("plain_text", "Continue Thread", false, false),
			),
		),
	}
	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (s *SlackService) SendWelcomeNewEmployeeForm(ctx context.Context, channelID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the candidate file link (google sheet)", false, false),
			nil,
			nil,
		),
		slack.NewInputBlock(
			"skill_file",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Skill File",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the skill file link (google sheet)",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "skill_file_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the candidate file link (google sheet)"},
				MinLength:   5,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewInputBlock(
			"personal_email",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Personal Email",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the personal email",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "personal_email_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the personal email"},
				MinLength:   5,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewActionBlock(
			"submit_welcome_new_employee",
			slack.NewButtonBlockElement(
				"submit_welcome_new_employee",
				"submit_welcome_new_employee",
				slack.NewTextBlockObject("plain_text", "Submit", false, false),
			),
		),
	}

	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}
