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
	XmlFile string
	Width   int
	Height  int
}

func NewVision(xmlPath string, w int, h int) Vision {
	return Vision{
		XmlFile: xmlPath,
		Width:   w,
		Height:  h,
	}
}

func (v Vision) Process(ctx context.Context, reader io.ReadCloser, debug string) (chan []byte, chan error) {
	result := make(chan error)
	imgchan := make(chan []byte)
	var win *gocv.Window
	if debug != "" {
		win = gocv.NewWindow(Catoso)
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
		for {
			n, err := io.ReadFull(reader, buf)
			if n == 0 || err == io.EOF {
				result <- errors.New("EOF")
				break
			} else if n != frameSize || err != nil {
				result <- err
				break
			}

			if lastConfirmed.After(time.Now()) {
				continue
			}

			img, err := gocv.NewMatFromBytes(v.Height, v.Width, gocv.MatTypeCV8UC3, buf)
			if err != nil {
				continue
			}

			if img.Empty() {
				img.Close()
				continue
			}
			img2 := gocv.NewMat()

			gocv.CvtColor(img, &img2, gocv.ColorBGRToRGB)

			rects := classifier.DetectMultiScale(img2)
			if len(rects) > 0 {
				detected += 1
			} else {
				detected = 0
			}

			if detected > 3 {
				log.Println(Catoso)
				lastConfirmed = time.Now().Add(time.Duration(time.Minute * 2))
				for _, r := range rects {
					gocv.Rectangle(&img2, r, blue, 3)

					size := gocv.GetTextSize(Catoso, gocv.FontHersheyPlain, 1.2, 2)
					pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
					gocv.PutText(&img2, Catoso, pt, gocv.FontHersheyPlain, 1.2, blue, 2)
				}
				buff, err := gocv.IMEncode(gocv.JPEGFileExt, img2)
				if err != nil {
					result <- err
					return
				}

				imgchan <- buff.GetBytes()
			}
			if debug != "" {
				win.IMShow(img2)
				win.WaitKey(10)
			}

			img.Close()
			img2.Close()

			select {
			case <-ctx.Done():
				if debug != "" {
					win.Close()
				}
				result <- nil
				return
			default:
				// no-op
			}
		}
	}()
	return imgchan, result
}
