package dialer

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Impl struct {
}

func New() *Impl {
	return &Impl{}
}

func (d *Impl) Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(urlStr, requestHeader)
}
