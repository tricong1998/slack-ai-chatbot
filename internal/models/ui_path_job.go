package models

import (
	"encoding/json"
	"time"
)

type UIPathJob struct {
	JobID        int             `json:"jobId" gorm:"column:job_id;unique;primaryKey"`
	State        string          `json:"state"`
	Error        string          `json:"error" gorm:"column:error;null"`
	Output       string          `json:"output" gorm:"column:output;null"`
	SlackChannel string          `json:"slackChannel" gorm:"column:slack_channel;not null"`
	JobType      string          `json:"jobType" gorm:"column:job_type;not null"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    time.Time       `json:"deleted_at"`
	Input        json.RawMessage `json:"input" gorm:"column:input;null"`
}

const (
	JobTypeGreeting              = "welcome_new_employee"
	JobTypeFillBuddyForm         = "fill_buddy_form"
	JobTypeIntegrateTrainingForm = "integrate_training_form"
	JobTypeCreateLeaveRequest    = "create_leave_request"
	JobTypePreOnboardEmail       = "pre_onboard_email"
)
