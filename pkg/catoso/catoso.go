package catoso

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/pedromol/catoso/pkg/camera"
	"github.com/pedromol/catoso/pkg/config"
	"github.com/pedromol/catoso/pkg/encoder"
	"github.com/pedromol/catoso/pkg/telegram"
	"github.com/pedromol/catoso/pkg/vision"
)

type Catoso struct {
	Config    *config.Config
	Telegram  *telegram.Telegram
	ChatId    int64
	Camera    *camera.Camera
	Encoder   *encoder.Encoder
	Vision    *vision.Vision
	Stream    *vision.Stream
	Frameskip int
	Context   context.Context
	Cancel    context.CancelFunc
	Handlers  int
}

func NewCatoso(cfg *config.Config) (*Catoso, error) {
	tel, err := telegram.NewTelegram(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	chatId, err := strconv.ParseInt(cfg.TelegramChat, 10, 64)
	if err != nil {
		return nil, err
	}

	cam := camera.NewCamera(cfg.OnvifIP, cfg.OnvifPort)

	enc := encoder.NewEncoder(cfg.InputImage)
	w, h, err := enc.GetVideoSize()
	if err != nil {
		return nil, err
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
			server.ListenAndServe()
		}()
	}

	fskip := 0
	if cfg.OutputFrameSkip != "" {
		fskip, err = strconv.Atoi(cfg.OutputFrameSkip)
		if err != nil {
			return nil, err
		}
	}

	return &Catoso{
		Config:    cfg,
		Telegram:  tel,
		ChatId:    chatId,
		Camera:    cam,
		Encoder:   enc,
		Vision:    vis,
		Stream:    st,
		Frameskip: fskip,
	}, nil

}

func (h *Catoso) Start() {
	for {
		if h.Config.CenterCamera != "" {
			if err := h.Camera.Centralize(h.Config.CenterCamera); err != nil {
				log.Println("Centralize error:  ", err)
			}
		}

		h.Context, h.Cancel = context.WithCancel(context.TODO())
		pr1, pw1 := io.Pipe()
		er1, ew1 := io.Pipe()

		ffchan := h.Encoder.ReadStream(h.Context, pw1, ew1, h.Config.InputFps)
		errchan := h.Encoder.Catch(h.Context, er1)
		cvimg, cvchan := h.Vision.Process(h.Context, pr1, h.Stream, h.Frameskip, h.Config.CatosoDebug)

		h.Handlers = 0
	loop:
		for {
			select {
			case ff := <-ffchan:
				h.Cancel()
				if ff != nil {
					log.Println("ffmpeg error: ", ff)
				} else {
					log.Println("ffmpeg finished with nil error")
				}
				h.Handlers += 1
			case cv := <-cvchan:
				h.Cancel()
				if cv != nil {
					log.Println("vision error: ", cv)
				} else {
					log.Println("vision finished with nil error")
				}
				close(cvimg)
				h.Handlers += 1
			case img := <-cvimg:
				if img != nil {
					if err := h.Telegram.SendPhoto(h.ChatId, img); err != nil {
						h.Cancel()
						log.Println("SendPhoto error: ", err)
					}
				}
			case err := <-errchan:
				h.Cancel()
				if err != nil {
					log.Print("duplicated frames error: ", err)
				} else {
					log.Println("duplicated frames finished with nil error")
				}
				h.Handlers += 1
			case <-h.Context.Done():
				h.Cancel()
				ctxErr := errors.New("context cancelled")
				pw1.CloseWithError(ctxErr)
				ew1.CloseWithError(ctxErr)
				er1.CloseWithError(ctxErr)
				pr1.CloseWithError(ctxErr)
				if h.Handlers > 2 {
					log.Println("context is clear")
					break loop
				}
			case <-time.After(1 * time.Second):
				//noop
			}
		}
	}
}