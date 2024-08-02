package main

import (
	"context"
	"github.com/varsotech/bskygo"
	"log"
	"os"
)

func main() {
	username := os.Getenv("BSKY_USERNAME")
	password := os.Getenv("BSKY_PASSWORD")

	client := bskygo.NewClient(username, password)

	ctx := context.Background()
	closer, err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	handle := "varso.org"
	handleDid, err := client.GetHandleDid(ctx, handle)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.FeedCreatePost(ctx, bskygo.NewFeedPost("Hello world! ").Mention(handle, handleDid))
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(post.Cid)
}