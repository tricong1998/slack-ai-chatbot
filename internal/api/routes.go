package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/api/handlers"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/gin/middleware"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/token"
	"gorm.io/gorm"
)

func SetupRoutes(
	routes *gin.Engine,
	db *gorm.DB,
	config *config.Config,
	log *zerolog.Logger,
	slackClient *slack.Client,
) {
	tokenMaker, err := token.NewJWTMaker(config.Auth.AccessTokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create token maker")
		return
	}
	jwtService := services.NewJwtService(tokenMaker, config.Auth)
	userRepo := repository.NewUserRepository(db)
	userPointRepo := repository.NewUserPointRepository(db)
	userPointService := services.NewUserPointService(userPointRepo)
	userService := services.NewUserService(userRepo, userPointService)
	userHandler := handlers.NewUserHandler(userService, jwtService)

	slackService := services.NewSlackService(&config.SlackConfig, slackClient)
	slackHandler := handlers.NewSlackHandler(slackService)

	threadRepo := repository.NewThreadRepository(db)
	threadService := services.NewThreadService(threadRepo)

	messageRepo := repository.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)

	aiChatbotService := services.NewAIChatbotService(config.AzureOpenAI, slackService, threadService, messageService)
	aiChatbotHandler := handlers.NewAIChatbotHandler(aiChatbotService)

	userGroup := routes.Group("users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.POST("/login", userHandler.Login)
	}
	authRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{}))
	{
		authRoutes.GET("/me", userHandler.ReadMe)
		authRoutes.GET("/:id", userHandler.ReadUser)
		authRoutes.PUT("/update-me", userHandler.UpdateMe)
		authRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	adminRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{"admin"}))
	{
		adminRoutes.GET("", userHandler.ListUsers)
	}

	//TODO: remove after testing AI Chatbot, Slack done
	slackRoutes := routes.Group("/slack")
	{
		slackRoutes.POST("/send-message", slackHandler.SendMessage)
	}

	aiAssistantRoutes := routes.Group("/ai-assistant")
	{
		aiAssistantRoutes.POST("/add-message", aiChatbotHandler.AddMessage)
	}
}
