package bskygo

import "github.com/varsotech/bskygo/internal/log"

type ClientOption = func(*Client)

// WithHost sets the PDS host the client should connect to
func WithHost(host string) ClientOption {
	return func(client *Client) {
		client.host = host
	}
}

// WithFirehoseRetryOnReset sets whether firehose connection should be retried on connection reset. Defaults to true.
func WithFirehoseRetryOnReset(retry bool) ClientOption {
	return func(client *Client) {
		client.firehoseRetryOnReset = retry
	}
}

// WithLogger sets the logger the client uses
func WithLogger(logger log.Logger) ClientOption {
	return func(client *Client) {
		client.logger = logger
	}
}
