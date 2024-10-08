package utils

import (
	"cms/utils/constants"
	"github.com/spf13/viper"
)

type Environment struct {
	DbUrl string `mapstructure:"DB_URL"`
	Mode  string `mapstructure:"MODE"`
	Port  string `mapstructure:"PORT"`
}

type Config interface {
	GetDbUrl() string
	GetMode() string
	GetPort() string
}

func NewConfig() (Config, error) {
	var config Environment
	viper.AddConfigPath("./")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&config)
	return &config, err
}

func (c *Environment) GetDbUrl() string {
	return c.DbUrl
}

func (c *Environment) GetMode() string {
	if c.Mode != "" {
		return c.Mode
	}

	return string(constants.Dev)
}

func (c *Environment) GetPort() string {
	if c.Port != "" {
		return c.Port
	}

	return constants.DefaultPort
}
