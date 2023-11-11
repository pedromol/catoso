package encoder

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	ffmpeg "github.com/pedromol/catoso/pkg/ffmpeg_go"
)

type VideoInfo struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int
		Height    int
	} `json:"streams"`
}

type Encoder struct {
	InputImage string
}

func NewEncoder(input string) Encoder {
	return Encoder{
		InputImage: input,
	}
}

func (h Encoder) GetVideoSize() (int, int, error) {
	data, err := ffmpeg.ProbeWithTimeout(h.InputImage, time.Duration(time.Second*30), ffmpeg.KwArgs{"rtsp_transport": "tcp"})
	if err != nil {
		return 0, 0, err
	}

	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		return 0, 0, err
	}

	for _, s := range vInfo.Streams {
		if s.CodecType == "video" {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, errors.New("could not get video size")
}

func (h Encoder) ReadStream(ctx context.Context, stdout io.WriteCloser, stderr io.WriteCloser, fps string) chan error {
	var output ffmpeg.KwArgs
	if fps != "" {
		output = ffmpeg.KwArgs{"filter:v": "fps=" + fps, "format": "rawvideo", "pix_fmt": "rgb24"}
	} else {
		output = ffmpeg.KwArgs{"format": "rawvideo", "pix_fmt": "rgb24"}
	}
	return ffmpeg.Input(h.InputImage, ffmpeg.KwArgs{"rtsp_transport": "tcp"}).
		Output("pipe:", output).
		WithOutput(stdout).
		WithErrorOutput(stderr).
		RunCtx(ctx)
}

func (h Encoder) Catch(ctx context.Context, er io.Reader) chan error {
	err := make(chan error)
	go func() {
		for {
			buf := make([]byte, 1024)
			_, er := er.Read(buf)
			if er != nil {
				err <- er

				return
			}
			if strings.Contains(string(buf), "More than 1000 frames duplicated") {
				err <- errors.New(string(buf))

				return
			}

			select {
			case <-ctx.Done():
				err <- errors.New("context cancelled")
				return
			case <-time.After(1 * time.Second):
				// no-op
			}
		}
	}()
	return err
}
