package repository

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"gorm.io/gorm"
)

type ThreadRepository struct {
	db *gorm.DB
}

type ThreadRepositoryInterface interface {
	CreateThread(thread *models.Thread) error
	GetThreadByID(threadID string) (*models.Thread, error)
	GetLatestOpenThreadByChannelAndUserID(channelID string, userID string) (*models.Thread, error)
	UpdateThreadStatus(threadID string, status string) error
}

func NewThreadRepository(db *gorm.DB) *ThreadRepository {
	return &ThreadRepository{db}
}

func (t *ThreadRepository) CreateThread(thread *models.Thread) error {
	thread.Status = models.ThreadStatusOpen
	return t.db.Create(thread).Error
}

func (t *ThreadRepository) GetThreadByID(threadID string) (*models.Thread, error) {
	var thread models.Thread
	return &thread, t.db.Where("id = ?", threadID).First(&thread).Error
}

func (t *ThreadRepository) GetLatestOpenThreadByChannelAndUserID(channelID string, userID string) (*models.Thread, error) {
	var thread models.Thread
	return &thread, t.db.Where("channel_id = ? AND slack_user_id = ? AND status = ?", channelID, userID, models.ThreadStatusOpen).Order("created_at DESC").First(&thread).Error
}

func (t *ThreadRepository) UpdateThreadStatus(threadID string, status string) error {
	return t.db.Model(&models.Thread{}).Where("id = ?", threadID).Update("status", status).Error
}
