package bskygo

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/varsotech/bskygo/internal/atproto"
	"net/http"
	"testing"
)

const (
	fakeHandle   = "test.example.com"
	fakePassword = "password"

	newAccessToken  = "new_access_token"
	newRefreshToken = "new_refresh_token"
)

func TestClient_Connect(t *testing.T) {
	client := NewClient(fakeHandle, fakePassword)
	client.atprotoClient = &atproto.Mock{}

	closer, err := client.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer closer()
}

type ClientTokenRefreshMock struct {
	atproto.Mock
	repoCreateRecordCalls int
}

func (m *ClientTokenRefreshMock) RepoCreateRecord(ctx context.Context, c *xrpc.Client, input *atproto.RepoCreateRecord_Input) (*atproto.RepoCreateRecord_Output, error) {
	m.repoCreateRecordCalls++

	if m.repoCreateRecordCalls == 1 {
		fmt.Println("repo create record called, returning unauthorized")
		return nil, &xrpc.Error{
			StatusCode: http.StatusUnauthorized,
			Wrapped:    fmt.Errorf("unauthorized"),
		}
	}

	fmt.Println("repo create record called, returning valid response")
	return &atproto.RepoCreateRecord_Output{
		Cid: atproto.FakeCid,
		Uri: atproto.FakeUri,
	}, nil
}

func (m *ClientTokenRefreshMock) ServerRefreshSession(ctx context.Context, c *xrpc.Client) (*atproto.ServerRefreshSession_Output, error) {
	fmt.Println("server refresh session called")
	return &atproto.ServerRefreshSession_Output{
		AccessJwt:  newAccessToken,
		RefreshJwt: newRefreshToken,
	}, nil
}

func TestClientTokenRefresh(t *testing.T) {
	client := NewClient(fakeHandle, fakePassword)
	client.atprotoClient = &ClientTokenRefreshMock{}

	closer, err := client.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer closer()
	
	_, err = client.FeedCreatePost(context.Background(), NewFeedPost("text"))
	if err != nil {
		t.Fatal(err)
	}

	if client.client.Auth.AccessJwt != newAccessToken {
		t.Fatal("Access jwt was not refreshed")
	}

	if client.client.Auth.RefreshJwt != newRefreshToken {
		t.Fatal("Refresh jwt was not refreshed")
	}
}
