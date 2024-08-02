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

func (p *FeedPost) Mention(handle, did string) *FeedPost {
	byteStart := len(p.record.Text)
	p.record.Text += "@" + handle
	p.record.Facets = append(p.record.Facets, &bsky.RichtextFacet{
		Features: []*bsky.RichtextFacet_Features_Elem{
			{
				RichtextFacet_Mention: &bsky.RichtextFacet_Mention{
					LexiconTypeID: "app.bsky.richtext.facet#mention",
					Did:           did,
				},
			},
		},
		Index: &bsky.RichtextFacet_ByteSlice{
			ByteStart: int64(byteStart),
			ByteEnd:   int64(len(p.record.Text)),
		},
	})

	return p
}
