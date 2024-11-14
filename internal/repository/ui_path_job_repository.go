package repository

import (
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"gorm.io/gorm"
)

type UIPathJobRepository struct {
	db *gorm.DB
}

func NewUIPathJobRepository(db *gorm.DB) *UIPathJobRepository {
	return &UIPathJobRepository{db: db}
}

func (r *UIPathJobRepository) CreateJob(job *models.UIPathJob) error {
	return r.db.Create(job).Error
}

func (r *UIPathJobRepository) GetJob(jobID int) (*models.UIPathJob, error) {
	var job models.UIPathJob
	if err := r.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *UIPathJobRepository) UpdateJob(job *models.UIPathJob) error {
	return r.db.Save(job).Error
}
