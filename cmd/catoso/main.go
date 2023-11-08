package main

import (
	"context"
	"log"
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

	cam := camera.NewCamera(cfg.OnvifIP, cfg.OnvifPort)

	enc := encoder.NewEncoder(cfg.InputImage)
	w, h, err := enc.GetVideoSize()
	if err != nil {
		panic(err)
	}

	vis := vision.NewVision(cfg.CascadePath, w, h)

	for {
		if cfg.CenterCamera != "" {
			if err := cam.Centralize(cfg.CenterCamera); err != nil {
				panic(err)
			}
		}

		ctx, cancel := context.WithCancel(context.TODO())
		pr1, pw1 := io.Pipe()
		er1, ew1 := io.Pipe()

		ffchan := enc.ReadStream(ctx, pw1, ew1)
		errchan := enc.Catch(ctx, er1)
		cvimg, cvchan := vis.Process(ctx, pr1, cfg.CatosoDebug)

		handlers := 0
	loop:
		for {
			select {
			case ff := <-ffchan:
				cancel()
				if ff != nil {
					log.Println("ffmpeg error: ", ff)
				} else {
					log.Println("ffmpeg finished with nil error")
				}
				pw1.Close()
				ew1.Close()
				er1.Close()
				handlers = handlers + 1
			case cv := <-cvchan:
				cancel()
				if cv != nil {
					log.Println("vision error: ", cv)
				} else {
					log.Println("vision finished with nil error")
				}
				close(cvimg)
				pr1.Close()
				handlers = handlers + 1
			case img := <-cvimg:
				if img != nil {
					if err := tel.SendPhoto(chatId, img); err != nil {
						cancel()
						log.Println("SendPhoto error: ", err)
					}
				}
			case err := <-errchan:
				cancel()
				if err != nil {
					log.Print("duplicated frames error: ", err)
				} else {
					log.Println("duplicated frames finished with nil error")
				}
				er1.Close()
				handlers = handlers + 1
			case <-ctx.Done():
				cancel()
				if handlers == 3 {
					log.Println("context is clear")
					break loop
				}
			case <-time.After(1 * time.Second):
				//noop
			}
		}
	}
}
