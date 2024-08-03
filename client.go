package bskygo

import (
	"context"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/varsotech/bskygo/internal/atproto"
	"github.com/varsotech/bskygo/internal/firehose"
	"github.com/varsotech/bskygo/internal/xrpc"
	"golang.org/x/sync/errgroup"
)

const defaultHost = "https://bsky.social"

type Client struct {
	xrpcClient     *xrpc.Client
	atprotoClient  atproto.ATProto
	firehoseClient *firehose.Firehose

	username             string
	password             string
	host                 string
	firehoseRetryOnReset bool
}

// NewClient creates a new Client instance. Call the Connect method on the returned client instance
// to create an authenticated session.
func NewClient(username, password string, options ...ClientOption) *Client {
	c := &Client{
		xrpcClient:           nil,
		atprotoClient:        atproto.New(),
		firehoseClient:       firehose.New(),
		username:             username,
		password:             password,
		host:                 defaultHost,
		firehoseRetryOnReset: true,
	}

	for _, option := range options {
		option(c)
	}

	c.xrpcClient = xrpc.New(c.host)
	return c
}

// Connect establishes an authenticated session with the server.
func (c *Client) Connect(ctx context.Context) error {
	sessionDetails, err := c.createSession(ctx)
	if err != nil {
		return err
	}

	c.xrpcClient.UpdateAuth(sessionDetails.AccessJwt, sessionDetails.RefreshJwt, sessionDetails.Handle, sessionDetails.Did)
	return nil
}

// ConnectAndListen establishes an authenticated session with the server, and listens to firehose events.
// This function is blocking.
func (c *Client) ConnectAndListen(ctx context.Context) error {
	err := c.Connect(ctx)
	if err != nil {
		return err
	}

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return c.firehoseClient.ConnectAndListen(ctx, c.firehoseRetryOnReset)
	})

	return errGroup.Wait()
}

func (c *Client) createSession(ctx context.Context) (sessionDetails *atproto.ServerCreateSession_Output, err error) {
	sessionInput := &atproto.ServerCreateSession_Input{
		Identifier: c.username,
		Password:   c.password,
	}

	_, err = c.xrpcClient.UseWithoutRefresh(func(client *xrpc.XRPC) (err error) {
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

// Firehose exposes methods to subscribe to network events.
// Callbacks should return quickly to avoid connection being reset by the server.
func (c *Client) Firehose() *firehose.Firehose {
	return c.firehoseClient
}
