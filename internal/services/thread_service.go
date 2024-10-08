package services

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
)

type ThreadService struct {
	threadRepo repository.ThreadRepositoryInterface
}

type ThreadServiceInterface interface {
	CreateThread(thread *models.Thread) error
	GetThreadByID(threadID string) (*models.Thread, error)
	CloseThreadStatus(threadID string) error
	GetLatestOpenThreadByChannelAndUserID(channelID string, userID string) (*models.Thread, error)
}

func NewThreadService(threadRepo repository.ThreadRepositoryInterface) *ThreadService {
	return &ThreadService{threadRepo}
}

func (t *ThreadService) CreateThread(thread *models.Thread) error {
	return t.threadRepo.CreateThread(thread)
}

func (t *ThreadService) GetThreadByID(threadID string) (*models.Thread, error) {
	return t.threadRepo.GetThreadByID(threadID)
}

func (t *ThreadService) GetLatestOpenThreadByChannelAndUserID(channelID string, userID string) (*models.Thread, error) {
	return t.threadRepo.GetLatestOpenThreadByChannelAndUserID(channelID, userID)
}

func (t *ThreadService) CloseThreadStatus(threadID string) error {
	return t.threadRepo.UpdateThreadStatus(threadID, models.ThreadStatusClosed)
}
