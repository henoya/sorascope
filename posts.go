package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	comatapi "github.com/bluesky-social/indigo/api"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util/cliutil"
	"github.com/henoya/sorascope/enum"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"regexp"
	"time"
)

var atDidRegexp = regexp.MustCompile(`at://(did:plc:[0-9a-z]+)/app\.bsky\.feed\.post/([0-9a-z]+)`)

func doGetPosts(cCtx *cli.Context) (err error) {
	if cCtx.Args().Present() {
		return cli.ShowSubcommandHelp(cCtx)
	}

	// DBファイルのオープン
	var db *gorm.DB
	db, err = openDB()
	if err != nil {
		return fmt.Errorf("failed to connect database")
	}

	err = migrateDB(db)
	if err != nil {
		return fmt.Errorf("failed to migrate database")
	}

	xrpcc, err := makeXRPCC(cCtx)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}
	ctx := context.Background()

	var feed []*bsky.FeedDefs_FeedViewPost

	n := cCtx.Int64("n")
	handle := cCtx.String("handle")
	handleDid := Did("")
	if handle == "" || handle == "self" {
		handleDid = Did(xrpcc.Auth.Did)
		s := cliutil.GetDidResolver(cCtx)
		phr := &comatapi.ProdHandleResolver{}
		handle, _, err = comatapi.ResolveDidToHandle(ctx, xrpcc, s, phr, string(handleDid))
		if err != nil {
			return fmt.Errorf("failed to resolve handle: %w", err)
		}
	} else {
		resolvHandle, err := comatproto.IdentityResolveHandle(ctx, xrpcc, handle)
		if err != nil {
			return fmt.Errorf("failed to resolve handle: %w", err)
		}
		handleDid = Did(resolvHandle.Did)
	}

	doAllPages := true

	ownerId := OwnerId(handleDid)
	recordCount := int64(0)

	var cursor string
	for {
		resp, err := bsky.FeedGetAuthorFeed(context.TODO(), xrpcc, string(ownerId), cursor, "posts_with_replies", n)
		if err != nil {
			return fmt.Errorf("cannot get author feed: %w", err)
		}

		feed = append(feed, resp.Feed...)

		for _, p := range resp.Feed {
			//if p.Post.Record.Val.(*bsky.FeedPost).Embed != (*bsky.FeedPost_Embed)(nil) {
			//	if p.Post.Record.Val.(*bsky.FeedPost).Embed.EmbedImages == (*bsky.EmbedImages)(nil) &&
			//		p.Post.Record.Val.(*bsky.FeedPost).Embed.EmbedExternal == (*bsky.EmbedExternal)(nil) &&
			//		p.Post.Record.Val.(*bsky.FeedPost).Embed.EmbedRecord == (*bsky.EmbedRecord)(nil) &&
			//		p.Post.Record.Val.(*bsky.FeedPost).Embed.EmbedRecordWithMedia == (*bsky.EmbedRecordWithMedia)(nil) {
			//		fmt.Printf("Post.Embed all nil embed: %s\n", p.Post.Cid)
			//		j, err := p.Post.Record.Val.(*bsky.FeedPost).Embed.MarshalJSON()
			//		if err != nil {
			//			return fmt.Errorf("cannot marshal json: %w", err)
			//		}
			//		fmt.Printf("%s\n", j)
			//	}
			//}
			//if p.Post.Cid == "bafyreihvvjvqjpighfkwvnsbag6jgbh2ijmpzq2tn65m7rjevwz5vjr4ai" {
			//	fmt.Printf("p: %v\n", p)
			//	fmt.Printf("p.Post: %v\n", p.Post)
			//	fmt.Printf("p.Post.Embed: %v\n", p.Post.Embed)
			//	fmt.Println()
			//}
			rawJson, err := json.Marshal(p)
			if err != nil {
				rawJson = []byte("{}")
				err = nil
			}
			_ = rawJson
			fmt.Println(string(rawJson))

			// PostRecordに記録
			recordCount = recordCount + 1
			//fmt.Printf("recordCount: %d\n", recordCount)
			postRecord, err := UpdateOrInsertPost(db, ownerId, p)
			if err != nil {
				return fmt.Errorf("failed to update or insert post: %w", err)
			}
			postId := postRecord.Id
			postUri := postRecord.Uri
			fmt.Printf("postId: %s, postUri: %s\n", postId, postUri)
		}

		if resp.Cursor != nil {
			if doAllPages {
				cursor = *resp.Cursor
			} else {
				break
			}
		} else {
			cursor = ""
		}

		if cursor == "" {
			break
		}
	}
	fmt.Printf("recordCount: %d\n", recordCount)

	//if len(feed) == 0 {
	//	return nil
	//}
	//for _, p := range feed {
	//	uri := p.Uri
	//	r := regexp.MustCompile(`at://(did:plc:[0-9a-z]+)/app\.bsky\.feed\.post/([0-9a-z]+)`)
	//	m := r.FindAllStringSubmatch(uri, -1)
	//	ownerDid := ""
	//	fmt.Println(m)
	//	if len(m) > 0 {
	//		ownerDid = m[0][1]
	//	}
	//	rawJson, err := json.Marshal(p)
	//	if err != nil {
	//		return fmt.Errorf("failed to marshal post: %w", err)
	//	}
	//	fmt.Println(string(rawJson))
	//	var count int64
	//	db.Model(&PostRecord{}).Where("cid = ?", p.Cid).Count(&count)
	//	if count == 0 {
	//		rec := p.Value.Val.(*bsky.FeedPost)
	//		post := &PostRecord{
	//			Uri:   AtUri(p.Uri),
	//			Cid:   AtCid(p.Cid),
	//			Owner: OwnerId(ownerDid),
	//
	//			Text: rec.Text,
	//			Json: string(rawJson),
	//		}
	//		db.Create(&post)
	//	} else {
	//		row := db.Model(&PostRecord{}).Where("cid = ?", p.Cid).Row()
	//		var post PostRecord
	//		err = row.Scan(&post)
	//		if err != nil {
	//			return fmt.Errorf("failed to scan row: %w", err)
	//		}
	//		post.Json = string(rawJson)
	//		rec := p.Value.Val.(*bsky.FeedPost)
	//		post.Text = rec.Text
	//		post.Owner = OwnerId(ownerDid)
	//	}
	//}
	//
	////sort.Slice(feed, func(i, j int) bool {
	////	ri := timep(feed[i].Value.Post.Record.Val.(*bsky.FeedPost).CreatedAt)
	////	rj := timep(feed[j].Post.Record.Val.(*bsky.FeedPost).CreatedAt)
	////	return ri.Before(rj)
	////})
	////if int64(len(feed)) > n {
	////	feed = feed[len(feed)-int(n):]
	////}
	//if cCtx.Bool("json") {
	//	for _, p := range feed {
	//		json.NewEncoder(os.Stdout).Encode(p)
	//	}
	//	//} else {
	//	//	for _, p := range feed {
	//	//		//if p.Reason != nil {
	//	//		//continue
	//	//		//}
	//	//		printPost(p.Post)
	//	//	}
	//}

	return nil
}

