package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configurations of the application.
// Values read by viper from a config file or environment variables.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig read configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// location of config file
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // type of config file

	// Tell Viper to automatically override config variables from file with environment variables
	// Useful for deploying to production or QA
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
