package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
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
				MinLength:   1,   // Minimum length for a valid email
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
				MinLength:   1,   // Minimum length for a valid email
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

func (s *SlackService) SendCreateBuddyForm(ctx context.Context, channelID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the transformation input and output file link (google sheet)", false, false),
			nil,
			nil,
		),
		slack.NewInputBlock(
			"transformation_input_file",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Transformation Input File",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the transformation input file link (google sheet)",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "transformation_input_file_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the transformation input file link (google sheet)"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewInputBlock(
			"transformation_output_file",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Transformation Output File",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the transformation output file link (google sheet)",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "transformation_output_file_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the transformation output file link (google sheet)"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewActionBlock(
			"submit_create_buddy",
			slack.NewButtonBlockElement(
				"submit_create_buddy",
				"submit_create_buddy",
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

func (s *SlackService) SendIntegrateTrainingForm(ctx context.Context, channelID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the sheet url and sheet name", false, false),
			nil,
			nil,
		),
		slack.NewInputBlock(
			"sheet_url",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Sheet URL",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the sheet url",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "sheet_url_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the sheet url"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewInputBlock(
			"sheet_name",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Sheet Name",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the sheet name",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "sheet_name_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the sheet name"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewActionBlock(
			"submit_integrate_training",
			slack.NewButtonBlockElement(
				"submit_integrate_training",
				"submit_integrate_training",
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

func (s *SlackService) SendCreateLeaveRequestForm(ctx context.Context, channelID string) error {
	leaveOptions := make([]*slack.OptionBlockObject, 0)
	for _, leave := range dto.AppMappingCodeLeave {
		leaveOptions = append(leaveOptions, &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject(slack.PlainTextType, leave.Name, false, false),
			Value: strconv.Itoa(leave.Code),
		})
	}
	workingTimeOptions := make([]*slack.OptionBlockObject, 0)
	for _, workingTime := range dto.AppMappingCodeWorkingTime {
		workingTimeOptions = append(workingTimeOptions, &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject(slack.PlainTextType, workingTime.Name, false, false),
			Value: strconv.Itoa(workingTime.Code),
		})
	}
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the leave request information", false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.MarkdownType, "*Leave Type*", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "*Working Time*", false, false),
			},
			nil,
		),
		slack.NewActionBlock(
			"leave_type",
			&slack.SelectBlockElement{
				Type:        slack.OptTypeStatic,
				ActionID:    "leave_type_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select leave type", Emoji: false},
				Options:     leaveOptions,
			},
			&slack.SelectBlockElement{
				Type:        slack.OptTypeStatic,
				ActionID:    "working_time_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select working time", Emoji: false},
				Options:     workingTimeOptions,
			},
		),
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.MarkdownType, "*Request Date From*", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "*Request Date To*", false, false),
			},
			nil,
		),
		slack.NewActionBlock(
			"date_pickers",
			&slack.DatePickerBlockElement{
				Type:        slack.METDatepicker,
				ActionID:    "request_date_from_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select start date"},
			},
			&slack.DatePickerBlockElement{
				Type:        slack.METDatepicker,
				ActionID:    "request_date_to_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select end date"},
			},
		),
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject(slack.MarkdownType, "*Hour From*", false, false),
				slack.NewTextBlockObject(slack.MarkdownType, "*Hour To*", false, false),
			},
			nil,
		),
		slack.NewActionBlock(
			"time_pickers",
			&slack.TimePickerBlockElement{
				Type:        slack.METTimepicker,
				ActionID:    "hour_from_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select start time"},
			},
			&slack.TimePickerBlockElement{ // Changed from DatePickerBlockElement to TimePickerBlockElement
				Type:        slack.METTimepicker,
				ActionID:    "hour_to_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Select end time"}, // Fixed text from "end date" to "end time"
			},
		),
		slack.NewInputBlock(
			"description",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Description",
				Emoji:    false,
				Verbatim: false,
			},
			nil,
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "description_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the description"},
				MinLength:   0,
				MaxLength:   254,
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewActionBlock(
			"submit_create_leave_request",
			slack.NewButtonBlockElement(
				"submit_create_leave_request",
				"submit_create_leave_request",
				slack.NewTextBlockObject("plain_text", "Submit", false, false),
			),
		),
	}

	_, _, err := s.slackClient.PostMessage(channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		s.slackClient.PostMessage(channelID, slack.MsgOptionText("Failed to send create leave request form: "+err.Error(), false))
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (s *SlackService) SendPreOnboardEmailForm(ctx context.Context, channelID string) error {
	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Please enter the sheet url and sheet name", false, false),
			nil,
			nil,
		),
		slack.NewInputBlock(
			"sheet_url",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Sheet URL",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the sheet url",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "sheet_url_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the sheet url"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewInputBlock(
			"sheet_name",
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Sheet Name",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.TextBlockObject{
				Type:     slack.PlainTextType,
				Text:     "Enter the sheet name",
				Emoji:    false,
				Verbatim: false,
			},
			&slack.PlainTextInputBlockElement{
				Type:        slack.METPlainTextInput,
				ActionID:    "sheet_name_input",
				Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter the sheet name"},
				MinLength:   1,   // Minimum length for a valid email
				MaxLength:   254, // Maximum length per RFC 5321
				DispatchActionConfig: &slack.DispatchActionConfig{
					TriggerActionsOn: []string{"on_enter_pressed"},
				},
			},
		),
		slack.NewActionBlock(
			"submit_pre_onboard_email",
			slack.NewButtonBlockElement(
				"submit_pre_onboard_email",
				"submit_pre_onboard_email",
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
