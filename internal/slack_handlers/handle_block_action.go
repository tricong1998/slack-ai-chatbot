package slack_handlers

import (
	"github.com/slack-go/slack"
)

func (s *SlackHandler) HandleBlockAction(payload slack.InteractionCallback) (string, error) {
	for _, action := range payload.ActionCallback.BlockActions {
		switch action.ActionID {
		case "submit_candidate_file":
			return "", s.handleCandidateSheetSubmission(payload)
		case "submit_welcome_new_employee":
			return "", s.handleGreetingNewEmployeeSubmission(payload)
		case "submit_create_buddy":
			return "", s.handleCreateBuddyFormFileEvent(payload)
		case "submit_create_leave_request":
			return "", s.handleCreateLeaveRequestSubmission(payload)
		case "submit_integrate_training":
			return "", s.handleCreateIntegrateTrainingSubmission(payload)
			// ... handle other action IDs as needed ...
		}
	}
	return "", nil
}
