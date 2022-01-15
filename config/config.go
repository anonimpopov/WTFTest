package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server Server
	Mongo  Mongo
}

type Server struct {
	Port string
	Mode string
}

type Mongo struct {
	URL      string
	Database string
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config-local")
	v.AddConfigPath("config")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable unmarshal, %v", err)
		return nil, err
	}

	return &c, nil
}
