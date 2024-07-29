package atproto

import (
	"context"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type Impl struct {
}

func New() *Impl {
	return &Impl{}
}

func (a *Impl) ServerCreateSession(ctx context.Context, c *xrpc.Client, input *ServerCreateSession_Input) (*ServerCreateSession_Output, error) {
	return atproto.ServerCreateSession(ctx, c, input)
}

func (a *Impl) RepoCreateRecord(ctx context.Context, c *xrpc.Client, input *RepoCreateRecord_Input) (*RepoCreateRecord_Output, error) {
	return atproto.RepoCreateRecord(ctx, c, input)
}

func (a *Impl) ServerRefreshSession(ctx context.Context, c *xrpc.Client) (*ServerRefreshSession_Output, error) {
	return atproto.ServerRefreshSession(ctx, c)
}
