package models

import "time"

const (
	ThreadStatusOpen   = "open"
	ThreadStatusClosed = "closed"
)

type Thread struct {
	ID          string    `json:"id"`
	Messages    []Message `json:"messages"`
	ChannelId   string    `json:"channel_id"`
	SlackUserId string    `json:"slack_user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Status      string    `json:"status"`
}
