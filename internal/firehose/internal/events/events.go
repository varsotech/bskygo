package events

import (
	"context"
	"github.com/bluesky-social/indigo/events"
	"github.com/gorilla/websocket"
)

type Impl struct {
}

func New() *Impl {
	return &Impl{}
}

func (i *Impl) HandleRepoStream(ctx context.Context, con *websocket.Conn, sched Scheduler) error {
	return events.HandleRepoStream(ctx, con, sched)
}
