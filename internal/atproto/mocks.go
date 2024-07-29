package atproto

import (
	"context"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

const (
	FakeExpiredToken = "eyJhbGciOiJFUzI1NksifQ.eyJzY29wZSI6ImNvbS5hdHByb3RvLmFwcFBhc3NQcml2aWxlZ2VkIiwic3ViIjoiZGlkOnBsYzpicXR2bjdteHNneG1lM2ZwejV4cGRhYXIiLCJpYXQiOjE3MjIxOTk3ODksImV4cCI6MTcyMjIwNjk4OSwiYXVkIjoiZGlkOndlYjpveXN0ZXIudXMtZWFzdC5ob3N0LmJza3kubmV0d29yayJ9.yIVsDiWhi0P3pjJrD4Ymb_HsoApWygsVLUZ9Ai03BqtZU5M6XOGg3YknPqlf2tbXng5Bhil2eEKhFgY2oRUSsw"
	FakeRefreshToken = "eyJhbGciOiJFUzI1NksifQ.eyJzY29wZSI6ImNvbS5hdHByb3RvLmFwcFBhc3NQcml2aWxlZ2VkIiwic3ViIjoiZGlkOnBsYzpicXR2bjdteHNneG1lM2ZwejV4cGRhYXIiLCJpYXQiOjE3MjIxOTk3ODksImV4cCI6MTcyMjIwNjk4OSwiYXVkIjoiZGlkOndlYjpveXN0ZXIudXMtZWFzdC5ob3N0LmJza3kubmV0d29yayJ9.yIVsDiWhi0P3pjJrD4Ymb_HsoApWygsVLUZ9Ai03BqtZU5M6XOGg3YknPqlf2tbXng5Bhil2eEKhFgY2oRUSsw"
	FakeDid          = "did:plc:abc222"
	FakeUri          = "at://did:plc:abc222"
	FakeCid          = "bafyreihmuibbe4u7jgti6kqcwvzzdewfxiacxe5cyimj5pnzvwn6arozj4"
)

type Mock struct {
	*Impl
}

func (a *Mock) ServerCreateSession(ctx context.Context, c *xrpc.Client, input *atproto.ServerCreateSession_Input) (*atproto.ServerCreateSession_Output, error) {
	t := true
	return &atproto.ServerCreateSession_Output{
		AccessJwt:  FakeExpiredToken,
		Active:     &t,
		Did:        FakeDid,
		Handle:     input.Identifier,
		RefreshJwt: FakeRefreshToken,
	}, nil
}

func (a *Mock) ServerRefreshSession(ctx context.Context, c *xrpc.Client) (*atproto.ServerRefreshSession_Output, error) {
	return &atproto.ServerRefreshSession_Output{
		AccessJwt: FakeExpiredToken,
	}, nil
}

func (a *Mock) RepoCreateRecord(ctx context.Context, c *xrpc.Client, input *atproto.RepoCreateRecord_Input) (*atproto.RepoCreateRecord_Output, error) {
	return &atproto.RepoCreateRecord_Output{
		Cid: FakeCid,
		Uri: FakeUri,
	}, nil
}
