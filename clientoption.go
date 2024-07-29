package bskygo

type Option = func(*Client)

// WithHost sets the PDS host the client should connect to
func WithHost(host string) Option {
	return func(client *Client) {
		client.host = host
	}
}
