package config

import (
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

type ServerConfig struct {
	Host string `mapstructure:"SERVER_HOST"`
	Port string `mapstructure:"SERVER_PORT"`
}

type AuthConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	AccessTokenSecret    string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret   string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type SlackConfig struct {
	WebhookURL    string `mapstructure:"SLACK_WEBHOOK_URL"`
	Channel       string `mapstructure:"SLACK_CHANNEL_ID"`
	Token         string `mapstructure:"SLACK_TOKEN"`
	BotToken      string `mapstructure:"SLACK_BOT_TOKEN"`
	SigningSecret string `mapstructure:"SLACK_SIGNING_SECRET"`
}

type AzureOpenAIConfig struct {
	Endpoint                 string `mapstructure:"AZURE_OPENAI_ENDPOINT"`
	Key                      string `mapstructure:"AZURE_OPENAI_KEY"`
	ApiVersion               string `mapstructure:"AZURE_OPENAI_API_VERSION"`
	AssistantIdDetectAction  string `mapstructure:"AZURE_OPENAI_ASSISTANT_ID_DETECT_ACTION"`
	AssistantIdHeaderMapping string `mapstructure:"AZURE_OPENAI_ASSISTANT_ID_HEADER_MAPPING"`
}

type GoogleConfig struct {
	Credentials string `mapstructure:"GOOGLE_CREDENTIALS"`
}

type UIPathConfig struct {
	Host                              string `mapstructure:"UI_PATH_HOST"`
	Tenant                            string `mapstructure:"UI_PATH_TENANT"`
	TenantID                          string `mapstructure:"UI_PATH_TENANT_ID"`
	ApiKey                            string `mapstructure:"UI_PATH_API_KEY"`
	GreetingNewEmployeeProcessKey     string `mapstructure:"UI_PATH_GREETING_NEW_EMPLOYEE_PROCESS_KEY"`
	FillBuddyProcessKey               string `mapstructure:"UI_PATH_FILL_BUDDY_PROCESS_KEY"`
	CreateLeaveRequestProcessKey      string `mapstructure:"UI_PATH_CREATE_LEAVE_REQUEST_PROCESS_KEY"`
	CreateIntegrateTrainingProcessKey string `mapstructure:"UI_PATH_CREATE_INTEGRATE_TRAINING_PROCESS_KEY"`
	PreOnboardEmailProcessKey         string `mapstructure:"UI_PATH_PRE_ONBOARD_EMAIL_PROCESS_KEY"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"AMQP_SERVER_HOST"`
	Port     string `mapstructure:"AMQP_SERVER_PORT"`
	User     string `mapstructure:"AMQP_SERVER_USER"`
	Password string `mapstructure:"AMQP_SERVER_PASSWORD"`
}

type Config struct {
	Server         ServerConfig
	DB             DBConfig
	Auth           AuthConfig
	SlackConfig    SlackConfig
	AzureOpenAI    AzureOpenAIConfig
	Google         GoogleConfig
	UIPath         UIPathConfig
	RabbitMQConfig RabbitMQConfig
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var server ServerConfig
	var db DBConfig
	var auth AuthConfig
	var slackConfig SlackConfig
	var azureOpenAI AzureOpenAIConfig
	var google GoogleConfig
	var uiPath UIPathConfig
	var rabbitMQ RabbitMQConfig
	err = viper.Unmarshal(&server)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&db)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&auth)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&slackConfig)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&azureOpenAI)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&google)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&uiPath)
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&rabbitMQ)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Server:         server,
		DB:             db,
		Auth:           auth,
		SlackConfig:    slackConfig,
		AzureOpenAI:    azureOpenAI,
		Google:         google,
		UIPath:         uiPath,
		RabbitMQConfig: rabbitMQ,
	}
	return config, nil
}
