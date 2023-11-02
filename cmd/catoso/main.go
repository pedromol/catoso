package main

import (
	"fmt"
	"strconv"
	"time"

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
		if err := cam.Centralize(cfg.CenterCamera); err != nil {
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
	errchan := enc.Catch(er1)
	cvimg, cvchan := vis.Process(pr1, cfg.CatosoDebug)

	for {
		select {
		case ff := <-ffchan:
			fmt.Println(ff)
			panic(ff)
		case cv := <-cvchan:
			panic(cv)
		case img := <-cvimg:
			if err := tel.SendPhoto(chatId, img); err != nil {
				panic(err)
			}
		case err := <-errchan:
			panic(err)
		case <-time.After(1 * time.Second):
			//noop
		}
	}
}
