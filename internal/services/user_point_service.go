package services

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
)

type UserPointService struct {
	UserPointRepo repository.IUserPointRepository
}

type IUserPointService interface {
	CreateUserPoint(input *models.UserPoint) error
	ReadUserPoint(id uint) (*models.UserPoint, error)
	ListUserPoints(
		perPage, page int32,
		userId *uint,
	) ([]models.UserPoint, int64, error)
	UpdateUserPoint(user *models.UserPoint) error
	DeleteUserPoint(id uint) error
}

func NewUserPointService(userRepo repository.IUserPointRepository) *UserPointService {
	return &UserPointService{userRepo}
}

func (us *UserPointService) CreateUserPoint(user *models.UserPoint) error {

	err := us.UserPointRepo.CreateUserPoint(user)
	return err
}

func (us *UserPointService) ReadUserPoint(id uint) (*models.UserPoint, error) {
	user, err := us.UserPointRepo.ReadUserPoint(id)
	return user, err
}

func (us *UserPointService) ListUserPoints(
	perPage, page int32,
	userId *uint,
) ([]models.UserPoint, int64, error) {
	return us.UserPointRepo.ListUserPoints(perPage, page, userId)
}

func (us *UserPointService) UpdateUserPoint(user *models.UserPoint) error {
	err := us.UserPointRepo.UpdateUserPoint(user)
	return err
}

func (us *UserPointService) DeleteUserPoint(id uint) error {
	return us.UserPointRepo.DeleteUserPoint(id)
}
