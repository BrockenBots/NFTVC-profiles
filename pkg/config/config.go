package config

import (
	"fmt"
	"nftvc-profile/pkg/logger"
	"nftvc-profile/pkg/mongodb"

	"github.com/spf13/viper"
)

type Config struct {
	Http             Http              `mapstructure:"http" validate:"required"`
	Logger           *logger.Config    `mapstructure:"logger" validate:"required"`
	Mongo            *mongodb.Config   `mapstructure:"mongo" validate:"required"`
	MongoCollections *MongoCollections `mapstructure:"mongoCollections" validate:"required"`
	AuthClientUrl    string            `mapstructure:"authClientURL" validate:"required"`
	S3Url            string            `mapstructure:"s3Url" validate:"required"`
}

type MongoCollections struct {
	Profiles string `mapstructure:"profiles" validate:"required"`
}

type Http struct {
	Port string `mapstructure:"port" validate:"required"`
}

func LoadConfig(pathToConfig string) (*Config, error) {
	if pathToConfig == "" {
		return nil, fmt.Errorf("path to cfg is empty")
	}

	cfg := &Config{}

	viper.SetConfigType("yml")
	viper.SetConfigFile(pathToConfig)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("(ReadInConfig) error: %v", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cfg")
	}

	return cfg, nil
}
