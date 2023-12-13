package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type DB struct {
	Host     string
	DbName   string
	Port     string
	User     string
	Password string
}

type Config struct {
	DB
	BaseURL          string
	DownloadBasePath string
}

func loadConfig() (*Config, error) {
	var conf Config

	_, err := os.Stat(filepath.Join(".", "conf", "config-local.yml"))
	if err == nil {
		viper.SetConfigName("config-local")
	} else {
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(".", "conf"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %w", err)
	}

	return &conf, nil
}
