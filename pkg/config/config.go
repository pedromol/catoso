package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken   string `mapstructure:"TELEGRAM_BOT"`
	TelegramChat    string `mapstructure:"TELEGRAM_CHAT"`
	OnvifIP         string `mapstructure:"ONVIF_IP"`
	OnvifPort       string `mapstructure:"ONVIF_PORT"`
	InputImage      string `mapstructure:"INPUT_IMAGE"`
	InputFps        string `mapstructure:"INPUT_FPS"`
	CascadePath     string `mapstructure:"CASCADE_PATH"`
	CenterCamera    string `mapstructure:"CENTER_CAMERA"`
	CatosoDebug     string `mapstructure:"CATOSO_DEBUG"`
	StreamPort      string `mapstructure:"STREAM_PORT"`
	OutputFrameSkip string `mapstructure:"OUTPUT_FRAMESKIP"`
}

func NewConfig() (Config, error) {
	cfg := Config{
		TelegramToken:   os.Getenv("TELEGRAM_BOT"),
		TelegramChat:    os.Getenv("TELEGRAM_CHAT"),
		OnvifIP:         os.Getenv("ONVIF_IP"),
		OnvifPort:       os.Getenv("ONVIF_PORT"),
		InputImage:      os.Getenv("INPUT_IMAGE"),
		CascadePath:     os.Getenv("CASCADE_PATH"),
		CenterCamera:    os.Getenv("CENTER_CAMERA"),
		CatosoDebug:     os.Getenv("CATOSO_DEBUG"),
		InputFps:        os.Getenv("INPUT_FPS"),
		StreamPort:      os.Getenv("STREAM_PORT"),
		OutputFrameSkip: os.Getenv("OUTPUT_FRAMESKIP"),
	}

	if cfg.TelegramToken == "" {
		return cfg, errors.New("missing TELEGRAM_BOT env")
	}
	if cfg.TelegramChat == "" {
		return cfg, errors.New("missing TELEGRAM_CHAT env")
	}
	if cfg.CenterCamera != "" {
		if cfg.OnvifIP == "" {
			return cfg, errors.New("missing ONVIF_IP env")
		}
		if cfg.OnvifPort == "" {
			return cfg, errors.New("missing ONVIF_PORT env")
		}
	}
	if cfg.InputImage == "" {
		return cfg, errors.New("missing INPUT_IMAGE env")
	}
	if cfg.CascadePath == "" {
		return cfg, errors.New("missing CASCADE_PATH env")
	}

	_, err := os.Stat(cfg.CascadePath)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
