package main

import (
	"context"
	"errors"
	"log"
	"net/http"
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

	var st *vision.Stream
	if cfg.StreamPort != "" {
		st = vision.NewStream()

		http.Handle("/", st)

		server := &http.Server{
			Addr:         "0.0.0.0:" + cfg.StreamPort,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		}

		defer server.Close()
		go func() {
			panic(server.ListenAndServe())
		}()
	}

	fskip := 0
	if cfg.OutputFrameSkip != "" {
		fskip, err = strconv.Atoi(cfg.OutputFrameSkip)
		if err != nil {
			panic(err)
		}
	}

	for {
		if cfg.CenterCamera != "" {
			if err := cam.Centralize(cfg.CenterCamera); err != nil {
				panic(err)
			}
		}

		ctx, cancel := context.WithCancel(context.TODO())
		pr1, pw1 := io.Pipe()
		er1, ew1 := io.Pipe()

		ffchan := enc.ReadStream(ctx, pw1, ew1, cfg.InputFps)
		errchan := enc.Catch(ctx, er1)
		cvimg, cvchan := vis.Process(ctx, pr1, st, fskip, cfg.CatosoDebug)

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
				handlers = handlers + 1
			case cv := <-cvchan:
				cancel()
				if cv != nil {
					log.Println("vision error: ", cv)
				} else {
					log.Println("vision finished with nil error")
				}
				close(cvimg)
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
				handlers = handlers + 1
			case <-ctx.Done():
				pw1.CloseWithError(errors.New("context cancelled"))
				ew1.CloseWithError(errors.New("context cancelled"))
				er1.CloseWithError(errors.New("context cancelled"))
				pr1.CloseWithError(errors.New("context cancelled"))
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