func extractDidFromAtUri(atUri AtUri) (did Did, err error) {
	m := atDidRegexp.FindAllStringSubmatch(string(atUri), -1)

	// DID取得
	did = Did("")
	if len(m) == 0 {
		return "", fmt.Errorf("invalid uri: %s", atUri)
	}
	did = Did(m[0][1])
	return did, nil
}

func extractTidFromAtUri(atUri AtUri) (tid Tid, err error) {
	m := atDidRegexp.FindAllStringSubmatch(string(atUri), -1)

	// DID取得
	tid = Tid("")
	if len(m) == 0 {
		return "", fmt.Errorf("invalid uri: %s", atUri)
	}
	tid = Tid(m[0][2])
	return tid, nil
}

type embedBlock struct {
	EmbedType enum.EmbedType
	EmbedDid  Did
	EmbedCid  Cid
	Uri       AtUri
	AuthorDid Did
	Blocked   bool
	Name      string
}

func inspectEmbedPost(e *bsky.FeedDefs_PostView_Embed, did Did, cid Cid) (eb *embedBlock, err error) {
	eb = &embedBlock{
		EmbedType: enum.EmbedNone,
		EmbedDid:  Did(""),
		EmbedCid:  Cid(""),
		Uri:       AtUri(""),
		AuthorDid: Did(""),
		Blocked:   false,
		Name:      "",
	}
	if e == (*bsky.FeedDefs_PostView_Embed)(nil) {
		eb.EmbedType = enum.EmbedNone
	} else {
		switch {
		case e.EmbedRecord_View != (*bsky.EmbedRecord_View)(nil):
			embedRecord := e.EmbedRecord_View.Record
			if embedRecord == (*bsky.EmbedRecord_View_Record)(nil) {
				eb.EmbedType = enum.EmbedUnknown
				fmt.Errorf("Embed type: EmbedRecord_View  invalid embed record e.EmbedRecord_View.Record is nil")
			} else {
				switch {
				case embedRecord.EmbedRecord_ViewRecord != (*bsky.EmbedRecord_ViewRecord)(nil):
					ebr := embedRecord.EmbedRecord_ViewRecord
					eb.EmbedType = enum.EmbedRecord
					eb.EmbedDid = Did(ebr.Author.Did)
					eb.EmbedCid = Cid(ebr.Cid)
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.EmbedRecord_ViewNotFound != (*bsky.EmbedRecord_ViewNotFound)(nil):
					ebr := embedRecord.EmbedRecord_ViewNotFound
					eb.EmbedType = enum.EmbedRecordNotFound
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.EmbedRecord_ViewBlocked != (*bsky.EmbedRecord_ViewBlocked)(nil):
					ebr := embedRecord.EmbedRecord_ViewBlocked
					eb.EmbedType = enum.EmbedRecordBlocked
					eb.AuthorDid = Did(ebr.Author.Did)
					eb.Blocked = ebr.Blocked
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.FeedDefs_GeneratorView != (*bsky.FeedDefs_GeneratorView)(nil):
					ebr := embedRecord.FeedDefs_GeneratorView
					eb.EmbedType = enum.EmbedFeedGenerator
					creatorDid := Did(ebr.Creator.Did)
					eb.EmbedDid = creatorDid
					eb.EmbedCid = Cid(ebr.Cid)
					eb.AuthorDid = creatorDid
					eb.Uri = AtUri(ebr.Uri)
					eb.Name = ebr.DisplayName
				case embedRecord.GraphDefs_ListView != (*bsky.GraphDefs_ListView)(nil):
					ebr := embedRecord.GraphDefs_ListView
					eb.EmbedType = enum.EmbedGraphListView
					creatorDid := Did(ebr.Creator.Did)
					eb.EmbedCid = Cid(ebr.Cid)
					eb.AuthorDid = creatorDid
					eb.Uri = AtUri(ebr.Uri)
					eb.Name = ebr.Name
				default:
					eb.EmbedType = enum.EmbedUnknown
					fmt.Errorf("Embed type: EmbedRecord_View  invalid embed record child of e.EmbedRecord_View.Record is all nil")
				}
			}
		case e.EmbedRecordWithMedia_View != (*bsky.EmbedRecordWithMedia_View)(nil):
			embedRecord := e.EmbedRecordWithMedia_View.Record
			if embedRecord == (*bsky.EmbedRecord_View)(nil) {
				eb.EmbedType = enum.EmbedUnknown
				fmt.Errorf("Embed type: EmbedRecordWithMedia_View  invalid embed record e.EmbedRecord_View.Record is nil")
			} else {
				switch {
				case embedRecord.Record.EmbedRecord_ViewRecord != (*bsky.EmbedRecord_ViewRecord)(nil):
					ebr := embedRecord.Record.EmbedRecord_ViewRecord
					eb.EmbedType = enum.EmbedRecordWithMedia
					eb.EmbedDid = Did(ebr.Author.Did)
					eb.EmbedCid = Cid(ebr.Cid)
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.Record.EmbedRecord_ViewNotFound != (*bsky.EmbedRecord_ViewNotFound)(nil):
					ebr := embedRecord.Record.EmbedRecord_ViewNotFound
					eb.EmbedType = enum.EmbedRecordNotFound
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.Record.EmbedRecord_ViewBlocked != (*bsky.EmbedRecord_ViewBlocked)(nil):
					ebr := embedRecord.Record.EmbedRecord_ViewBlocked
					eb.EmbedType = enum.EmbedRecordBlocked
					eb.AuthorDid = Did(ebr.Author.Did)
					eb.Blocked = ebr.Blocked
					eb.Uri = AtUri(ebr.Uri)
				case embedRecord.Record.FeedDefs_GeneratorView != (*bsky.FeedDefs_GeneratorView)(nil):
					ebr := embedRecord.Record.FeedDefs_GeneratorView
					eb.EmbedType = enum.EmbedFeedGenerator
					creatorDid := Did(ebr.Creator.Did)
					eb.EmbedDid = creatorDid
					eb.EmbedCid = Cid(ebr.Cid)
					eb.AuthorDid = creatorDid
					eb.Uri = AtUri(ebr.Uri)
					eb.Name = ebr.DisplayName
				case embedRecord.Record.GraphDefs_ListView != (*bsky.GraphDefs_ListView)(nil):
					ebr := embedRecord.Record.GraphDefs_ListView
					eb.EmbedType = enum.EmbedGraphListView
					creatorDid := Did(ebr.Creator.Did)
					eb.EmbedCid = Cid(ebr.Cid)
					eb.AuthorDid = creatorDid
					eb.Uri = AtUri(ebr.Uri)
					eb.Name = ebr.Name
				default:
					eb.EmbedType = enum.EmbedUnknown
					fmt.Errorf("Embed type: EmbedRecordWithMedia_View  invalid embed record child of e.EmbedRecord_View.Record is all nil")
				}
			}
		case e.EmbedImages_View != (*bsky.EmbedImages_View)(nil):
			eb.EmbedType = enum.EmbedImages
			eb.EmbedDid = did
			eb.EmbedCid = cid
		case e.EmbedExternal_View != (*bsky.EmbedExternal_View)(nil):
			eb.EmbedType = enum.EmbedExternal
			eb.EmbedDid = did
			eb.EmbedCid = cid
		default:
			eb.EmbedType = enum.EmbedUnknown
			eb.EmbedDid = Did("")
			eb.EmbedCid = Cid("")
			fmt.Errorf("invalid embed record child of e is all nil")
		}
	}
	return eb, nil
}

