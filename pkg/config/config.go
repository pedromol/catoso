package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken       string `json:"-"`
	TelegramChat        string `json:"telegramChat"`
	OnvifIP             string `json:"onvifIP"`
	OnvifPort           string `json:"onvifPort"`
	InputImage          string `json:"inputImage"`
	InputFps            string `json:"inputFps"`
	CascadePath         string `json:"cascadePath"`
	CenterCamera        string `json:"centerCamera"`
	CatosoDebug         string `json:"catosoDebug"`
	StreamPort          string `json:"streamPort"`
	OutputFrameSkip     string `json:"outputFrameSkip"`
	DelayAfterDetectMin string `json:"delayAfterDetectMin"`
	DrawOverFace        string `json:"drawOverFace"`
	ExitAfterMin        string `json:"exitAfterMin"`
}

func NewConfig() (*Config, error) {
	cfg := Config{
		TelegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		TelegramChat:        os.Getenv("TELEGRAM_CHAT"),
		OnvifIP:             os.Getenv("ONVIF_IP"),
		OnvifPort:           os.Getenv("ONVIF_PORT"),
		InputImage:          os.Getenv("INPUT_IMAGE"),
		CascadePath:         os.Getenv("CASCADE_PATH"),
		CenterCamera:        os.Getenv("CENTER_CAMERA"),
		CatosoDebug:         os.Getenv("CATOSO_DEBUG"),
		InputFps:            os.Getenv("INPUT_FPS"),
		StreamPort:          os.Getenv("STREAM_PORT"),
		OutputFrameSkip:     os.Getenv("OUTPUT_FRAMESKIP"),
		DelayAfterDetectMin: os.Getenv("DELAY_AFTER_DETECT_MIN"),
		DrawOverFace:        os.Getenv("DRAW_OVER_FACE"),
		ExitAfterMin:        os.Getenv("EXIT_AFTER_MIN"),
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
