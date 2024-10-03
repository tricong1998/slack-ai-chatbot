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
	WebhookURL string `mapstructure:"SLACK_WEBHOOK_URL"`
	Channel    string `mapstructure:"SLACK_CHANNEL_ID"`
	Token      string `mapstructure:"SLACK_TOKEN"`
	BotToken   string `mapstructure:"SLACK_BOT_TOKEN"`
}

type AzureOpenAIConfig struct {
	Endpoint    string `mapstructure:"AZURE_OPENAI_ENDPOINT"`
	Key         string `mapstructure:"AZURE_OPENAI_KEY"`
	ApiVersion  string `mapstructure:"AZURE_OPENAI_API_VERSION"`
	AssistantId string `mapstructure:"AZURE_OPENAI_ASSISTANT_ID"`
}

type Config struct {
	Server      ServerConfig
	DB          DBConfig
	Auth        AuthConfig
	SlackConfig SlackConfig
	AzureOpenAI AzureOpenAIConfig
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

	config := Config{
		Server:      server,
		DB:          db,
		Auth:        auth,
		SlackConfig: slackConfig,
		AzureOpenAI: azureOpenAI,
	}
	return config, nil
}
