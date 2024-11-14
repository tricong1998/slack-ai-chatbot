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
	SlackService          *SlackService
	pollingCheckPublisher rabbitmq.IPublisher
}

func NewUIPathJobService(uiPathJobRepository *repository.UIPathJobRepository, pollingCheckPublisher rabbitmq.IPublisher, uiPathService *UIPathService, slackService *SlackService) *UIPathJobService {
	return &UIPathJobService{
		uiPathJobRepository:   uiPathJobRepository,
		pollingCheckPublisher: pollingCheckPublisher,
		UIPathService:         uiPathService,
		SlackService:          slackService,
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
			s.SlackService.SendMessage(context.Background(), &job.SlackChannel, err.Error())
			return false, err
		}
		switch job.JobType {
		case models.JobTypeGreeting:
			completed, err := s.HandleCheckGreetingJobPolling(job)
			if err != nil {
				s.SlackService.SendMessage(context.Background(), &job.SlackChannel, err.Error())
				return false, err
			}
			if completed {
				return true, nil
			}
		case models.JobTypeFillBuddyForm:
			completed, err := s.HandleCheckFillBuddyFormJobPolling(job)
			if err != nil {
				s.SlackService.SendMessage(context.Background(), &job.SlackChannel, err.Error())
				return false, err
			}
			if completed {
				return true, nil
			}
		case models.JobTypeCreateLeaveRequest:
			completed, err := s.HandleCheckCreateLeaveRequestJobPolling(job)
			if err != nil {
				s.SlackService.SendMessage(context.Background(), &job.SlackChannel, err.Error())
				return false, err
			}
			if completed {
				return true, nil
			}
		case models.JobTypeIntegrateTrainingForm:
			completed, err := s.HandleCheckCreateIntegrateTrainingJobPolling(job)
			if err != nil {
				s.SlackService.SendMessage(context.Background(), &job.SlackChannel, err.Error())
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
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, err
	}
	if status == JobStatusCompleted {
		uiPathGreetingOutput := dto.UIPathGreetingOutput{}
		err := json.Unmarshal([]byte(output), &uiPathGreetingOutput)
		if err != nil {
			str := "Sorry, something went wrong. Please try again later."
			s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
			return true, err
		}
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, uiPathGreetingOutput.Greeting)
		return true, nil
	} else if status == JobStatusFailed {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, nil
	}
	return false, nil
}

func (s *UIPathJobService) HandleCheckFillBuddyFormJobPolling(job *models.UIPathJob) (bool, error) {
	status, output, err := s.CheckAndUpdateJobStatus(job.JobID)
	if err != nil {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, err
	}
	if status == JobStatusCompleted {
		uiPathFillBuddyOutput := dto.UIPathFillBuddyOutput{}
		err := json.Unmarshal([]byte(output), &uiPathFillBuddyOutput)
		if err != nil {
			str := "Sorry, something went wrong. Please try again later."
			s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
			return true, err
		}
		response := fmt.Sprintf("Buddy form created successfully. Please check file %s", uiPathFillBuddyOutput.BuddyFormName)
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, response)
		return true, nil
	} else if status == JobStatusFailed {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, nil
	}
	return false, nil
}

func (s *UIPathJobService) HandleCheckCreateLeaveRequestJobPolling(job *models.UIPathJob) (bool, error) {
	status, output, err := s.CheckAndUpdateJobStatus(job.JobID)
	if err != nil {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, err
	}
	if status == JobStatusCompleted {
		uiPathCreateLeaveRequestOutput := dto.UIPathLeaveOutput{}
		err := json.Unmarshal([]byte(output), &uiPathCreateLeaveRequestOutput)
		if err != nil {
			return true, err
		}
		var uiPathCreateLeaveRequestOutputResponse dto.UIPathLeaveOutputResponse
		err = json.Unmarshal([]byte(uiPathCreateLeaveRequestOutput.Response), &uiPathCreateLeaveRequestOutputResponse)
		if err != nil {
			return true, err
		}
		if uiPathCreateLeaveRequestOutputResponse.Result != nil {
			if uiPathCreateLeaveRequestOutputResponse.Result.Code == 200 {
				str := "Leave request created successfully. Please check your calendar."
				s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
				return true, nil
			}
		}
		if uiPathCreateLeaveRequestOutputResponse.Error != nil {
			err = fmt.Errorf(uiPathCreateLeaveRequestOutputResponse.Error.Data.Message)
			return true, err
		}
		err = fmt.Errorf("sorry, something went wrong - please try again later")
		return true, err
	} else if status == JobStatusFailed {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, nil
	}
	return false, nil
}

