package health

import (
	"bytes"
	"errors"
	"net"
	"net/http"
	"regexp"
)

type Health struct {
	Url  string
	Host string
}

func NewHealth(url string) (*Health, error) {
	s := regexp.MustCompile(`^[a-z]+://([^:]+):([0-9]+)`).FindStringSubmatch(url)
	if len(s) != 3 {
		return nil, errors.New("invalid url")
	}
	return &Health{
		Url:  url,
		Host: s[1] + ":" + s[2],
	}, nil
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := net.Dial("tcp", h.Host)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer c.Close()

	_, err = c.Write([]byte("OPTIONS * RTSP/1.0\r\nCSeq: 1\r\n\r\n"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	buff := []byte{}
	for {
		message := make([]byte, 1024)
		n, err := c.Read(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		buff = append(buff, message[:n]...)
		if bytes.Contains(message[:n], []byte("\r\n\r\n")) {
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buff)
}
