package services

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
)

type MessageService struct {
	messageRepo repository.MessageRepositoryInterface
}

type MessageServiceInterface interface {
	CreateMessage(message *models.Message) error
	GetMessagesByThreadID(threadID string) ([]models.Message, error)
}

func NewMessageService(messageRepo repository.MessageRepositoryInterface) *MessageService {
	return &MessageService{messageRepo}
}

func (m *MessageService) CreateMessage(message *models.Message) error {
	return m.messageRepo.CreateMessage(message)
}

func (m *MessageService) GetMessagesByThreadID(threadID string) ([]models.Message, error) {
	return m.messageRepo.GetMessagesByThreadID(threadID)
}
