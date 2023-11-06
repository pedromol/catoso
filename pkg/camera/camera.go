package camera

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	UP    = "UP"
	DOWN  = "DWON"
)

type move struct {
	direction string
	quantity  int
}

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
	log.Println("moving " + d)
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

func (h Camera) Centralize(moves string) error {
	mm := strings.Split(moves, ",")
	mvs := make([]move, 0)
	for _, m := range mm {
		nm := strings.Split(m, "=")
		if len(nm) < 2 {
			return errors.New(m + " is an invalid move")
		}
		qtd, err := strconv.ParseInt(nm[1], 10, 0)
		if err != nil {
			return errors.New(m + " has an invalid quantity")
		}
		mvs = append(mvs, move{
			direction: nm[0],
			quantity:  int(qtd),
		})
	}

	for _, m := range mvs {
		for i := 0; i < m.quantity; i++ {
			if err := h.Move(m.direction); err != nil {
				return err
			}
			time.Sleep(time.Duration(time.Millisecond * 500))
		}
	}

	return nil
}
