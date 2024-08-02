package bskygo

import (
	"context"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/varsotech/bskygo/internal/atproto"
	"github.com/varsotech/bskygo/internal/xrpc"
)

const defaultHost = "https://bsky.social"

type Client struct {
	xrpcClient    *xrpc.Client
	atprotoClient atproto.ATProto

	username string
	password string
	host     string
}

func NewClient(username, password string, options ...Option) *Client {
	c := &Client{
		atprotoClient: atproto.New(),
		username:      username,
		password:      password,
		host:          defaultHost,
	}

	for _, option := range options {
		option(c)
	}

	c.xrpcClient = xrpc.New(c.host)
	return c
}

// Connect establishes a session with the server.
func (c *Client) Connect(ctx context.Context) (func(), error) {
	sessionDetails, err := c.createSession(ctx)
	if err != nil {
		return nil, err
	}

	c.xrpcClient.UpdateAuth(sessionDetails.AccessJwt, sessionDetails.RefreshJwt, sessionDetails.Handle, sessionDetails.Did)

	// Start routines
	ctx, cancel := context.WithCancel(ctx)

	return func() {
		cancel()
	}, nil
}

func (c *Client) createSession(ctx context.Context) (sessionDetails *atproto.ServerCreateSession_Output, err error) {
	sessionInput := &atproto.ServerCreateSession_Input{
		Identifier: c.username,
		Password:   c.password,
	}

	err = c.xrpcClient.Use(ctx, c.atprotoClient, func(client *xrpc.XRPC) (err error) {
		sessionDetails, err = c.atprotoClient.ServerCreateSession(ctx, client, sessionInput)
		return
	})

	return
}

type CreateFeedPostOutput struct {
	Cid string
	Uri string
}

func (c *Client) FeedCreatePost(ctx context.Context, post *FeedPost) (*CreateFeedPostOutput, error) {
	var response *atproto.RepoCreateRecord_Output

	err := c.xrpcClient.Use(ctx, c.atprotoClient, func(xrpc *xrpc.XRPC) (err error) {
		createRecordInput := &atproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       xrpc.Auth.Did,
			Record:     &lexutil.LexiconTypeDecoder{Val: post.record},
		}

		response, err = c.atprotoClient.RepoCreateRecord(ctx, xrpc, createRecordInput)
		return
	})

	if err != nil {
		return nil, err
	}

	return &CreateFeedPostOutput{
		Cid: response.Cid,
		Uri: response.Uri,
	}, nil
}

func (c *Client) GetHandleDid(ctx context.Context, handle string) (string, error) {
	var response *atproto.IdentityResolveHandle_Output

	err := c.xrpcClient.Use(ctx, c.atprotoClient, func(xrpc *xrpc.XRPC) (err error) {
		response, err = c.atprotoClient.IdentityResolveHandle(ctx, xrpc, handle)
		return
	})

	if err != nil {
		return "", err
	}

	return response.Did, nil
}
