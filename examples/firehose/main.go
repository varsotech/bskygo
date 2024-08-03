package main

import (
	"context"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/varsotech/bskygo"
	"log"
	"os"
)

func main() {
	username := os.Getenv("BSKY_USERNAME")
	password := os.Getenv("BSKY_PASSWORD")

	client := bskygo.NewClient(username, password)

	client.Firehose().OnRepoCommit(func(evt *atproto.SyncSubscribeRepos_Commit) error {
		log.Println(evt.Time)
		return nil
	})

	ctx := context.Background()
	err := client.ConnectAndListen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