func setupPostRecord(p *bsky.FeedDefs_FeedViewPost, postCid Cid, postDid Did) (postRecord *PostRecord, err error) {
	postRecord = &PostRecord{}
	postView := p.Post
	record := postView.Record.Val.(*bsky.FeedPost)
	atUri := AtUri(postView.Uri)
	postRecord.Cid = postCid
	postRecord.Did = postDid
	postRecord.Uri = atUri
	postRecord.Tid, err = extractTidFromAtUri(atUri)
	if err != nil {
		return nil, err
	}
	postRecord.Text = record.Text
	eb, err := inspectEmbedPost(postView.Embed, postDid, postCid)
	postRecord.EmbedType = eb.EmbedType
	postRecord.EmbedDid = eb.EmbedDid
	postRecord.EmbedCid = eb.EmbedCid
	postRecord.EmbedUri = eb.Uri
	postRecord.EmbedAuthorDid = eb.AuthorDid
	postRecord.EmbedName = eb.Name
	postRecord.CreatedAt, err = timepWithError(record.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CreatedAt: %w", err)
	}
	postRecord.IndexedAt, err = timepWithError(postView.IndexedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IndexedAt: %w", err)
	}
	postRecord.PostType = enum.PostTypePost
	if p.Reply != nil {
		postRecord.PostType = enum.PostTypeReply
	}
	postRecord.Langs = postView.Record.Val.(*bsky.FeedPost).Langs
	return postRecord, nil
}

