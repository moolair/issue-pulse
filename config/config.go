package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	GitHub struct {
		Token  string   `json:"token"`
		Owner  string   `json:"owner"`
		Repo   string   `json:"repo"`
		Labels []string `json:"labels"`
	} `json:"github"`
	Slack struct {
		WebhookURL string `json:"webhook_url"`
	} `json:"slack"`
	PollIntervalSeconds int `json:"poll_interval_seconds"`
}

func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	if cfg.PollIntervalSeconds == 0 {
		cfg.PollIntervalSeconds = 60
	}
	return &cfg, nil
}
