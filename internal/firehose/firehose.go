package firehose

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	"github.com/google/uuid"
	"github.com/varsotech/bskygo/internal/firehose/internal/dialer"
	"github.com/varsotech/bskygo/internal/firehose/internal/events"
	"github.com/varsotech/bskygo/internal/log"
	"net/http"
	"strings"
)

type Firehose struct {
	*RepoStreamCallbacks

	logger    log.Logger
	dialer    dialer.Dialer
	events    events.Events
	identity  string
	scheduler *sequential.Scheduler
}

// New creates a new Firehose instance
func New(logger log.Logger) *Firehose {
	identity := uuid.NewString()
	callbacks := newRepoStreamCallbacks()
	scheduler := sequential.NewScheduler(identity, callbacks.GetEventHandler())

	return &Firehose{
		logger:              logger,
		dialer:              dialer.New(),
		events:              events.New(),
		identity:            identity,
		RepoStreamCallbacks: callbacks,
		scheduler:           scheduler,
	}
}

// ConnectAndListen is a blocking function that listens to events from the bsky network and dispatches
// the configured handlers.
func (f *Firehose) ConnectAndListen(ctx context.Context, retryOnReset bool) error {
	for {
		uri := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
		conn, _, err := f.dialer.Dial(uri, http.Header{})
		if err != nil {
			return fmt.Errorf("websocket dial to bsky network failed: %w", err)
		}

		err = f.events.HandleRepoStream(ctx, conn, f.scheduler)
		if err == nil {
			return nil
		}

		if retryOnReset && strings.Contains(err.Error(), ": connection reset by peer") {
			f.logger.Error("firehose connection reset by peer, attempting reconnect. this may happen when event handlers take too long")
			continue
		}

		return fmt.Errorf("handle repo stream failed: %w", err)
	}
}
