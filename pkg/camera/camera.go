package camera

import (
	"net"
	"time"
)

const (
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	UP    = "UP"
	DOWN  = "DWON"
)

type Camera struct {
	OnvifIP   string
	OnvifPort string
}

func NewCamera(ip string, port string) Camera {
	return Camera{
		OnvifIP:   ip,
		OnvifPort: port,
	}
}

func (h Camera) Move(d string) error {
	c, err := net.Dial("tcp", h.OnvifIP+":"+h.OnvifPort)
	if err != nil {
		return err
	}
	defer c.Close()

	c.SetDeadline(time.Now().Add(time.Duration(time.Second * 5)))

	_, err = c.Write([]byte(`
	SET_PARAMETER rtsp://` + h.OnvifIP + `/onvif1 RTSP/1.0
	CSeq: 1
	Content-type: ptzCmd: ` + d + `
  `))

	return err
}

func (h Camera) Centralize() error {
	for i := 0; i <= 20; i++ {
		if err := h.Move(LEFT); err != nil {
			return err
		}

		time.Sleep(time.Duration(time.Millisecond * 1000))
	}
	for i := 0; i <= 20; i++ {
		if err := h.Move(DOWN); err != nil {
			return err
		}

		time.Sleep(time.Duration(time.Millisecond * 1000))
	}
	for i := 0; i <= 9; i++ {
		if err := h.Move(RIGHT); err != nil {
			return err
		}

		time.Sleep(time.Duration(time.Millisecond * 1000))
	}
	for i := 0; i <= 2; i++ {
		if err := h.Move(UP); err != nil {
			return err
		}

		time.Sleep(time.Duration(time.Millisecond * 1000))
	}

	return nil
}
