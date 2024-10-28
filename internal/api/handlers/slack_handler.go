package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
)

type SlackHandler struct {
	slackService   services.ISlackService
	ggSheetService services.IGSheetService
}

func NewSlackHandler(
	slackService services.ISlackService,
	ggSheetService services.IGSheetService,
) *SlackHandler {
	return &SlackHandler{
		slackService:   slackService,
		ggSheetService: ggSheetService,
	}
}

func (h *SlackHandler) VerifySlackRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		slackSigningSecret := h.slackService.GetSigningSecret() // Implement this method in your slack service

		timestamp := c.GetHeader("X-Slack-Request-Timestamp")
		slackSignature := c.GetHeader("X-Slack-Signature")

		// Check if the timestamp is within 5 minutes of current time
		currentTime := time.Now().Unix()
		t, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || currentTime-t > 300 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid timestamp"})
			return
		}

		// Read the request body
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		// Restore the body for later use
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Create the signature base string
		baseString := fmt.Sprintf("v0:%s:%s", timestamp, string(body))

		// Create the HMAC-SHA256 signature
		mac := hmac.New(sha256.New, []byte(slackSigningSecret))
		mac.Write([]byte(baseString))
		calculatedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

		// Compare signatures
		if !hmac.Equal([]byte(calculatedSignature), []byte(slackSignature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		c.Next()
	}
}

func (h *SlackHandler) SendMessage(ctx *gin.Context) {
	var input dto.SendMessageDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.slackService.SendMessage(ctx, input.ChannelID, input.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message sent to Slack"})
}

func (s *SlackHandler) HandleBlockActions(c *gin.Context) {
	var payload slack.InteractionCallback
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse action response"})
		return
	}

	// Check if this is a rating submission
	if payload.BlockID == "rating_input" {
		ratingStr := payload.ActionCallback.BlockActions[0].Value
		rating, err := strconv.Atoi(ratingStr)
		if err != nil || rating < 1 || rating > 5 {
			// Handle invalid input
			s.sendErrorMessage(payload.Channel.ID, "Invalid rating. Please enter a number between 1 and 5.")
			c.Status(http.StatusOK)
			return
		}

		// Process the rating
		if err := s.processRating(payload.User.ID, rating); err != nil {
			s.sendErrorMessage(payload.Channel.ID, "Failed to process rating. Please try again.")
			c.Status(http.StatusOK)
			return
		}

		// Send a thank you message
		thankYouMessage := fmt.Sprintf("Thank you for your rating of %d stars!", rating)
		if _, _, err := s.slackService.PostMessage(payload.Channel.ID, slack.MsgOptionText(thankYouMessage, false)); err != nil {
			// Handle error
			c.Status(http.StatusInternalServerError)
			return
		}
	} else if payload.BlockID == "candidate_file" {
		candidateFile := payload.ActionCallback.BlockActions[0].Value
		if err := s.processCandidateFile(payload.User.ID, candidateFile); err != nil {
			s.sendErrorMessage(payload.Channel.ID, "Failed to process candidate file. Please try again.")
			c.Status(http.StatusOK)
			return
		}
	}

	c.Status(http.StatusOK)
}

func (s *SlackHandler) sendErrorMessage(channelID, message string) {
	s.slackService.PostMessage(channelID, slack.MsgOptionText(message, false))
}

func (s *SlackHandler) processCandidateFile(_, fileLink string) error {

	// Here you would typically handle the file link, possibly saving it or processing it
	// For example, you might want to validate the file link, download the file, or extract data
	// This is a placeholder for actual processing logic

	return nil
}

func (s *SlackHandler) processRating(_ string, _ int) error {
	// Here you would typically handle the rating, possibly saving it or processing it
	// For example, you might want to validate the rating, store it, or use it for analysis
	// This is a placeholder for actual processing logic
	return nil
}