func setupPostStatus(postView *bsky.FeedDefs_PostView, postRecord *PostRecord) (postStatus *PostStatus, err error) {
	postStatus = &PostStatus{}
	postStatus.Cid = postRecord.Cid
	postStatus.Did = postRecord.Did
	postStatus.Uri = postRecord.Uri
	postStatus.CreatedAt = postRecord.CreatedAt
	updatedAt := time.Now().Local()
	postStatus.UpdatedAt = &updatedAt

	if postView.LikeCount != nil {
		postStatus.LikeCount = *postView.LikeCount
	} else {
		postStatus.LikeCount = 0
	}
	if postView.ReplyCount != nil {
		postStatus.ReplyCount = *postView.ReplyCount
	} else {
		postStatus.ReplyCount = 0
	}
	if postView.RepostCount != nil {
		postStatus.RepostCount = *postView.RepostCount
	} else {
		postStatus.RepostCount = 0
	}
	postStatus.Labels = []string{}
	if postView.Record.Val.(*bsky.FeedPost).Labels != nil {
		labels := postView.Record.Val.(*bsky.FeedPost).Labels.LabelDefs_SelfLabels.Values
		if len(labels) > 0 {
			postStatus.Labels = make([]string, len(labels))
			for i, l := range labels {
				postStatus.Labels[i] = l.Val
			}
		}
	}

	return postStatus, nil
}

