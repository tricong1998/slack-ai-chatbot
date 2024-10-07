package slack_handlers

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleEventMessage will take an event and handle it properly based on the type of event
func (s *SlackHandler) HandleEventMessage(event slackevents.EventsAPIEvent) error {
	fmt.Println("asdfjaldkfjlkadjsflkds", event.Data)
	fmt.Println("asdfjaldkfjlkadjsflkds2", event.Type)
	fmt.Printf("asdfjaldkfjlkadjsflkds3: %+v\n", event)
	fmt.Printf("asdfjaldkfjlkadjsflkds4: %+v\n", reflect.TypeOf(event.InnerEvent.Data))
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		fmt.Println("ffff", innerEvent)
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			return s.handleAppMentionEvent(ev)
		case *slackevents.MessageEvent:
			return s.handleMessageEvent(ev)
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func (s *SlackHandler) handleAppMentionEvent(event *slackevents.AppMentionEvent) error {
	user, err := s.slackClient.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	text := strings.ToLower(event.Text)
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") {
		// Greet the user
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings"
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "How can I be of service"
		attachment.Color = "#3d3d3d"
	}
	_, _, err = s.slackClient.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func (s *SlackHandler) handleMessageEvent(event *slackevents.MessageEvent) error {
	// user, err := s.slackClient.GetUserInfo(event.User)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("asdfjaldkfjlkadjsflkds", event.Text)
	s.aiChatbotService.AddAndRunMessage(context.Background(), &event.Channel, event.Text)
	return nil
}
