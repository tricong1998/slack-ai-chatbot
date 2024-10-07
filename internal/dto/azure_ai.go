package dto

type AzureAIChatbotMessage struct {
	ID      string                         `json:"id"`
	Role    string                         `json:"role"`
	Content []AzureAIChatbotMessageContent `json:"content"`
}

type AzureAIChatbotMessageContent struct {
	Type string `json:"type"`
	Text struct {
		Value       string `json:"value"`
		Annotations []struct {
			Type         string `json:"type"`
			FileCitation struct {
				FileID string `json:"file_id"`
			} `json:"file_citation"`
			Text string `json:"text"`
		} `json:"annotations"`
	} `json:"text"`
}
