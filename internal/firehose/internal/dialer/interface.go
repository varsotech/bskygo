package dialer

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Dialer interface {
	Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}
