package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config stores all information from config file.
type Config struct {
	Port     string `toml:"port"`
	LogLevel string `toml:"loglevel"`
	Nats     struct {
		ServerID string `toml:"serverID"`
		ClientID string `toml:"clientID"`
		NatsURL  string `toml:"natsUrl"`
	} `toml:"nats"`
	DB struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	} `toml:"db"`
}

// GetConfig reads configuration file and stores it in Config.
func GetConfig() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalln(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatalln(err)
	}

	return cfg
}