func setupPostHistory(p *bsky.FeedDefs_FeedViewPost, postRecord *PostRecord, owner OwnerId) (postHistory *PostHistory, err error) {
	postHistory = &PostHistory{}
	postFeedType := enum.PostFeedTypePostView
	indexedAt := postRecord.CreatedAt
	if p.Reason != nil {
		if p.Reason.FeedDefs_ReasonRepost != nil {
			postFeedType = enum.PostFeedTypeReasonRepost
			indexedAt, err = timepWithError(p.Reason.FeedDefs_ReasonRepost.IndexedAt)
		} else {
			postFeedType = enum.PostFeedTypeUnknown
		}
	}
	postHistory.Owner = owner
	postHistory.Cid = postRecord.Cid
	postHistory.Did = postRecord.Did
	postHistory.Uri = postRecord.Uri
	postHistory.Tid = postRecord.Tid
	postHistory.PostFeedType = postFeedType
	postHistory.CreatedAt = postRecord.CreatedAt
	postHistory.IndexedAt = indexedAt
	postHistory.Text = postRecord.Text
	return postHistory, nil
}

func setupPostHistoryStatus(p *bsky.FeedDefs_FeedViewPost, postRecord *PostRecord, owner OwnerId) (postHistoryStatus *PostHistoryStatus, err error) {
	postHistoryStatus = &PostHistoryStatus{}
	postHistoryStatus.Owner = owner
	postHistoryStatus.Cid = postRecord.Cid
	postHistoryStatus.Uri = postRecord.Uri

	postHistoryStatus.BlockedBy = false
	if p.Post.Author.Viewer.BlockedBy != nil {
		postHistoryStatus.BlockedBy = *p.Post.Author.Viewer.BlockedBy
	}
	postHistoryStatus.Muted = false
	if p.Post.Author.Viewer.Muted != nil {
		postHistoryStatus.Muted = *p.Post.Author.Viewer.Muted
	}
	return postHistoryStatus, nil
}

func setupAuthorRecord(p *bsky.FeedDefs_FeedViewPost) (authorRecord *AuthorRecord, err error) {
	authorRecord = &AuthorRecord{}
	authorRecord.DisplayName = ""
	authorRecord.AvatarUrl = ""
	authorRecord.Description = ""
	authorRecord.ChangeDateTime = nil
	authorRecord.Json = ""
	authorRecord.Hash = ""

	if p.Post.Author == (*bsky.ActorDefs_ProfileViewBasic)(nil) {
		fmt.Errorf("author record is nil")
	}
	author := p.Post.Author
	authorRecord.Did = Did(author.Did)
	authorRecord.Revision = 0
	if author.DisplayName != nil {
		authorRecord.DisplayName = *author.DisplayName
	}
	authorRecord.Handle = Handle(author.Handle)
	if author.Avatar != nil {
		authorRecord.AvatarUrl = *author.Avatar
	}
	return authorRecord, nil
}

