package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
)

type AIChatbotService struct {
	azureOpenAIConfig config.AzureOpenAIConfig
}

func NewAIChatbotService(azureOpenAIConfig config.AzureOpenAIConfig) *AIChatbotService {
	return &AIChatbotService{azureOpenAIConfig: azureOpenAIConfig}
}

func (s *AIChatbotService) CreateThread(ctx context.Context) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.azureOpenAIConfig.Endpoint+"/openai/threads", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	thread := struct {
		ID string `json:"id"`
	}{}
	err = json.Unmarshal(body, &thread)
	if err != nil {
		return "", err
	}
	return thread.ID, nil
}

func (s *AIChatbotService) CreateMessage(ctx context.Context, threadID string, message string) (string, error) {
	client := &http.Client{}
	requestBody := struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: message,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		s.azureOpenAIConfig.Endpoint+"/openai/threads/"+threadID+"/messages",
		bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	newMessage := struct {
		ID string `json:"id"`
	}{}
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		return "", err
	}
	return newMessage.ID, nil
}

func (s *AIChatbotService) CreateRun(ctx context.Context, threadID string, assistantID string, instructions string) (string, error) {
	client := &http.Client{}
	requestBody := struct {
		AssistantID  string `json:"assistant_id"`
		Instructions string `json:"instructions"`
	}{
		AssistantID:  assistantID,
		Instructions: instructions,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", s.azureOpenAIConfig.Endpoint+"/openai/threads/"+threadID+"/runs", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	run := struct {
		ID string `json:"id"`
	}{}
	err = json.Unmarshal(body, &run)
	if err != nil {
		return "", err
	}
	return run.ID, nil
}

func (s *AIChatbotService) GetRun(ctx context.Context, threadID string, runID string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.azureOpenAIConfig.Endpoint+"/openai/threads/"+threadID+"/runs/"+runID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	run := struct {
		Status string `json:"status"`
	}{}
	err = json.Unmarshal(body, &run)
	if err != nil {
		return "", err
	}
	return run.Status, nil
}

func (s *AIChatbotService) ListMessages(ctx context.Context, threadID string) ([]dto.AzureAIChatbotMessage, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.azureOpenAIConfig.Endpoint+"/openai/threads/"+threadID+"/messages", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	messages := []dto.AzureAIChatbotMessage{}
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *AIChatbotService) GetFileContent(ctx context.Context, fileID string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.azureOpenAIConfig.Endpoint+"/openai/files/"+fileID+"/content", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *AIChatbotService) GetFileInformation(ctx context.Context, fileID string) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"%s/openai/files/%s?api-version=%s",
		s.azureOpenAIConfig.Endpoint,
		fileID,
		s.azureOpenAIConfig.ApiVersion,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("api-key", s.azureOpenAIConfig.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var fileRes map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&fileRes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return fileRes, nil
}

// func (s *AIChatbotService) GetFileStorageUrl(ctx context.Context, fileId string) (string, error) {
// 	fileInfo, err := s.GetFileInformation(ctx, fileId)
// 	if err != nil {
// 		return "", err
// 	}

// 	fileContent, err := s.GetFileContent(ctx, fileId)
// 	if err != nil {
// 		return "", err
// 	}
// }
