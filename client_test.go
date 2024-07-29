package bskygo

import (
	"context"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/varsotech/bskygo/internal/atproto"
	"testing"
)

const (
	fakeHandle       = "test.example.com"
	fakePassword     = "password"
	fakeExpiredToken = "eyJhbGciOiJFUzI1NksifQ.eyJzY29wZSI6ImNvbS5hdHByb3RvLmFwcFBhc3NQcml2aWxlZ2VkIiwic3ViIjoiZGlkOnBsYzpicXR2bjdteHNneG1lM2ZwejV4cGRhYXIiLCJpYXQiOjE3MjIxOTk3ODksImV4cCI6MTcyMjIwNjk4OSwiYXVkIjoiZGlkOndlYjpveXN0ZXIudXMtZWFzdC5ob3N0LmJza3kubmV0d29yayJ9.yIVsDiWhi0P3pjJrD4Ymb_HsoApWygsVLUZ9Ai03BqtZU5M6XOGg3YknPqlf2tbXng5Bhil2eEKhFgY2oRUSsw"
	fakeRefreshToken = "eyJhbGciOiJFUzI1NksifQ.eyJzY29wZSI6ImNvbS5hdHByb3RvLmFwcFBhc3NQcml2aWxlZ2VkIiwic3ViIjoiZGlkOnBsYzpicXR2bjdteHNneG1lM2ZwejV4cGRhYXIiLCJpYXQiOjE3MjIxOTk3ODksImV4cCI6MTcyMjIwNjk4OSwiYXVkIjoiZGlkOndlYjpveXN0ZXIudXMtZWFzdC5ob3N0LmJza3kubmV0d29yayJ9.yIVsDiWhi0P3pjJrD4Ymb_HsoApWygsVLUZ9Ai03BqtZU5M6XOGg3YknPqlf2tbXng5Bhil2eEKhFgY2oRUSsw"
	fakeDid          = "did:plc:abc222"
	fakeUri          = "at://did:plc:abc222"
	fakeCid          = "bafyreihmuibbe4u7jgti6kqcwvzzdewfxiacxe5cyimj5pnzvwn6arozj4"
)

type ATProtoMock struct {
	*atproto.Impl
}

func (a *ATProtoMock) ServerCreateSession(ctx context.Context, c *xrpc.Client, input *atproto.ServerCreateSession_Input) (*atproto.ServerCreateSession_Output, error) {
	t := true
	return &atproto.ServerCreateSession_Output{
		AccessJwt:  fakeExpiredToken,
		Active:     &t,
		Did:        fakeDid,
		Handle:     input.Identifier,
		RefreshJwt: fakeRefreshToken,
	}, nil
}

func (a *ATProtoMock) ServerRefreshSession(ctx context.Context, c *xrpc.Client) (*atproto.ServerRefreshSession_Output, error) {
	return &atproto.ServerRefreshSession_Output{
		AccessJwt: fakeExpiredToken,
	}, nil
}

func (a *ATProtoMock) RepoCreateRecord(ctx context.Context, c *xrpc.Client, input *atproto.RepoCreateRecord_Input) (*atproto.RepoCreateRecord_Output, error) {
	return &atproto.RepoCreateRecord_Output{
		Cid: fakeCid,
		Uri: fakeUri,
	}, nil
}

func TestClient_Connect(t *testing.T) {
	client := NewClient(fakeHandle, fakePassword)
	client.atprotoClient = &ATProtoMock{}

	err := client.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if client.client.Auth == nil {
		t.Fatal("client.Auth is nil")
	}

	if client.client.Auth.Did != fakeDid {
		t.Fatalf("client.Auth.Did != %s, is: %s", fakeDid, client.client.Auth.Did)
	}

	if client.client.Auth.Handle != fakeHandle {
		t.Fatalf("client.Auth.Handle != %s, is: %s", fakeHandle, client.client.Auth.Handle)
	}

	if client.client.Auth.AccessJwt != fakeExpiredToken {
		t.Fatalf("client.Auth.AccessJwt != %s, is: %s", fakeExpiredToken, client.client.Auth.AccessJwt)
	}

	if client.client.Auth.RefreshJwt != fakeRefreshToken {
		t.Fatalf("client.client.Auth.RefreshJwt != %s, is: %s", fakeRefreshToken, client.client.Auth.RefreshJwt)
	}
}

func TestClientTokenRefresh(t *testing.T) {
	client := NewClient(fakeHandle, fakePassword)
	client.atprotoClient = &ATProtoMock{}

	err := client.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateFeedPost(context.Background(), NewFeedPost("text"))
	if err == nil {
		t.Fatal("expected unauthorized error creating post")
	}
}
