package vision

import (
	"context"
	"errors"
	"image"
	"image/color"

	"io"
	"log"
	"time"

	"gocv.io/x/gocv"
)

const Catoso = "Catoso"

type Vision struct {
	XmlFile             string
	Width               int
	Height              int
	DelayAfterDetectMin int64
	Frameskip           int
	Debug               bool
	DrawOverFace        bool
}

func NewVision(xmlPath string, w int, h int, delay int64, fskip int, debug bool, draw bool) *Vision {
	return &Vision{
		XmlFile:             xmlPath,
		Width:               w,
		Height:              h,
		DelayAfterDetectMin: delay,
		Frameskip:           fskip,
		Debug:               debug,
		DrawOverFace:        draw,
	}
}

func (v *Vision) Process(ctx context.Context, reader io.ReadCloser, stream *Stream) (chan []byte, chan error) {
	result := make(chan error)
	imgchan := make(chan []byte)
	var win *gocv.Window
	if v.Debug {
		win = gocv.NewWindow(Catoso)
		defer win.Close()
	}
	go func() {
		detected := 0
		lastConfirmed := time.Date(1, time.January, 1, 1, 1, 1, 1, time.Local)

		blue := color.RGBA{B: 255}

		classifier := gocv.NewCascadeClassifier()
		defer classifier.Close()

		if !classifier.Load(v.XmlFile) {
			result <- errors.New("error reading cascade file: " + v.XmlFile)
			return
		}

		frameSize := v.Width * v.Height * 3
		buf := make([]byte, frameSize)
		cf := 0
		for {
			select {
			case <-ctx.Done():
				result <- nil
				return
			default:
				// no-op
			}

			n, err := io.ReadFull(reader, buf)
			if n == 0 || err == io.EOF {
				result <- errors.New("EOF")
				break
			} else if n != frameSize || err != nil {
				result <- err
				break
			}

			cf += 1
			if cf < v.Frameskip {
				continue
			}
			cf = 0

			if lastConfirmed.After(time.Now()) && stream == nil {
				continue
			}

			img, err := gocv.NewMatFromBytes(v.Height, v.Width, gocv.MatTypeCV8UC3, buf)
			if err != nil {
				log.Println("NewMatFromBytes error: " + err.Error())
				continue
			}

			if img.Empty() {
				log.Println("img.Empty")
				img.Close()
				continue
			}

			img2 := gocv.NewMat()

			gocv.CvtColor(img, &img2, gocv.ColorBGRToRGB)
			img.Close()

			if stream != nil {
				buf, err := gocv.IMEncode(gocv.JPEGFileExt, img2)
				if err == nil {
					stream.UpdateJPEG(buf.GetBytes())
					buf.Close()
				}
			}

			if lastConfirmed.After(time.Now()) {
				img2.Close()
				continue
			}

			rects := classifier.DetectMultiScale(img2)
			if len(rects) > 0 {
				detected += 1
			} else {
				detected = 0
			}

			if detected > 3 {
				log.Println(Catoso)
				lastConfirmed = time.Now().Add(time.Minute * time.Duration(v.DelayAfterDetectMin))
				if v.DrawOverFace {
					for _, r := range rects {
						gocv.Rectangle(&img2, r, blue, 3)
						size := gocv.GetTextSize(Catoso, gocv.FontHersheyPlain, 1.2, 2)
						pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
						gocv.PutText(&img2, Catoso, pt, gocv.FontHersheyPlain, 1.2, blue, 2)
					}
				}
				buff, err := gocv.IMEncode(gocv.JPEGFileExt, img2)
				if err != nil {
					img2.Close()
					result <- err
					return
				}

				imgchan <- buff.GetBytes()
				buff.Close()
			}
			if win != nil {
				win.IMShow(img2)
				win.WaitKey(10)
			}

			img2.Close()
		}
	}()
	return imgchan, result
}
