package bskygo

import (
	"context"
	"errors"
	"fmt"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/varsotech/bskygo/internal/atproto"
	"net/http"
	"sync"
	"time"
)

const defaultHost = "https://bsky.social"

type Client struct {
	client      *xrpc.Client
	clientMutex sync.RWMutex

	atprotoClient atproto.ATProto

	username string
	password string
	host     string
}

func NewClient(username, password string, options ...Option) *Client {
	c := &Client{
		username: username,
		password: password,
		host:     defaultHost,
	}

	for _, option := range options {
		option(c)
	}

	c.client = &xrpc.Client{
		Host: c.host,
	}

	return c
}

// Connect establishes a session with the server.
func (c *Client) Connect(ctx context.Context) error {
	sessionInput := &atproto.ServerCreateSession_Input{
		Identifier: c.username,
		Password:   c.password,
	}

	sessionDetails, err := c.atprotoClient.ServerCreateSession(ctx, c.client, sessionInput)
	if err != nil {
		return fmt.Errorf("failed creating session: %w", err)
	}

	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  sessionDetails.AccessJwt,
		RefreshJwt: sessionDetails.RefreshJwt,
		Handle:     sessionDetails.Handle,
		Did:        sessionDetails.Did,
	}

	return nil
}

func (c *Client) GetXRPCClient() *xrpc.Client {
	return c.client
}

func (c *Client) shouldRetryWithRefreshedToken(ctx context.Context, err *error) bool {
	if err == nil {
		return false
	}

	var xrpcErr *xrpc.Error
	if !errors.As(*err, &xrpcErr) {
		return false
	}

	if xrpcErr.StatusCode != http.StatusUnauthorized {
		return false
	}

	refreshErr := c.refreshToken(ctx)
	if refreshErr != nil {
		*err = refreshErr
		return false
	}

	return true
}

func (c *Client) isAccessTokenExpired(token string) bool {
	claims := jwt.MapClaims{}
	jwtToken, _, _ := jwt.NewParser().ParseUnverified(token, &claims)
	if jwtToken == nil {
		return true
	}

	expiration, err := claims.GetExpirationTime()
	if err != nil {
		return true
	}

	if expiration.UTC().Before(time.Now()) {
		return true
	}

	return false
}

func (c *Client) refreshToken(ctx context.Context) error {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	// Access token without read lock, as we are write locking
	currentToken := c.client.Auth.AccessJwt
	if c.isAccessTokenExpired(currentToken) {
		// Another routine already refreshed the token
		return nil
	}

	sessionDetails, err := c.atprotoClient.ServerRefreshSession(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed refreshing session: %w", err)
	}

	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  sessionDetails.AccessJwt,
		RefreshJwt: sessionDetails.RefreshJwt,
		Handle:     sessionDetails.Handle,
		Did:        sessionDetails.Did,
	}

	return nil
}

func (c *Client) makeRLockedRequest(f func() error) error {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()
	return f()
}

// makeAuthenticatedRequest attempts making the request with a read lock on the session. If 401 is returned
// it will write lock the session and refresh the token, unless another routine already did.
func (c *Client) makeAuthenticatedRequest(ctx context.Context, f func() error) error {
	err := c.makeRLockedRequest(f)
	if c.shouldRetryWithRefreshedToken(ctx, &err) {
		return f()
	}
	return err
}

type CreateFeedPostOutput struct {
	Cid string
	Uri string
}

func (c *Client) CreateFeedPost(ctx context.Context, post *FeedPost) (*CreateFeedPostOutput, error) {
	createRecordInput := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       c.client.Auth.Did,
		Record:     &lexutil.LexiconTypeDecoder{Val: post.record},
	}

	var response *atproto.RepoCreateRecord_Output
	err := c.makeAuthenticatedRequest(ctx, func() (err error) {
		response, err = c.atprotoClient.RepoCreateRecord(ctx, c.client, createRecordInput)
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
