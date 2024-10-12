package clients

import (
	"cms/utils/constants"
	"github.com/spf13/viper"
)

type Environment struct {
	SystemDbUrl          string `mapstructure:"SYSTEM_DB_URL"`
	UserDbUrl            string `mapstructure:"USER_DB_URL"`
	Mode                 string `mapstructure:"MODE"`
	Port                 string `mapstructure:"PORT"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	AccessTokenValidity  int    `mapstructure:"ACCESS_TOKEN_VALIDITY"`
	RefreshTokenValidity int    `mapstructure:"REFRESH_TOKEN_VALIDITY"`
}

type Config interface {
	GetSystemDbUrl() string
	GetMode() string
	GetPort() string
	GetUserDbUrl() string
	GetJWTSecret() string
	GetAccessTokenValidity() int
	GetRefreshTokenValidity() int
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

func (c *Environment) GetSystemDbUrl() string {
	return c.SystemDbUrl
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

func (c *Environment) GetUserDbUrl() string {
	return c.UserDbUrl
}

func (c *Environment) GetJWTSecret() string {
	return c.JWTSecret
}

func (c *Environment) GetAccessTokenValidity() int {
	return c.AccessTokenValidity
}

func (c *Environment) GetRefreshTokenValidity() int {
	return c.RefreshTokenValidity
}
