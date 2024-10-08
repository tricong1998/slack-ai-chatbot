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
		return handleWasChatbotUsefulCommand(command, client)
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
