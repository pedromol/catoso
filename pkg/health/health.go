package health

import (
	"bytes"
	"net"
	"net/http"
)

type Health struct {
	Host string
}

func NewHealth(url string) (*Health, error) {
	return &Health{
		Host: url,
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
