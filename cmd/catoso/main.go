package main

import (
	"fmt"
	"strconv"
	"strings"

	"io"

	"github.com/pedromol/catoso/pkg/camera"
	"github.com/pedromol/catoso/pkg/config"
	"github.com/pedromol/catoso/pkg/encoder"
	"github.com/pedromol/catoso/pkg/telegram"
	"github.com/pedromol/catoso/pkg/vision"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	tel, err := telegram.NewTelegram(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}

	chatId, err := strconv.ParseInt(cfg.TelegramChat, 10, 64)
	if err != nil {
		panic(err)
	}

	if cfg.CenterCamera != "" {
		cam := camera.NewCamera(cfg.OnvifIP, cfg.OnvifPort)
		if err := cam.Centralize(); err != nil {
			panic(err)
		}
	}

	enc := encoder.NewEncoder(cfg.InputImage)
	w, h, err := enc.GetVideoSize()
	if err != nil {
		panic(err)
	}

	vis := vision.NewVision(cfg.CascadePath, w, h)

	pr1, pw1 := io.Pipe()
	er1, ew1 := io.Pipe()

	ffchan := enc.ReadStream(pw1, ew1)
	cvimg, cvchan := vis.Process(pr1)

	for {
		select {
		case ff := <-ffchan:
			panic(ff)
		case cv := <-cvchan:
			panic(cv)
		case img := <-cvimg:
			if err := tel.SendPhoto(chatId, img); err != nil {
				panic(err)
			}
		default:
			buf := make([]byte, 1024)
			er1.Read(buf)
			if strings.Contains(string(buf), "More than 1000 frames duplicated") {
				fmt.Print(string(buf))
				return
			}
		}
	}
}