func UpdateOrInsertPost(db *gorm.DB, owner OwnerId, p *bsky.FeedDefs_FeedViewPost) (postRecord *PostRecord, err error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if p == nil {
		return nil, fmt.Errorf("post is nil")
	}
	postView := p.Post
	//record := postView.Record.Val.(*bsky.FeedPost)

	rawJson, err := json.Marshal(p)

	//postRecord = &PostRecord{}
	//postStatus := &PostStatus{}
	//postHistory := &PostHistory{}
	//postHistoryStatus := &PostHistoryStatus{}

	postCid := Cid(postView.Cid)
	postDid, err := extractDidFromAtUri(AtUri(postView.Uri))
	if err != nil {
		return nil, fmt.Errorf("cannot extract did from uri: %w", err)
	}

	// PostRecord の設定
	postRecord, err = setupPostRecord(p, postCid, postDid)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect post record: %w", err)
	}

	// PostStatus の設定
	postStatus, err := setupPostStatus(postView, postRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect post status: %w", err)
	}
	postStatus.Json = string(rawJson)

	// PostHistory の設定
	postHistory, err := setupPostHistory(p, postRecord, owner)

	// postHistoryStatus の設定
	postHistoryStatus, err := setupPostHistoryStatus(p, postRecord, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect post history status: %w", err)
	}

	// AuthorRecord の設定
	authorRecord, err := setupAuthorRecord(p)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect author record: %w", err)
	}
	_ = authorRecord

	// PostRecordに記録
	err = db.Transaction(func(tx *gorm.DB) error {
		// PostRecordが存在するかチェック
		var count int64
		idHash, err := calcPostRecordHash(postRecord)
		postId := PostRecordId(idHash)
		if err := tx.Model(&PostRecord{}).Where("id = ? OR (cid = ? AND did = ?)", postId, postRecord.Cid, postRecord.Did).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			// PostRecordが存在しない場合、新規にレコードを作成する
			fmt.Printf("create new post record: %s\n", postRecord.Cid)
			postRecord.Id = postId
			if err := tx.Create(&postRecord).Error; err != nil {
				return err
			}
			fmt.Printf("Create post record id: %d\n", postId)
		} else {
			fmt.Printf("update post record: %s\n", count)
			rows, err := tx.Model(&PostRecord{}).Where("id = ? OR (cid = ? AND did = ?)", postId, postRecord.Cid, postRecord.Did).Rows()
			if err != nil {
				return err
			}
			for rows.Next() {
				fmt.Printf("Select post record cid: %s  did:%s\n", postRecord.Cid, postRecord.Did)
				var post PostRecord
				err := db.ScanRows(rows, &post)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				postId = post.Id
				fmt.Printf("Get post record id: %d\n", postId)
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
		}
		if len(postId) == 0 {
			return fmt.Errorf("postId is 0")
		}
		// PostStatusが存在するかチェック
		if err := tx.Model(&PostStatus{}).Where("id = ?", postId).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			// PostStatusが存在しない場合、新規にレコードを作成する
			postStatus.Id = postId
			if err := tx.Create(&postStatus).Error; err != nil {
				return err
			}
		} else {
			rows, err := tx.Model(&PostStatus{}).Where("id = ?", postId).Rows()
			if err != nil {
				return err
			}
			var postSt PostStatus
			labelDiff := false
			for rows.Next() {
				err := db.ScanRows(rows, &postSt)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				if postSt.Cid != postStatus.Cid || postSt.Did != postStatus.Did || postSt.Uri != postStatus.Uri {
					return fmt.Errorf("postId: %d postStatus.Cid(%s) != postSt.Cid(%s) || postStatus.Did(%s) != postSt.Did(%s) || postStatus.Uri(%s) != postSt.Uri(%s)", postId, postStatus.Cid, postSt.Cid, postStatus.Did, postSt.Did, postStatus.Uri, postSt.Uri)
				}
				if len(postSt.Labels) != len(postStatus.Labels) {
					labelDiff = true
				} else {
					for i, labelItem := range postStatus.Labels {
						if labelItem != postSt.Labels[i] {
							labelDiff = true
							break
						}
					}
				}
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
			if postSt.LikeCount != postStatus.LikeCount || postSt.ReplyCount != postStatus.ReplyCount || postSt.RepostCount != postSt.RepostCount || labelDiff {
				postSt.LikeCount = postStatus.LikeCount
				postSt.ReplyCount = postStatus.ReplyCount
				postSt.RepostCount = postSt.RepostCount
				postSt.Labels = postStatus.Labels
				t := time.Now().Local()
				postSt.UpdatedAt = &t
				if err := tx.Save(&postSt).Error; err != nil {
					return err
				}
			}
		}
		// PostHistoryが存在するかチェック
		postHistoryId, err := calcPostHistoryHash(postHistory)
		if err != nil {
			return err
		}
		if len(postHistoryId) == 0 {
			return fmt.Errorf("postHistoryId is 0")
		}
		postHistory.Id = PostHistroyId(postHistoryId)
		if err := tx.Model(&PostHistory{}).Where("id = ?", postHistoryId).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			// PostStatusが存在しない場合、新規にレコードを作成する
			if err := tx.Create(&postHistory).Error; err != nil {
				return err
			}
		} else {
			rows, err := tx.Model(&PostHistory{}).Where("id = ?", postHistoryId).Rows()
			if err != nil {
				return err
			}
			var postHi PostHistory
			for rows.Next() {
				err := db.ScanRows(rows, &postHi)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
		}
		// PostHistoryStatusが存在するかチェック
		if err := tx.Model(&PostHistoryStatus{}).Where("id = ?", postHistoryId).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			// PostHistoryStatusが存在しない場合、新規にレコードを作成する
			postHistoryStatus.Id = PostHistroyId(postHistoryId)
			if err := tx.Create(&postHistoryStatus).Error; err != nil {
				return err
			}
		} else {
			rows, err := tx.Model(&PostHistoryStatus{}).Where("id = ?", postHistoryId).Rows()
			if err != nil {
				return err
			}
			var postSt PostHistoryStatus
			for rows.Next() {
				err := db.ScanRows(rows, &postSt)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
			if postSt.BlockedBy != postHistoryStatus.BlockedBy || postSt.Muted != postHistoryStatus.Muted {
				postSt.BlockedBy = postHistoryStatus.BlockedBy
				postSt.Muted = postHistoryStatus.Muted
				if err := tx.Save(&postSt).Error; err != nil {
					return err
				}
			}
		}

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return nil, err
	}
	return postRecord, nil
}
func calcHash(str string) (hash string, err error) {
	idHash := sha256.Sum256([]byte(str))
	hash = hex.EncodeToString(idHash[:])
	return hash, nil
}

