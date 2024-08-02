package atproto

import (
	"context"
	"github.com/bluesky-social/indigo/xrpc"
)

type ATProto interface {
	ServerCreateSession(ctx context.Context, c *xrpc.Client, input *ServerCreateSession_Input) (*ServerCreateSession_Output, error)
	RepoCreateRecord(ctx context.Context, c *xrpc.Client, input *RepoCreateRecord_Input) (*RepoCreateRecord_Output, error)
	ServerRefreshSession(ctx context.Context, c *xrpc.Client) (*ServerRefreshSession_Output, error)
	IdentityResolveHandle(ctx context.Context, c *xrpc.Client, handle string) (*IdentityResolveHandle_Output, error)
}
