package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
)

const (
	JobStatusPending   = "Pending"
	JobStatusRunning   = "Running"
	JobStatusCompleted = "Successful"
	JobStatusFailed    = "Faulted"
)

type UIPathService struct {
	client *http.Client
	config config.UIPathConfig
}

func NewUIPathService(client *http.Client, config config.UIPathConfig) *UIPathService {
	return &UIPathService{
		client: client,
		config: config,
	}
}

func (s *UIPathService) GetJobDetails(jobID int) (dto.UIPathJobDetails, error) {
	url := s.GetUrlJobDetails(jobID)
	resp, err := s.Call("GET", url, nil)
	if err != nil {
		return dto.UIPathJobDetails{}, err
	}
	defer resp.Body.Close()

	var data dto.UIPathJobDetails
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return dto.UIPathJobDetails{}, err
	}

	return data, nil
}

func (s *UIPathService) GetUrlTrigger(processName string) string {
	return fmt.Sprintf("%s/%s/orchestrator_/t/%s/%s", s.config.Host, s.config.Tenant, s.config.TenantID, processName)
}

func (s *UIPathService) GetUrlJobDetails(jobID int) string {
	return fmt.Sprintf("%s/%s/orchestrator_/odata/Jobs(%d)", s.config.Host, s.config.Tenant, jobID)
}

func (s *UIPathService) Call(method string, path string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.ApiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp, nil
}

func (s *UIPathService) GreetingNewEmployee(body dto.UIPathGreetingNewEmployee) (*dto.UIPathTriggerResponse, error) {
	return s.CallPostTriggerJob(body, s.config.GreetingNewEmployeeProcessKey)
}

func (s *UIPathService) FillBuddyForm(body dto.UIPathFillBuddyInput) (*dto.UIPathTriggerResponse, error) {
	return s.CallPostTriggerJob(body, s.config.FillBuddyProcessKey)
}

func (s *UIPathService) CreateLeaveRequestOnOdoo(body dto.UIPathCreateLeaveRequestInput) (*dto.UIPathTriggerResponse, error) {
	return s.CallPostTriggerJob(body, s.config.CreateLeaveRequestProcessKey)
}

func (s *UIPathService) CreateIntegrateTraining(body dto.UIPathCreateIntegrateTrainingInput) (*dto.UIPathTriggerResponse, error) {
	return s.CallPostTriggerJob(body, s.config.CreateIntegrateTrainingProcessKey)
}

func (s *UIPathService) CallPostTriggerJob(body interface{}, key string) (*dto.UIPathTriggerResponse, error) {
	url := s.GetUrlTrigger(key)
	fmt.Println("url-=-=-=-", url, body)
	resp, err := s.Call("POST", url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		var error dto.UIPathErrorTriggerJob
		err = json.NewDecoder(resp.Body).Decode(&error)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(error.Message)
	}
	var data dto.UIPathTriggerResponse
	fmt.Println(resp)
	fmt.Println(url)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
