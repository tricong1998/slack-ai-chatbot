package dto

type SendMessageDto struct {
	Message string `json:"message" binding:"required"`
}
