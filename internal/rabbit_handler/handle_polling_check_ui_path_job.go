package rabbit_handler

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
	"github.com/streadway/amqp"
)

type PollingCheckUIPathJobDependencies struct {
	UIPathJobService *services.UIPathJobService
	Logger           *zerolog.Logger
}

func HandlePollingCheckUIPathJob(queue string, msg amqp.Delivery, dependencies *PollingCheckUIPathJobDependencies) error {
	dependencies.Logger.Info().Msgf("Message received on queue: %s with message: %s", queue, string(msg.Body))

	var input dto.UIPathCheckingJobInput
	if err := json.Unmarshal(msg.Body, &input); err != nil {
		return err
	}

	dependencies.Logger.Info().Msgf("Job ID: %d", input.JobID)
	completed, err := dependencies.UIPathJobService.PollingCheck(input.JobID)
	if err != nil {
		return err
	}
	dependencies.Logger.Info().Msgf("Job completed: %t", completed)
	return nil
}
