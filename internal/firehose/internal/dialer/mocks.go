package dialer

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type DialerSanityMock struct {
}

func (d *DialerSanityMock) Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
	return nil, nil, nil
}
