package firehose

import (
	"context"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/events"
	"golang.org/x/sync/errgroup"
	"sync"
)

type RepoStreamCallbacks struct {
	repoCommit   []func(evt *atproto.SyncSubscribeRepos_Commit) error
	repoCommitMu sync.RWMutex

	repoHandle   []func(evt *atproto.SyncSubscribeRepos_Handle) error
	repoHandleMu sync.RWMutex

	repoIdentity   []func(evt *atproto.SyncSubscribeRepos_Identity) error
	repoIdentityMu sync.RWMutex

	repoAccount   []func(evt *atproto.SyncSubscribeRepos_Account) error
	repoAccountMu sync.RWMutex

	repoInfo   []func(evt *atproto.SyncSubscribeRepos_Info) error
	repoInfoMu sync.RWMutex

	repoMigrate   []func(evt *atproto.SyncSubscribeRepos_Migrate) error
	repoMigrateMu sync.RWMutex

	repoTombstone   []func(evt *atproto.SyncSubscribeRepos_Tombstone) error
	repoTombstoneMu sync.RWMutex

	labelLabels   []func(evt *atproto.LabelSubscribeLabels_Labels) error
	labelLabelsMu sync.RWMutex

	labelInfo   []func(evt *atproto.LabelSubscribeLabels_Info) error
	labelInfoMu sync.RWMutex

	error   []func(evt *events.ErrorFrame) error
	errorMu sync.RWMutex

	// repoStreamCallbacks contains the callbacks provided externally. These can be read freely as
	// we do not write to them after the constructor.
	repoStreamCallbacks *events.RepoStreamCallbacks
}

func newRepoStreamCallbacks() *RepoStreamCallbacks {
	c := &RepoStreamCallbacks{}

	c.repoStreamCallbacks = &events.RepoStreamCallbacks{
		RepoCommit: func(evt *atproto.SyncSubscribeRepos_Commit) error {
			errGroup := getErrGroup(&c.repoCommitMu, evt, c.repoCommit)
			return errGroup.Wait()
		},
		RepoHandle: func(evt *atproto.SyncSubscribeRepos_Handle) error {
			errGroup := getErrGroup(&c.repoHandleMu, evt, c.repoHandle)
			return errGroup.Wait()
		},
		RepoIdentity: func(evt *atproto.SyncSubscribeRepos_Identity) error {
			errGroup := getErrGroup(&c.repoIdentityMu, evt, c.repoIdentity)
			return errGroup.Wait()
		},
		RepoAccount: func(evt *atproto.SyncSubscribeRepos_Account) error {
			errGroup := getErrGroup(&c.repoAccountMu, evt, c.repoAccount)
			return errGroup.Wait()
		},
		RepoInfo: func(evt *atproto.SyncSubscribeRepos_Info) error {
			errGroup := getErrGroup(&c.repoInfoMu, evt, c.repoInfo)
			return errGroup.Wait()
		},
		RepoMigrate: func(evt *atproto.SyncSubscribeRepos_Migrate) error {
			errGroup := getErrGroup(&c.repoMigrateMu, evt, c.repoMigrate)
			return errGroup.Wait()
		},
		RepoTombstone: func(evt *atproto.SyncSubscribeRepos_Tombstone) error {
			errGroup := getErrGroup(&c.repoTombstoneMu, evt, c.repoTombstone)
			return errGroup.Wait()
		},
		LabelLabels: func(evt *atproto.LabelSubscribeLabels_Labels) error {
			errGroup := getErrGroup(&c.labelLabelsMu, evt, c.labelLabels)
			return errGroup.Wait()
		},
		LabelInfo: func(evt *atproto.LabelSubscribeLabels_Info) error {
			errGroup := getErrGroup(&c.labelInfoMu, evt, c.labelInfo)
			return errGroup.Wait()
		},
		Error: func(evt *events.ErrorFrame) error {
			errGroup := getErrGroup(&c.errorMu, evt, c.error)
			return errGroup.Wait()
		},
	}

	return c
}

func getErrGroup[T any](mutex *sync.RWMutex, evt T, funcs []func(evt T) error) *errgroup.Group {
	mutex.RLock()
	defer mutex.RUnlock()

	errGroup := &errgroup.Group{}
	for _, f := range funcs {
		errGroup.Go(func() error {
			return f(evt)
		})
	}

	return errGroup
}

func (c *RepoStreamCallbacks) GetEventHandler() func(ctx context.Context, xev *events.XRPCStreamEvent) error {
	return c.repoStreamCallbacks.EventHandler
}

func (c *RepoStreamCallbacks) OnRepoCommit(f func(evt *atproto.SyncSubscribeRepos_Commit) error) {
	c.repoCommitMu.Lock()
	defer c.repoCommitMu.Unlock()
	c.repoCommit = append(c.repoCommit, f)
}

func (c *RepoStreamCallbacks) OnRepoHandle(f func(evt *atproto.SyncSubscribeRepos_Handle) error) {
	c.repoHandleMu.Lock()
	defer c.repoHandleMu.Unlock()
	c.repoHandle = append(c.repoHandle, f)
}

func (c *RepoStreamCallbacks) OnRepoIdentity(f func(evt *atproto.SyncSubscribeRepos_Identity) error) {
	c.repoIdentityMu.Lock()
	defer c.repoIdentityMu.Unlock()
	c.repoIdentity = append(c.repoIdentity, f)
}

func (c *RepoStreamCallbacks) OnRepoAccount(f func(evt *atproto.SyncSubscribeRepos_Account) error) {
	c.repoAccountMu.Lock()
	defer c.repoAccountMu.Unlock()
	c.repoAccount = append(c.repoAccount, f)
}

func (c *RepoStreamCallbacks) OnRepoInfo(f func(evt *atproto.SyncSubscribeRepos_Info) error) {
	c.repoInfoMu.Lock()
	defer c.repoInfoMu.Unlock()
	c.repoInfo = append(c.repoInfo, f)
}

func (c *RepoStreamCallbacks) OnRepoMigrate(f func(evt *atproto.SyncSubscribeRepos_Migrate) error) {
	c.repoMigrateMu.Lock()
	defer c.repoMigrateMu.Unlock()
	c.repoMigrate = append(c.repoMigrate, f)
}

func (c *RepoStreamCallbacks) OnRepoTombstone(f func(evt *atproto.SyncSubscribeRepos_Tombstone) error) {
	c.repoTombstoneMu.Lock()
	defer c.repoTombstoneMu.Unlock()
	c.repoTombstone = append(c.repoTombstone, f)
}

func (c *RepoStreamCallbacks) OnLabelLabels(f func(evt *atproto.LabelSubscribeLabels_Labels) error) {
	c.labelLabelsMu.Lock()
	defer c.labelLabelsMu.Unlock()
	c.labelLabels = append(c.labelLabels, f)
}

func (c *RepoStreamCallbacks) OnError(f func(evt *events.ErrorFrame) error) {
	c.errorMu.Lock()
	defer c.errorMu.Unlock()
	c.error = append(c.error, f)
}
