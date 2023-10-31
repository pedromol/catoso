package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken string `mapstructure:"TELEGRAM_BOT"`
	TelegramChat  string `mapstructure:"TELEGRAM_CHAT"`
	OnvifIP       string `mapstructure:"ONVIF_IP"`
	OnvifPort     string `mapstructure:"ONVIF_PORT"`
	InputImage    string `mapstructure:"INPUT_IMAGE"`
	CascadePath   string `mapstructure:"CASCADE_PATH"`
}

func NewConfig() (Config, error) {
	cfg := Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT"),
		TelegramChat:  os.Getenv("TELEGRAM_CHAT"),
		OnvifIP:       os.Getenv("ONVIF_IP"),
		OnvifPort:     os.Getenv("ONVIF_PORT"),
		InputImage:    os.Getenv("INPUT_IMAGE"),
		CascadePath:   os.Getenv("CASCADE_PATH"),
	}

	if cfg.TelegramToken == "" {
		return cfg, errors.New("missing TELEGRAM_BOT env")
	}
	if cfg.TelegramChat == "" {
		return cfg, errors.New("missing TELEGRAM_CHAT env")
	}
	if cfg.OnvifIP == "" {
		return cfg, errors.New("missing ONVIF_IP env")
	}
	if cfg.OnvifPort == "" {
		return cfg, errors.New("missing ONVIF_PORT env")
	}
	if cfg.InputImage == "" {
		return cfg, errors.New("missing INPUT_IMAGE env")
	}
	if cfg.CascadePath == "" {
		return cfg, errors.New("missing CASCADE_PATH env")
	}

	return cfg, nil
}
