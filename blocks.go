package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	//"github.com/bluesky-social/indigo/api/bsky"
	"github.com/urfave/cli/v2"
	//gorm.io/gorm"
	"os"
	//"regexp"
)

func doGetBlocks(cCtx *cli.Context) (err error) {
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
			var collection string = "app.bsky.graph.block"
			resp, err := atproto.RepoListRecords(context.TODO(), xrpcc, collection, cursor, n, handle, true, "", "")
			//func GraphGetBlocks(ctx context.Context, c *xrpc.Client, cursor string, limit int64) (*GraphGetBlocks_Output, error) {
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
			var collection string = "app.bsky.graph.block"
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

	if len(feed) == 0 {
		return fmt.Errorf("no posts found")
	}
	// for _, p := range feed {
	// 	uri := p.Uri
	// 	r := regexp.MustCompile(`at://(did:plc:[0-9a-z]+)/app\.bsky\.feed\.post/([0-9a-z]+)`)
	// 	m := r.FindAllStringSubmatch(uri, -1)
	// 	ownerDid := ""
	// 	postId := ""
	// 	fmt.Println(m)
	// 	if len(m) > 0 {
	// 		ownerDid = m[0][1]
	// 		postId = m[0][2]
	// 	}
	// 	rawJson, err := json.Marshal(p)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to marshal post: %w", err)
	// 	}
	// 	fmt.Println(string(rawJson))
	// 	var count int64
	// 	db.Model(&PostRecord{}).Where("cid = ?", p.Cid).Count(&count)
	// 	if count == 0 {
	// 		rec := p.Value.Val.(*bsky.FeedPost)
	// 		post := &PostRecord{
	// 			Cid:    p.Cid,
	// 			Uri:    p.Uri,
	// 			Owner:  ownerDid,
	// 			PostId: postId,
	// 			Text:   rec.Text,
	// 			Json:   string(rawJson),
	// 		}
	// 		db.Create(&post)
	// 	} else {
	// 		row := db.Model(&PostRecord{}).Where("cid = ?", p.Cid).Row()
	// 		var post PostRecord
	// 		err = row.Scan(&post)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to scan row: %w", err)
	// 		}
	// 		post.Json = string(rawJson)
	// 		rec := p.Value.Val.(*bsky.FeedPost)
	// 		post.Text = rec.Text
	// 		post.Owner = ownerDid
	// 		post.PostId = postId
	// 	}
	// }

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
