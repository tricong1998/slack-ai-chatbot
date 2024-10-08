package slack_handlers

import (
	"fmt"

	"github.com/slack-go/slack"
)

func handleEventInteraction(event *slack.InteractionCallback, client *slack.Client) error {
	fmt.Printf("Interaction event: %v", event.ActionID)
	fmt.Printf("Interaction event type: %s", event.Type)

	switch event.Type {
	case slack.InteractionTypeBlockActions:
		fmt.Printf("Interaction event type: %s", event.Type)
		for _, action := range event.ActionCallback.BlockActions {
			fmt.Printf("%+v", action)
			fmt.Printf("%+v", action.SelectedOption)
		}
	}
	return nil
}
