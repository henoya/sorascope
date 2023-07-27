package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/urfave/cli/v2"
	"os"
)

func doShowPosts(cCtx *cli.Context) error {
	if cCtx.Args().Present() {
		return cli.ShowSubcommandHelp(cCtx)
	}

	xrpcc, err := makeXRPCC(cCtx)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}

	var feed []*atproto.RepoListRecords_Record

	n := cCtx.Int64("n")
	handle := cCtx.String("handle")

	var cursor string
	for {
		if handle != "" {
			if handle == "self" {
				handle = xrpcc.Auth.Did
			}
			//func RepoListRecords(ctx context.Context, c *xrpc.Client, collection string, cursor string, limit int64, repo string, reverse bool, rkeyEnd string, rkeyStart string) (*RepoListRecords_Output, error) {
			var collection string = "app.bsky.feed.post"
			resp, err := atproto.RepoListRecords(context.TODO(), xrpcc, collection, cursor, n, handle, true, "", "")
			//resp, err := bsky.FeedGetAuthorFeed(context.TODO(), xrpcc, handle, cursor, n)
			if err != nil {
				return fmt.Errorf("cannot get author feed: %w", err)
			}
			feed = append(feed, resp.Records...)
			if resp.Cursor != nil {
				cursor = *resp.Cursor
			} else {
				cursor = ""
			}
		} else {
			var collection string = "app.bsky.feed.post"
			resp, err := atproto.RepoListRecords(context.TODO(), xrpcc, collection, cursor, n, handle, true, "", "")
			if err != nil {
				return fmt.Errorf("cannot get timeline: %w", err)
			}
			feed = append(feed, resp.Records...)
			if resp.Cursor != nil {
				cursor = *resp.Cursor
			} else {
				cursor = ""
			}
		}
		//if cursor == "" || int64(len(feed)) > n {
		//	break
		//}
		if cursor == "" {
			break
		}
	}

	//sort.Slice(feed, func(i, j int) bool {
	//	ri := timep(feed[i].Value.Post.Record.Val.(*bsky.FeedPost).CreatedAt)
	//	rj := timep(feed[j].Post.Record.Val.(*bsky.FeedPost).CreatedAt)
	//	return ri.Before(rj)
	//})
	//if int64(len(feed)) > n {
	//	feed = feed[len(feed)-int(n):]
	//}
	if cCtx.Bool("json") {
		for _, p := range feed {
			json.NewEncoder(os.Stdout).Encode(p)
		}
		//} else {
		//	for _, p := range feed {
		//		//if p.Reason != nil {
		//		//continue
		//		//}
		//		printPost(p.Post)
		//	}
	}

	return nil
}
