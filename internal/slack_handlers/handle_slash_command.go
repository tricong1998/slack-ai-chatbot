package slack_handlers

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

func (s *SlackHandler) HandleSlashCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// We need to switch depending on the command
	switch command.Command {
	case "/hello":
		// This was a hello command, so pass it along to the proper function
		return handleHelloCommand(command, client), nil

	case "/was-chatbot-useful":
		return handleWasChatbotUsefulCommandWithForm(command, client)
	}

	return nil, nil
}

// handleHelloCommand will take care of /hello submissions
func handleHelloCommand(command slack.SlashCommand, client *slack.Client) error {
	// The Input is found in the text field so
	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: command.UserName,
		},
	}

	// Greet the user
	attachment.Text = fmt.Sprintf("Hello %s", command.Text)
	attachment.Color = "#4af030"

	// Send the message to the channel
	// The Channel is available in the command.ChannelID
	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func handleWasChatbotUsefulCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	attachment := slack.Attachment{}
	checkbox := slack.NewCheckboxGroupsBlockElement("answer",
		slack.NewOptionBlockObject(
			"yes",
			&slack.TextBlockObject{Text: "Yes", Type: slack.MarkdownType},
			&slack.TextBlockObject{Text: "Did you Enjoy it?", Type: slack.MarkdownType},
		),
		slack.NewOptionBlockObject(
			"no",
			&slack.TextBlockObject{Text: "No", Type: slack.MarkdownType},
			&slack.TextBlockObject{Text: "Did you Dislike it?", Type: slack.MarkdownType},
		),
	)
	accessory := slack.NewAccessory(checkbox)
	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			// Create a new section block element and add some text and the accessory to it
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Did you think this chatbot was helpful?",
				},
				nil,
				accessory,
			),
		},
	}

	attachment.Text = "Rate the tutorial"
	attachment.Color = "#4af030"
	return attachment, nil
}

func handleWasChatbotUsefulCommandWithButton(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	attachment := slack.Attachment{}
	yesButton := slack.NewButtonBlockElement("yes", "yes", &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Yes"})
	noButton := slack.NewButtonBlockElement("no", "no", &slack.TextBlockObject{Type: slack.PlainTextType, Text: "No"})

	numberOfDays := 3
	date := time.Now().Format("2006-01-02")
	reason := "illness"
	text := fmt.Sprintf("Do you want to take %d days off from %s due to %s?", numberOfDays, date, reason)

	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: text,
				},
				nil,
				nil,
			),
			slack.NewActionBlock(
				"chatbot_rating",
				yesButton,
				noButton,
			),
		},
	}

	attachment.Text = text
	attachment.Color = "#4af030"
	return attachment, nil
}

func handleWasChatbotUsefulCommandWithForm(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// Create a new modal view
	modalView := slack.ModalViewRequest{
		Type: slack.ViewType("modal"),
		Title: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Rate the Chatbot",
		},
		Close: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Cancel",
		},
		Submit: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Submit",
		},
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"rating_block",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Enter the rating",
					},
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Enter a rating from 1 to 5",
					},
					&slack.PlainTextInputBlockElement{
						Type:        slack.METPlainTextInput,
						ActionID:    "rating_input",
						Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Enter a rating from 1 to 5"},
						MaxLength:   1,
					},
				),
			},
		},
	}

	// Open the modal
	_, err := client.OpenView(command.TriggerID, modalView)
	if err != nil {
		return nil, fmt.Errorf("failed to open modal: %w", err)
	}

	// Return a response to acknowledge the command
	return &slack.Msg{
		ResponseType: slack.ResponseTypeEphemeral,
		Text:         "Please rate the chatbot in the modal that just opened.",
	}, nil
}
