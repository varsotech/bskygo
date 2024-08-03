package events

import (
	"context"
	"github.com/gorilla/websocket"
)

type Events interface {
	HandleRepoStream(ctx context.Context, con *websocket.Conn, sched Scheduler) error
}
