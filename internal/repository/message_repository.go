package repository

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

type MessageRepositoryInterface interface {
	CreateMessage(message *models.Message) error
	GetMessagesByThreadID(threadID string) ([]models.Message, error)
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db}
}

func (m *MessageRepository) CreateMessage(message *models.Message) error {
	return m.db.Create(message).Error
}

func (m *MessageRepository) GetMessagesByThreadID(threadID string) ([]models.Message, error) {
	var messages []models.Message
	return messages, m.db.Where("thread_id = ?", threadID).Find(&messages).Error
}
