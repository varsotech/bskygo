package bskygo

import (
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util"
	"time"
)

type FeedPost struct {
	record *bsky.FeedPost
}

func NewFeedPost(text string) *FeedPost {
	return &FeedPost{
		record: &bsky.FeedPost{
			LexiconTypeID: "app.bsky.feed.post",
			CreatedAt:     time.Now().Format(util.ISO8601),
			Text:          text,
		},
	}
}

func NewFeedPostFromRecord(post *bsky.FeedPost) *FeedPost {
	return &FeedPost{
		record: post,
	}
}
