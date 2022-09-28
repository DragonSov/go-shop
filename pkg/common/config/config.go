package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Addr         string `mapstructure:"ADDR"`
	DSN          string `mapstructure:"DSN"`
	JWTLifeTime  int    `mapstructure:"JWT_LIFETIME"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./pkg/common/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
