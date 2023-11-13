package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken   string
	TelegramChat    string
	OnvifIP         string
	OnvifPort       string
	InputImage      string
	InputFps        string
	CascadePath     string
	CenterCamera    string
	CatosoDebug     string
	StreamPort      string
	OutputFrameSkip string
}

func NewConfig() (*Config, error) {
	cfg := Config{
		TelegramToken:   os.Getenv("TELEGRAM_TOKEN"),
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
		return nil, errors.New("missing TELEGRAM_TOKEN env")
	}
	if cfg.TelegramChat == "" {
		return nil, errors.New("missing TELEGRAM_CHAT env")
	}
	if cfg.CenterCamera != "" {
		if cfg.OnvifIP == "" {
			return nil, errors.New("missing ONVIF_IP env")
		}
		if cfg.OnvifPort == "" {
			return nil, errors.New("missing ONVIF_PORT env")
		}
	}
	if cfg.InputImage == "" {
		return nil, errors.New("missing INPUT_IMAGE env")
	}
	if cfg.CascadePath == "" {
		return nil, errors.New("missing CASCADE_PATH env")
	}

	_, err := os.Stat(cfg.CascadePath)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
