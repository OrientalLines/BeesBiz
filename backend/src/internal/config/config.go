package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresUser     string `mapstructure:"PSQL_USER"`
	PostgresPassword string `mapstructure:"PSQL_PASSWORD"`
	PostgresDB       string `mapstructure:"PSQL_DB"`
	PostgresHost     string `mapstructure:"PSQL_HOST"`
	PostgresPort     string `mapstructure:"PSQL_PORT"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	TiKV             TiKVConfig
	RabbitMQ         RabbitMQConfig
	App              AppConfig
	API              APIConfig
}

type AppConfig struct {
	Environment string `mapstructure:"APP_ENV"`
}

type APIConfig struct {
	LimitEnabled    bool `mapstructure:"API_LIMIT_ENABLED"`
	LimitAmount     int  `mapstructure:"API_LIMIT_AMOUNT"`
	LimitExpiration int  `mapstructure:"API_LIMIT_EXPIRATION"`
}

type TiKVConfig struct {
	PDEndpoints []string `mapstructure:"TIKV_PD_ENDPOINTS"`
	Username    string   `mapstructure:"TIKV_USERNAME"`
	Password    string   `mapstructure:"TIKV_PASSWORD"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"RABBITMQ_HOST"`
	Port     string `mapstructure:"RABBITMQ_PORT"`
	Username string `mapstructure:"RABBITMQ_USERNAME"`
	Password string `mapstructure:"RABBITMQ_PASSWORD"`
	VHost    string `mapstructure:"RABBITMQ_VHOST"`
}

var GlobalConfig Config

func LoadConfig() error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.ReadInConfig()

	// Populate GlobalConfig using Viper
	GlobalConfig = Config{
		PostgresUser:     v.GetString("POSTGRES_USER"),
		PostgresPassword: v.GetString("POSTGRES_PASSWORD"),
		PostgresDB:       v.GetString("POSTGRES_DB"),
		PostgresHost:     v.GetString("POSTGRES_HOST"),
		PostgresPort:     v.GetString("POSTGRES_PORT"),
		JwtSecret:        v.GetString("JWT_SECRET"),
		TiKV: TiKVConfig{
			PDEndpoints: v.GetStringSlice("TIKV_PD_ENDPOINTS"),
			Username:    v.GetString("TIKV_USERNAME"),
			Password:    v.GetString("TIKV_PASSWORD"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     v.GetString("RABBITMQ_HOST"),
			Port:     v.GetString("RABBITMQ_PORT"),
			Username: v.GetString("RABBITMQ_USERNAME"),
			Password: v.GetString("RABBITMQ_PASSWORD"),
			VHost:    v.GetString("RABBITMQ_VHOST"),
		},
		App: AppConfig{
			Environment: v.GetString("APP_ENV"),
		},
		API: APIConfig{
			LimitAmount:     v.GetInt("API_LIMIT_AMOUNT"),
			LimitExpiration: v.GetInt("API_LIMIT_EXPIRATION"),
			LimitEnabled:    v.GetBool("API_LIMIT_ENABLED"),
		},
	}

	return nil
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
