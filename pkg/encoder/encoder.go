package encoder

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func (h Encoder) ReadStream(stdout io.WriteCloser, stderr io.WriteCloser) chan error {
	result := make(chan error)
	go func() {
		defer stdout.Close()
		defer stderr.Close()
		result <- ffmpeg.Input(h.InputImage, ffmpeg.KwArgs{"rtsp_transport": "tcp"}).
			Output("pipe:",
				ffmpeg.KwArgs{
					"format": "rawvideo", "pix_fmt": "rgb24",
				}).
			WithOutput(stdout).
			WithErrorOutput(stderr).
			Run()
	}()

	return result
}

func (h Encoder) Catch(er io.Reader) chan error {
	err := make(chan error)
	go func() {
		for {
			buf := make([]byte, 1024)
			er.Read(buf)
			if strings.Contains(string(buf), "More than 1000 frames duplicated") {
				err <- errors.New(string(buf))
				return
			}
		}
	}()
	return err
}
