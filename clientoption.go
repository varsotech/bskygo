package bskygo

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
