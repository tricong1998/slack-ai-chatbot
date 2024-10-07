package dto

type AddMessageRequest struct {
	ChannelID *string `json:"channel_id"`
	Message   string  `json:"message"`
}
