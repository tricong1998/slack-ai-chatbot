package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/models"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/util"
	"gorm.io/gorm"
)

type AIChatbotService struct {
	azureOpenAIConfig config.AzureOpenAIConfig
	slackService      *SlackService
	threadService     *ThreadService
	messageService    *MessageService
}

func NewAIChatbotService(azureOpenAIConfig config.AzureOpenAIConfig, slackService *SlackService, threadService *ThreadService, messageService *MessageService) *AIChatbotService {
	return &AIChatbotService{azureOpenAIConfig: azureOpenAIConfig, slackService: slackService, threadService: threadService, messageService: messageService}
}

func (s *AIChatbotService) CreateThread(ctx context.Context) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.getUrl("threads"), nil)
	if err != nil {
		return "", err
	}
	s.addHeader(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
	url := s.getUrl("threads/" + threadID + "/messages")
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	s.addHeader(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

func (s *AIChatbotService) CreateRun(ctx context.Context, threadID string, assistantID string) (string, error) {
	client := &http.Client{}
	requestBody := struct {
		AssistantID string `json:"assistant_id"`
	}{
		AssistantID: assistantID,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	url := s.getUrl("threads/" + threadID + "/runs")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	s.addHeader(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
	url := s.getUrl("threads/" + threadID + "/runs/" + runID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	s.addHeader(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
	url := s.getUrl("threads/" + threadID + "/messages")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addHeader(req, true)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := struct {
		Data   []dto.AzureAIChatbotMessage `json:"data"`
		Object string                      `json:"object"`
	}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (s *AIChatbotService) GetFileContent(ctx context.Context, fileID string) ([]byte, error) {
	client := &http.Client{}
	url := s.getUrl("files/" + fileID + "/content")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addHeader(req, false)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *AIChatbotService) GetFileInformation(ctx context.Context, fileID string) (map[string]interface{}, error) {
	url := s.getUrl("files/" + fileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	s.addHeader(req, false)

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

func (s *AIChatbotService) AddAndRunMessage(ctx context.Context, channelID *string, message string, userID string) (string, string, error) {
	thread, err := s.threadService.GetLatestOpenThreadByChannelAndUserID(*channelID, userID)
	var threadID string
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return "", "", err
		} else {
			threadID, err = s.CreateThread(ctx)
			if err != nil {
				return "", "", err
			}
			s.threadService.CreateThread(&models.Thread{
				ID:          threadID,
				ChannelId:   *channelID,
				SlackUserId: userID,
			})
		}
	} else {
		threadID = thread.ID
	}

	messageID, err := s.CreateMessage(ctx, threadID, message)
	if err != nil {
		return "", "", err
	}
	runID, err := s.CreateRun(ctx, threadID, s.azureOpenAIConfig.AssistantIdDetectAction)
	if err != nil {
		return "", "", err
	}
	var action string
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-ticker.C:
				runStatus, err := s.GetRun(ctx, threadID, runID)
				if err != nil {
					return
				}
				if runStatus != "queued" && runStatus != "in_progress" {
					if runStatus == "completed" {
						listMessages, err := s.ListMessages(ctx, threadID)
						if err != nil {
							done <- true
							return
						}
						consecutiveAssistantMessages := s.GetFirstConsecutiveAssistantMessages(listMessages)
						for _, message := range consecutiveAssistantMessages {
							var textContent *dto.AzureAIChatbotMessageContent
							for _, content := range message.Content {
								if content.Type == "text" {
									textContent = &content
									break
								}
							}
							if textContent == nil {
								continue
							}
							// annotations := textContent.Text.Annotations
							// if len(annotations) > 0 {
							// regexPattern := regexp.MustCompile(`assistant-[a-zA-Z0-9]+`)
							// matches := regexPattern.FindAllString(textContent.Text.Value, -1)
							// for _, match := range matches {
							// }
							// }
							// for _, annotation := range annotations {
							if textContent.Text.Value != "" {
								action = util.DetectAction(textContent.Text.Value)
								s.slackService.SendMessage(ctx, channelID, textContent.Text.Value)
							}
						}
					}
					done <- true
					return
				}
			case <-timeout:
				done <- true
				return
			}
		}
	}()

	<-done
	// go s.SendMessageCloseThreadAfter5Minutes(ctx, *channelID, threadID)
	return messageID, action, nil
}

func (s *AIChatbotService) SendMessageCloseThreadAfter5Minutes(ctx context.Context, channelID string, threadID string) error {
	timer := time.NewTimer(5 * time.Minute)
	// eventChan <- struct{}{}
	eventChan := make(chan struct{})

	for {
		select {
		case <-timer.C:
			// Timer expired, close the thread
			// err := s.threadService.CloseThread(threadID)
			// if err != nil {
			// 	return fmt.Errorf("failed to close thread: %w", err)
			// }
			// message := "This thread has been inactive for 5 minutes and will now be closed."
			// err = s.slackService.SendMessage(ctx, &channelID, message)
			// if err != nil {
			// 	return fmt.Errorf("failed to send closing message: %w", err)
			// }
			return nil

		case <-eventChan:
			// Event received, reset the timer
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(5 * time.Minute)

		case <-ctx.Done():
			// Context cancelled, stop the timer and exit
			timer.Stop()
			return ctx.Err()
		}
	}
}

func (s *AIChatbotService) GetFirstConsecutiveAssistantMessages(messages []dto.AzureAIChatbotMessage) []dto.AzureAIChatbotMessage {
	consecutiveAssistantMessages := []dto.AzureAIChatbotMessage{}
	for _, message := range messages {
		if message.Role == "assistant" {
			consecutiveAssistantMessages = append(consecutiveAssistantMessages, message)
		} else {
			if len(consecutiveAssistantMessages) > 0 {
				return consecutiveAssistantMessages
			}
		}
	}
	return consecutiveAssistantMessages
}

func (s *AIChatbotService) addHeader(req *http.Request, isContentJson bool) {
	if isContentJson {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("api-key", s.azureOpenAIConfig.Key)
}

func (s *AIChatbotService) getUrl(path string) string {
	return fmt.Sprintf(
		"%s/openai/%s?api-version=%s",
		s.azureOpenAIConfig.Endpoint,
		path,
		s.azureOpenAIConfig.ApiVersion,
	)
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
