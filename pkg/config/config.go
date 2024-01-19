package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logrus.Error("no such config file")
		} else {
			// Config file was found but another error was produced
			logrus.Error("read config error")
		}
		logrus.Fatal(err)
	}
}
