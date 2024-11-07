package shared

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/google_internal"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/logger"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type AppDependencies struct {
	UIPathJobService *services.UIPathJobService
	Logger           *zerolog.Logger
	UIPathJobRepo    *repository.UIPathJobRepository
	DB               *gorm.DB
	RabbitConn       *amqp.Connection
	AiChatbotService *services.AIChatbotService
	SlackService     *services.SlackService
	GgSheetService   *services.GSheetService
	UiPathService    *services.UIPathService
	SlackClient      *slack.Client
	ThreadService    *services.ThreadService
	MessageService   *services.MessageService
	ThreadRepo       *repository.ThreadRepository
	MessageRepo      *repository.MessageRepository
	Config           *config.Config
}

func InitDependencies(db *gorm.DB, rabbitConn *amqp.Connection, cfg *config.Config) AppDependencies {
	uiPathJobRepo := repository.NewUIPathJobRepository(db)
	logger := logger.NewLogger()
	slackClient := slack.New(
		cfg.SlackConfig.BotToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(cfg.SlackConfig.Token),
	)
	slackService := services.NewSlackService(&cfg.SlackConfig, slackClient)
	threadRepo := repository.NewThreadRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	threadService := services.NewThreadService(threadRepo)
	messageService := services.NewMessageService(messageRepo)
	aiChatbotService := services.NewAIChatbotService(cfg.AzureOpenAI, slackService, threadService, messageService)
	ggSheetService := services.NewGSheetService(
		google_internal.GetSheetService(&cfg.Google),
		google_internal.GetDriveService(&cfg.Google),
	)
	uiPathService := services.NewUIPathService(http.DefaultClient, cfg.UIPath)

	return AppDependencies{
		UIPathJobService: services.NewUIPathJobService(uiPathJobRepo,
			rabbitmq.NewPublisher(context.Background(),
				&rabbitmq.RabbitMQConfig{
					Host:     cfg.RabbitMQConfig.Host,
					Port:     cfg.RabbitMQConfig.Port,
					User:     cfg.RabbitMQConfig.User,
					Password: cfg.RabbitMQConfig.Password,
				},
				rabbitConn,
				logger,
				rabbitmq.HYPER_AUTOMATE_CHATBOT,
				"direct",
				rabbitmq.WELCOME_NEW_EMPLOYEE_QUEUE,
			),
			uiPathService,
			slackService,
		),
		Logger:           &logger,
		UIPathJobRepo:    uiPathJobRepo,
		DB:               db,
		RabbitConn:       rabbitConn,
		AiChatbotService: aiChatbotService,
		SlackService:     slackService,
		GgSheetService:   ggSheetService,
		UiPathService:    uiPathService,
		ThreadService:    threadService,
		MessageService:   messageService,
		ThreadRepo:       threadRepo,
		MessageRepo:      messageRepo,
		Config:           cfg,
		SlackClient:      slackClient,
	}
}
