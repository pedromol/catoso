package config

import (
	"errors"
	"os"
	"regexp"
)

type Config struct {
	TelegramToken       string `json:"-"`
	TelegramChat        string `json:"telegramChat"`
	TelegramAlways      string `json:"telegramAlways"`
	OnvifIP             string `json:"onvifIP"`
	OnvifPort           string `json:"onvifPort"`
	InputImage          string `json:"inputImage"`
	InputProtocol       string `json:"inputProtocol"`
	InputFps            string `json:"inputFps"`
	InputPrefix         string `json:"inputPrefix"`
	InputHost           string `json:"inputHost"`
	InputSuffix         string `json:"inputSuffix"`
	CascadePath         string `json:"cascadePath"`
	CenterCamera        string `json:"centerCamera"`
	CatosoDebug         string `json:"catosoDebug"`
	StreamPort          string `json:"streamPort"`
	OutputFrameSkip     string `json:"outputFrameSkip"`
	DelayAfterDetectMin string `json:"delayAfterDetectMin"`
	RequiredDetections  string `json:"requiredDetections"`
	DrawOverFace        string `json:"drawOverFace"`
	ExitAfterMin        string `json:"exitAfterMin"`
	BucketURI           string `json:"bucketURI"`
	BucketName          string `json:"bucketName"`
	BucketKey           string `json:"bucketKey"`
	BucketSecret        string `json:"-"`
	HealthURI           string `json:"healthURI"`
}

func NewConfig() (*Config, error) {
	cfg := Config{
		TelegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		TelegramChat:        os.Getenv("TELEGRAM_CHAT"),
		TelegramAlways:      os.Getenv("TELEGRAM_ALWAYS"),
		OnvifIP:             os.Getenv("ONVIF_IP"),
		OnvifPort:           os.Getenv("ONVIF_PORT"),
		InputImage:          os.Getenv("INPUT_IMAGE"),
		InputProtocol:       os.Getenv("INPUT_PROTOCOL"),
		CascadePath:         os.Getenv("CASCADE_PATH"),
		CenterCamera:        os.Getenv("CENTER_CAMERA"),
		CatosoDebug:         os.Getenv("CATOSO_DEBUG"),
		InputFps:            os.Getenv("INPUT_FPS"),
		InputPrefix:         os.Getenv("INPUT_PREFIX"),
		InputHost:           os.Getenv("INPUT_HOST"),
		InputSuffix:         os.Getenv("INPUT_SUFFIX"),
		StreamPort:          os.Getenv("STREAM_PORT"),
		OutputFrameSkip:     os.Getenv("OUTPUT_FRAMESKIP"),
		DelayAfterDetectMin: os.Getenv("DELAY_AFTER_DETECT_MIN"),
		RequiredDetections:  os.Getenv("REQUIRED_DETECTIONS"),
		DrawOverFace:        os.Getenv("DRAW_OVER_FACE"),
		ExitAfterMin:        os.Getenv("EXIT_AFTER_MIN"),
		BucketURI:           os.Getenv("BUCKET_URI"),
		BucketName:          os.Getenv("BUCKET_NAME"),
		BucketKey:           os.Getenv("BUCKET_KEY"),
		BucketSecret:        os.Getenv("BUCKET_SECRET"),
		HealthURI:           os.Getenv("HEALTH_URI"),
	}

	if cfg.InputImage == "" {
		cfg.InputImage = cfg.InputPrefix + cfg.InputHost + cfg.InputSuffix
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

	if cfg.HealthURI == "" {
		s := regexp.MustCompile(`^[a-z]+://([^:]+):([0-9]+)`).FindStringSubmatch(cfg.InputImage)
		if len(s) != 3 {
			return nil, errors.New("invalid health url")
		}
		cfg.HealthURI = s[1] + ":" + s[2]
	}

	return &cfg, nil
}
