package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/api"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/database"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/dto"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/rabbit_handler"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/shared"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/slack_handlers"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/logger"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/rabbitmq"
)

func main() {
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	// Initialize db
	db, err := database.Initialize(&cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize database")
	}

	// Migrate db
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot migrate database")
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect rabbit")
	}

	rabbitConfig := rabbitmq.RabbitMQConfig{
		Host:     cfg.RabbitMQConfig.Host,
		Port:     cfg.RabbitMQConfig.Port,
		User:     cfg.RabbitMQConfig.User,
		Password: cfg.RabbitMQConfig.Password,
	}
	rabbitConn, err := rabbitmq.NewRabbitMQConn(&rabbitConfig, context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect rabbit")
	}
	dependencies := shared.InitDependencies(db, rabbitConn, &cfg)
	welcomeNewEmployeeDependencies := rabbit_handler.WelcomeNewEmployeeDependencies{
		UIPathJobService: dependencies.UIPathJobService,
		Logger:           dependencies.Logger,
	}
	greetingConsumer := rabbitmq.NewConsumer[*rabbit_handler.WelcomeNewEmployeeDependencies](
		context.Background(),
		&rabbitConfig,
		rabbitConn,
		*dependencies.Logger,
		rabbit_handler.HandlePollingCheckUIPathJob,
		rabbitmq.HYPER_AUTOMATE_CHATBOT,
		"direct",
		rabbitmq.WELCOME_NEW_EMPLOYEE_QUEUE,
		rabbitmq.WELCOME_NEW_EMPLOYEE_QUEUE,
	)
	go func() {
		err := greetingConsumer.ConsumeMessage(dto.CreateUserPoint{}, &welcomeNewEmployeeDependencies)
		if err != nil {
			log.Error().Err(err).Msg("Consume message error")
		}
	}()

	go socket(
		context.Background(),
		&dependencies,
	)
	// dependencies.UIPathJobService.PollingCheck(418824812)
	runGinServer(&dependencies)
}

func runGinServer(
	dependencies *shared.AppDependencies,
) {
	// Initialize router
	routes := gin.Default()
	api.SetupRoutes(routes, dependencies)

	// Start server
	address := fmt.Sprintf("%s:%s", dependencies.Config.Server.Host, dependencies.Config.Server.Port)
	err := routes.Run(address)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Cannot run server")
	}
}

func socket(ctx context.Context,
	dependencies *shared.AppDependencies,
) {
	slackHandler := slack_handlers.NewSlackHandler(
		dependencies.SlackClient,
		dependencies.SlackService,
		dependencies.AiChatbotService,
		dependencies.GgSheetService,
		dependencies.UIPathJobService,
	)
	socketClient := socketmode.New(
		dependencies.SlackClient,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				dependencies.Logger.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						dependencies.Logger.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := slackHandler.HandleEventMessage(eventsAPIEvent)
					if err != nil {
						dependencies.Logger.Error().Err(err).Msg("Cannot handle event message")
					}
				case socketmode.EventTypeSlashCommand:
					// Just like before, type cast to the correct event type, this time a SlashEvent
					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
						continue
					}
					// handleSlashCommand will take care of the command
					payload, err := slackHandler.HandleSlashCommand(command, client)
					if err != nil {
						log.Fatal(err)
					}
					// Dont forget to acknowledge the request
					socketClient.Ack(*event.Request, payload)
				case socketmode.EventTypeInteractive:
					interactionCallback, ok := event.Data.(slack.InteractionCallback)
					if !ok {
						dependencies.Logger.Printf("Could not type cast the event to an InteractionCallback: %v\n", event)
						continue
					}
					// handleSlashCommand will take care of the command
					payload, err := slackHandler.HandleBlockAction(interactionCallback)
					if err != nil {
						dependencies.Logger.Error().Err(err).Msg("Cannot handle block actions")
					}
					socketClient.Ack(*event.Request, payload)
				}
			}
		}
	}(ctx, dependencies.SlackClient, socketClient)

	socketClient.Run()
}