func (s *UIPathJobService) HandleCheckCreateIntegrateTrainingJobPolling(job *models.UIPathJob) (bool, error) {
	status, output, err := s.CheckAndUpdateJobStatus(job.JobID)
	if err != nil {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, err
	}
	if status == JobStatusCompleted {
		uiPathCreateIntegrateTrainingOutput := dto.UIPathCreateIntegrateTrainingOutput{}
		err := json.Unmarshal([]byte(output), &uiPathCreateIntegrateTrainingOutput)
		if err != nil {
			return true, err
		}
		if uiPathCreateIntegrateTrainingOutput.CalendarId != "" {
			s.SlackService.SendMessage(context.Background(), &job.SlackChannel, uiPathCreateIntegrateTrainingOutput.CalendarId)
			return true, nil
		}
		err = fmt.Errorf(uiPathCreateIntegrateTrainingOutput.ErrMessage)
		return true, err
	} else if status == JobStatusFailed {
		str := "Sorry, something went wrong. Please try again later."
		s.SlackService.SendMessage(context.Background(), &job.SlackChannel, str)
		return true, nil
	}
	return false, nil
}

func (s *UIPathJobService) CheckAndUpdateJobStatus(jobID int) (string, string, error) {
	job, err := s.GetJob(jobID)
	if err != nil {
		return "", "", err
	}
	jobDetails, err := s.UIPathService.GetJobDetails(job.JobID)
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
	s.pollingCheckPublisher.PublishMessage(dto.UIPathCheckingJobInput{JobID: job.JobID})
	return nil
}

func (s *UIPathJobService) CreateFillBuddyJob(input dto.UIPathFillBuddyInput, slackChannel string) error {
	uiJob, err := s.UIPathService.FillBuddyForm(input)
	if err != nil {
		return err
	}
	job := &models.UIPathJob{
		JobID:        uiJob.ID,
		JobType:      models.JobTypeFillBuddyForm,
		SlackChannel: slackChannel,
	}
	err = s.uiPathJobRepository.CreateJob(job)
	if err != nil {
		return err
	}
	s.pollingCheckPublisher.PublishMessage(dto.UIPathCheckingJobInput{JobID: job.JobID})
	return nil
}

func (s *UIPathJobService) CreateIntegrateTrainingJob(input dto.UIPathCreateIntegrateTrainingInput, slackChannel string) error {
	uiJob, err := s.UIPathService.CreateIntegrateTraining(input)
	if err != nil {
		return err
	}
	job := &models.UIPathJob{
		JobID:        uiJob.ID,
		JobType:      models.JobTypeIntegrateTrainingForm,
		SlackChannel: slackChannel,
	}
	err = s.uiPathJobRepository.CreateJob(job)
	if err != nil {
		return err
	}
	s.pollingCheckPublisher.PublishMessage(dto.UIPathCheckingJobInput{JobID: job.JobID})
	return nil
}

func (s *UIPathJobService) CreateLeaveRequestJob(input dto.UIPathCreateLeaveRequestInput, slackChannel string) error {
	uiJob, err := s.UIPathService.CreateLeaveRequestOnOdoo(input)
	if err != nil {
		return err
	}
	job := &models.UIPathJob{
		JobID:        uiJob.ID,
		JobType:      models.JobTypeCreateLeaveRequest,
		SlackChannel: slackChannel,
	}
	err = s.uiPathJobRepository.CreateJob(job)
	if err != nil {
		return err
	}
	s.pollingCheckPublisher.PublishMessage(dto.UIPathCheckingJobInput{JobID: job.JobID})
	return nil
}
