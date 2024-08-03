package xrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/varsotech/bskygo/internal/atproto"
	"net/http"
	"sync"
	"time"
)

type XRPC = xrpc.Client

type Client struct {
	client *xrpc.Client

	// clientMutex protects client when auth token is refreshed
	clientMutex sync.RWMutex
}

func New(host string) *Client {
	return &Client{
		client: &xrpc.Client{
			Host: host,
		},
	}
}

// UseWithoutRefresh locks the client for reading and provides access to it.
// It also returns the access token that was used in the request.
// Unlike Use, it does not refresh the token and retries when getting unauthorized error.
func (c *Client) UseWithoutRefresh(f func(client *xrpc.Client) error) (string, error) {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()

	usedAccessToken := ""
	if c.client.Auth != nil {
		usedAccessToken = c.client.Auth.AccessJwt
	}

	return usedAccessToken, f(c.client)
}

// Use locks the client for reading and provides access to it.
// If it gets an unauthorized error, it refreshes the token and retries
func (c *Client) Use(ctx context.Context, atprotoClient atproto.ATProto, f func(client *xrpc.Client) error) error {
	usedAccessToken, err := c.UseWithoutRefresh(f)
	if err == nil {
		return nil
	}

	var xrpcErr *xrpc.Error
	if !errors.As(err, &xrpcErr) {
		return err
	}

	if xrpcErr.StatusCode != http.StatusUnauthorized {
		return err
	}

	refreshErr := c.refreshToken(ctx, atprotoClient, usedAccessToken)
	if refreshErr != nil {
		return fmt.Errorf("failed to refresh the token: %w, initial error was: %w", refreshErr, err)
	}

	// Retry after having refreshed the token
	_, err = c.UseWithoutRefresh(f)
	return err
}

func (c *Client) refreshToken(ctx context.Context, atprotoClient atproto.ATProto, usedAccessToken string) error {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	if c.client.Auth == nil {
		return fmt.Errorf("client auth has not been initialized yet")
	}

	currentToken := c.client.Auth.AccessJwt
	if currentToken != usedAccessToken {
		// Another routine already refreshed the token
		return nil
	}

	sessionDetails, err := atprotoClient.ServerRefreshSession(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed refreshing session: %w", err)
	}

	c.updateAuthWithoutLock(sessionDetails.AccessJwt, sessionDetails.RefreshJwt, sessionDetails.Handle, sessionDetails.Did)
	return nil
}

func isAccessTokenExpired(token string) bool {
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

func (c *Client) updateAuthWithoutLock(accessToken, refreshToken, handle, did string) {
	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  accessToken,
		RefreshJwt: refreshToken,
		Handle:     handle,
		Did:        did,
	}
}

func (c *Client) UpdateAuth(accessToken, refreshToken, handle, did string) {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	c.updateAuthWithoutLock(accessToken, refreshToken, handle, did)
}
