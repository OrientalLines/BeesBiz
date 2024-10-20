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
		PostgresUser:     v.GetString("PSQL_USER"),
		PostgresPassword: v.GetString("PSQL_PASSWORD"),
		PostgresDB:       v.GetString("PSQL_DB"),
		PostgresHost:     v.GetString("PSQL_HOST"),
		PostgresPort:     v.GetString("PSQL_PORT"),
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
