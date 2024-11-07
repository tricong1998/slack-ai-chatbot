package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/rabbitmq"
)

type UIPathJobService struct {
	uiPathJobRepository   *repository.UIPathJobRepository
	UIPathService         *UIPathService
	slackService          *SlackService
	pollingCheckPublisher rabbitmq.IPublisher
}

func NewUIPathJobService(uiPathJobRepository *repository.UIPathJobRepository, pollingCheckPublisher rabbitmq.IPublisher, uiPathService *UIPathService, slackService *SlackService) *UIPathJobService {
	return &UIPathJobService{
		uiPathJobRepository:   uiPathJobRepository,
		pollingCheckPublisher: pollingCheckPublisher,
		UIPathService:         uiPathService,
		slackService:          slackService,
	}
}

func (s *UIPathJobService) CreateJob(job *models.UIPathJob) error {
	return s.uiPathJobRepository.CreateJob(job)
}

func (s *UIPathJobService) PollingCheck(jobID int) (bool, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		job, err := s.GetJob(jobID)
		if err != nil {
			return false, err
		}
		switch job.JobType {
		case models.JobTypeGreeting:
			completed, err := s.HandleCheckGreetingJobPolling(job)
			if err != nil {
				return false, err
			}
			if completed {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *UIPathJobService) HandleCheckGreetingJobPolling(job *models.UIPathJob) (bool, error) {
	status, output, err := s.CheckAndUpdateJobStatus(job.JobID)
	if err != nil {
		str := "Sorry, something went wrong. Please try again later."
		s.slackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, err
	}
	if status == JobStatusCompleted {
		uiPathGreetingOutput := dto.UIPathGreetingOutput{}
		err := json.Unmarshal([]byte(output), &uiPathGreetingOutput)
		fmt.Println("err----", err)
		if err != nil {
			str := "Sorry, something went wrong. Please try again later."
			s.slackService.SendMessage(context.Background(), &job.SlackChannel, str)
			return true, err
		}
		fmt.Println("uiPathGreetingOutput.Greeting----", uiPathGreetingOutput.Greeting)
		fmt.Println("job.SlackChannel----", job.SlackChannel)
		s.slackService.SendMessage(context.Background(), &job.SlackChannel, uiPathGreetingOutput.Greeting)
		return true, nil
	} else if status == JobStatusFailed {
		str := "Sorry, something went wrong. Please try again later."
		s.slackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, nil
	}
	return false, nil
}

func (s *UIPathJobService) CheckAndUpdateJobStatus(jobID int) (string, string, error) {
	fmt.Println("CheckAndUpdateJobStatus----", jobID)
	job, err := s.GetJob(jobID)
	fmt.Println("job----", job)
	fmt.Println("err----", err)
	if err != nil {
		return "", "", err
	}
	jobDetails, err := s.UIPathService.GetJobDetails(job.JobID)
	fmt.Println("jobDetails----", jobDetails)
	fmt.Println("jobDetailsState----", jobDetails.State)
	fmt.Println("err----", err)
	if err != nil {
		job.State = JobStatusFailed
		job.Error = err.Error()
		s.UpdateJob(job)
		return JobStatusFailed, "", err
	}
	if jobDetails.State == JobStatusCompleted {
		job.State = JobStatusCompleted
		job.Output = jobDetails.OutputArguments
		s.UpdateJob(job)
		return JobStatusCompleted, jobDetails.OutputArguments, nil
	} else if jobDetails.State == JobStatusFailed {
		job.State = JobStatusFailed
		s.UpdateJob(job)
		return JobStatusFailed, "", nil
	}
	return JobStatusPending, "", nil
}

func (s *UIPathJobService) GetJob(jobID int) (*models.UIPathJob, error) {
	return s.uiPathJobRepository.GetJob(jobID)
}

func (s *UIPathJobService) UpdateJob(job *models.UIPathJob) error {
	return s.uiPathJobRepository.UpdateJob(job)
}

func (s *UIPathJobService) CreateGreetingJob(input dto.UIPathGreetingNewEmployee, slackChannel string) error {
	uiJob, err := s.UIPathService.GreetingNewEmployee(input)
	if err != nil {
		return err
	}
	job := &models.UIPathJob{
		JobID:        uiJob.ID,
		JobType:      models.JobTypeGreeting,
		SlackChannel: slackChannel,
	}
	err = s.uiPathJobRepository.CreateJob(job)
	if err != nil {
		return err
	}
	s.pollingCheckPublisher.PublishMessage(dto.UIPathGreetingJobInput{JobID: job.JobID})
	return nil
}
