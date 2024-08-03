package firehose

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/varsotech/bskygo/internal/firehose/internal/dialer"
	"github.com/varsotech/bskygo/internal/firehose/internal/events"
	"github.com/varsotech/bskygo/internal/log"
	"strings"
	"testing"
	"time"
)

const (
	someFatalError       = "some fatal error"
	connectionResetError = "some error: connection reset by peer"
)

type EventsFatalErrorMock struct {
	*events.Impl
}

func (e *EventsFatalErrorMock) HandleRepoStream(ctx context.Context, con *websocket.Conn, sched events.Scheduler) error {
	return fmt.Errorf(someFatalError)
}

func TestFirehose_ConnectAndListen_FatalError(t *testing.T) {
	firehose := New(log.NewSlog())
	firehose.dialer = &dialer.DialerSanityMock{}
	firehose.events = &EventsFatalErrorMock{}

	err := firehose.ConnectAndListen(context.Background(), true)
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if !strings.Contains(err.Error(), someFatalError) {
		t.Fatal("expected error, got none")
	}
}

type EventsRetryOnResetMock struct {
	count int
	*events.Impl
}

func (e *EventsRetryOnResetMock) HandleRepoStream(ctx context.Context, con *websocket.Conn, sched events.Scheduler) error {
	e.count++
	if e.count == 1 {
		fmt.Println("call to HandleRepoStream, returning connection reset error")
		time.Sleep(1 * time.Second)
		return fmt.Errorf(connectionResetError)
	}

	fmt.Println("call to HandleRepoStream, returning nil")
	return nil
}

func TestFirehose_ConnectAndListen_RetryOnResetTrue(t *testing.T) {
	firehose := New(log.NewSlog())
	firehose.dialer = &dialer.DialerSanityMock{}
	firehose.events = &EventsRetryOnResetMock{}

	err := firehose.ConnectAndListen(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFirehose_ConnectAndListen_RetryOnResetFalse(t *testing.T) {
	firehose := New(log.NewSlog())
	firehose.dialer = &dialer.DialerSanityMock{}
	firehose.events = &EventsRetryOnResetMock{}

	err := firehose.ConnectAndListen(context.Background(), false)
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if !strings.Contains(err.Error(), connectionResetError) {
		t.Fatalf("expected connection reset error, got: %s", err.Error())
	}
}