func calcPostRecordHash(postRecord *PostRecord) (hash string, err error) {
	createdAtUTC := postRecord.CreatedAt.UTC().Format(time.RFC3339)
	idString := string(postRecord.Did) + string(postRecord.Cid) + createdAtUTC
	id, err := calcHash(idString)
	if err != nil {
		return "", err
	}
	return id, nil
}

func calcPostHistoryHash(postHistory *PostHistory) (hash string, err error) {
	indexdAtUTC := postHistory.IndexedAt.UTC().Format(time.RFC3339)
	idString := string(postHistory.Owner) + string(postHistory.Did) + string(postHistory.Cid) + indexdAtUTC
	id, err := calcHash(idString)
	if err != nil {
		return "", err
	}
	return id, nil
}

func getPost(cCtx *cli.Context, uri string) (postData []*bsky.FeedDefs_PostView, err error) {
	xrpcc, err := makeXRPCC(cCtx)
	if err != nil {
		return nil, fmt.Errorf("cannot create client: %w", err)
	}

	//r := regexp.MustCompile(`at://(did:plc:[0-9a-z]+)/app\.bsky\.feed\.post/([0-9a-z]+)`)
	//m := r.FindAllStringSubmatch(uri, -1)
	//ownerDid := ""
	//fmt.Println(m)
	//if len(m) > 0 {
	//	ownerDid = m[0][1]
	//}

	//ids := strings.Split(uri, "/app.bsky.feed.post/")
	//if len(ids) != 2 {
	//	return nil, fmt.Errorf("uri does not contain 2 parts: %s", uri)
	//}
	//did := ids[0]
	//postID := ids[1]

	//ctx := context.Background()
	//idResolve, err := comatproto.IdentityResolveHandle(ctx, xrpcc, handle)
	//if err != nil {
	//	return nil, err
	//}
	//out := idResolve.Did
	//did := "at://" + out + "/app.bsky.feed.post/" + postID
	dids := []string{uri}
	resp, err := bsky.FeedGetPosts(context.TODO(), xrpcc, dids)
	if err != nil {
		return nil, fmt.Errorf("cannot get post thread: %w", err)
	}
	if len(resp.Posts) == 0 {
		return nil, fmt.Errorf("no posts found")
	}
	return resp.Posts, nil
}
