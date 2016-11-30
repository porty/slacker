package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	SlackMessageURL string   `required:"true"`
	SlackChannels   []string `required:"true"`
	Port            int      `default:"8000"`
	Username        string   `required:"true"`
	Password        string   `required:"true"`
}

func loadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("slacker", &c)
	return &c, err
}
