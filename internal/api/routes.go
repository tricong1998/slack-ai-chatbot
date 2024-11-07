package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/api/handlers"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/google_internal"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/repository"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/services"
	"github.com/sotatek-dev/hyper-automation-chatbot/internal/shared"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/gin/middleware"
	"github.com/sotatek-dev/hyper-automation-chatbot/pkg/token"
)

func SetupRoutes(
	routes *gin.Engine,
	dependencies *shared.AppDependencies,
) {
	tokenMaker, err := token.NewJWTMaker(dependencies.Config.Auth.AccessTokenSecret)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Cannot create token maker")
		return
	}
	jwtService := services.NewJwtService(tokenMaker, dependencies.Config.Auth)
	userRepo := repository.NewUserRepository(dependencies.DB)
	userPointRepo := repository.NewUserPointRepository(dependencies.DB)
	userPointService := services.NewUserPointService(userPointRepo)
	userService := services.NewUserService(userRepo, userPointService)
	userHandler := handlers.NewUserHandler(userService, jwtService)

	slackService := services.NewSlackService(&dependencies.Config.SlackConfig, dependencies.SlackClient)
	ggSheetService := services.NewGSheetService(
		google_internal.GetSheetService(&dependencies.Config.Google),
		google_internal.GetDriveService(&dependencies.Config.Google),
	)
	slackHandler := handlers.NewSlackHandler(slackService, ggSheetService)

	threadRepo := repository.NewThreadRepository(dependencies.DB)
	threadService := services.NewThreadService(threadRepo)

	messageRepo := repository.NewMessageRepository(dependencies.DB)
	messageService := services.NewMessageService(messageRepo)

	aiChatbotService := services.NewAIChatbotService(dependencies.Config.AzureOpenAI, slackService, threadService, messageService)
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
	slackRoutes.Use(slackHandler.VerifySlackRequest())
	{
		slackRoutes.POST("/send-message", slackHandler.SendMessage)
		slackRoutes.POST("/actions", slackHandler.HandleBlockActions)
	}

	aiAssistantRoutes := routes.Group("/ai-assistant")
	{
		aiAssistantRoutes.POST("/add-message", aiChatbotHandler.AddMessage)
	}

	sheetHandler := handlers.NewSheetHandler(ggSheetService)
	sheetRoutes := routes.Group("/sheets")
	{
		sheetRoutes.POST("/candidate-offer", sheetHandler.ReadCandidateOffer)
		// sheetRoutes.POST("/create-new-sheet", sheetHandler.CreateNewSheet)
		sheetRoutes.POST("/handle-file-candidate-offer", sheetHandler.HandleFileCandidateOffer)
	}

	uiPathService := services.NewUIPathService(http.DefaultClient, dependencies.Config.UIPath)
	uiPathHandler := handlers.NewUIPathHandler(uiPathService)
	uiPathRoutes := routes.Group("/ui-path")
	{
		uiPathRoutes.POST("/greeting-new-employee", uiPathHandler.GreetingNewEmployee)
		uiPathRoutes.GET("/job-details/:jobID", uiPathHandler.GetJobDetails)
	}
}
