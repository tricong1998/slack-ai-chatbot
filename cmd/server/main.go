package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/api"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/database"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/slack_handlers"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/logger"
	"gorm.io/gorm"
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

	runGinServer(&cfg, db, &log)
}

func runGinServer(
	cfg *config.Config,
	db *gorm.DB,
	log *zerolog.Logger,
) {
	// Initialize router
	routes := gin.Default()
	slackClient := slack.New(
		cfg.SlackConfig.BotToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(cfg.SlackConfig.Token),
	)
	api.SetupRoutes(routes, db, cfg, log, slackClient)
	slackService := services.NewSlackService(&cfg.SlackConfig, slackClient)
	aiChatbotService := services.NewAIChatbotService(cfg.AzureOpenAI, slackService)
	go socket(slackClient, context.Background(), log, aiChatbotService)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err := routes.Run(address)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run server")
	}
}

func socket(client *slack.Client, ctx context.Context, zerolog *zerolog.Logger, aiChatbotService *services.AIChatbotService) {
	slackHandler := slack_handlers.NewSlackHandler(client, aiChatbotService)
	socketClient := socketmode.New(
		client,
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
				zerolog.Println("Shutting down socketmode listener")
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
						zerolog.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := slackHandler.HandleEventMessage(eventsAPIEvent)
					if err != nil {
						zerolog.Error().Err(err).Msg("Cannot handle event message")
					}
				case socketmode.EventTypeSlashCommand:
					// Just like before, type cast to the correct event type, this time a SlashEvent
					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
						continue
					}
					// Dont forget to acknowledge the request
					socketClient.Ack(*event.Request)
					// handleSlashCommand will take care of the command
					err := slackHandler.HandleSlashCommand(command, client)
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}(ctx, client, socketClient)

	socketClient.Run()
}
