package post

import (
	"time"

	"github.com/henoya/sorascope/enum"
	"github.com/henoya/sorascope/typedef"

	_ "github.com/mattn/go-sqlite3"
)

type PostRecords []*PostRecord

type Image struct {
	Id       int64       `json:"id" gorm:"type:integer;primary_key;auto_increment"`
	Did      typedef.Did `json:"did" gorm:"type:text;index:idx_image_did"`
	Cid      typedef.Cid `json:"cid" gorm:"type:text;index:idx_image_cid"`
	MimeType string      `json:"mime_type" gorm:"type:text"`
	Size     int64       `json:"size" gorm:"type:integer"`
	URL      string      `json:"url" gorm:"type:text;index:idx_image_url"`
	ThumbURL string      `json:"thumb_url" gorm:"type:text"`
	Alt      string      `json:"alt" gorm:"type:text"`
}

type EmbedImages struct {
	Id       int64       `json:"id" gorm:"type:integer;primary_key"`
	Did      typedef.Did `json:"did" gorm:"type:text;index:idx_embed_image_did"`
	PostCid  typedef.Cid `json:"cid" gorm:"type:text;index:idx_embed_image_cid"`
	Index    int         `json:"index" gorm:"type:integer"`
	ImageCid typedef.Cid `json:"image_cid" gorm:"type:text"`
}

type EmbedExternal struct {
	Id          int64       `json:"id" gorm:"type:integer;primary_key"`
	Did         typedef.Did `json:"did" gorm:"type:text;index:idx_embed_external_did"`
	PostCid     typedef.Cid `json:"cid" gorm:"type:text;index:idx_embed_external_cid"`
	ThumbCid    typedef.Cid `json:"thumb_cid" gorm:"type:text;index:idx_embed_external_thumb_cid"`
	Description string      `json:"description" gorm:"type:text"`
	Title       string      `json:"title" gorm:"type:text"`
	Uri         string      `json:"uri" gorm:"type:text;index:idx_embed_external_uri"`
}

type AuthorRecord struct {
	Id             int64          `json:"id" gorm:"type:integer;primary_key"`
	Did            typedef.Did    `json:"did" gorm:"type:text;index:idx_author_record_did,unique"`
	Revision       int            `json:"revision" gorm:"type:integer;index:idx_author_record_did,unique"`
	DisplayName    string         `json:"display_name" gorm:"type:text"`
	Handle         typedef.Handle `json:"handle" gorm:"type:text"`
	AvatarUrl      string         `json:"avatar_url" gorm:"type:text"`
	Description    string         `json:"description" gorm:"type:text"`
	ChangeDateTime *time.Time     `json:"change_date_time" gorm:"type:datetime;nullable"`
	Json           string         `json:"json" gorm:"type:text"`
	Hash           string         `json:"hash" gorm:"type:text"`
}

type PostStatus struct {
	Id          typedef.PostRecordId `json:"id" gorm:"type:text;primary_key"`
	Cid         typedef.Cid          `json:"cid" gorm:"type:text;index:idx_post_status_cid"`
	Did         typedef.Did          `json:"did" gorm:"type:text;index:idx_post_status_did"`
	Uri         typedef.AtUri        `json:"uri" gorm:"type:text;index:idx_post_status_uri"`
	LikeCount   int64                `json:"like_count" gorm:"type:integer"`
	ReplyCount  int64                `json:"reply_count" gorm:"type:integer"`
	RepostCount int64                `json:"repost_count" gorm:"type:integer"`
	Labels      typedef.StringArray  `json:"labels" gorm:"type:text;serializer:json"`
	Json        string               `json:"json" gorm:"type:text"`
	CreatedAt   *time.Time           `json:"created_at" gorm:"type:datetime;index:idx_post_status_created_at"`
	UpdatedAt   *time.Time           `json:"updated_at" gorm:"type:datetime;index:idx_post_status_updated_at"`
	DeletedAr   *time.Time           `json:"deleted_ar" gorm:"type:datetime;nullable:index:idx_post_status_deleted_at"`
}

type PostRecord struct {
	Id             typedef.PostRecordId `json:"id" gorm:"type:text;primary_key"`
	Cid            typedef.Cid          `json:"cid" gorm:"type:text;index:idx_post_record_cid"`
	Did            typedef.Did          `json:"did" gorm:"type:text;index:idx_post_record_did"`
	Uri            typedef.AtUri        `json:"uri" gorm:"type:text;index:idx_post_record_uri"`
	Tid            typedef.Tid          `json:"tid" gorm:"type:text;index:idx_post_record_tid"`
	AuthorRevision int                  `json:"author_revision" gorm:"type:integer"`
	CreatedAt      *time.Time           `json:"created_at" gorm:"type:datetime;index:idx_post_record_created_at"`
	IndexedAt      *time.Time           `json:"indexed_at" gorm:"type:datetime;index:idx_post_record_indexed_at"`
	Text           string               `json:"text" gorm:"type:text"`
	PostType       enum.PostType        `json:"type" gorm:"type:integer"`
	EmbedType      enum.EmbedType       `json:"embed_type" gorm:"type:text"`
	EmbedDid       typedef.Did          `json:"embed_did" gorm:"type:text;nullable"`
	EmbedCid       typedef.Cid          `json:"embed_cid" gorm:"type:text;nullable"`
	EmbedUri       typedef.AtUri        `json:"embed_uri" gorm:"type:text"`
	EmbedAuthorDid typedef.Did          `json:"embed_author_did" gorm:"type:text"`
	EmbedBlocked   bool                 `json:"embed_blocked" gorm:"type:boolean"`
	EmbedName      string               `json:"embed_name" gorm:"type:text"`
	Langs          typedef.StringArray  `json:"langs" gorm:"type:text;serializer:json"`
	DeletedAr      *time.Time           `json:"deleted_ar" gorm:"type:datetime;nullable;index:idx_post_record_deleted_at"`
}

type PostHistoryStatus struct {
	Id        typedef.PostHistroyId `json:"id" gorm:"type:text;primary_key"`                            // id
	Owner     typedef.OwnerId       `json:"owner" gorm:"type:text;index:idx_post_history_status_owner"` // ポスト履歴の所持者(所持者のDID
	Cid       typedef.Cid           `json:"cid" gorm:"type:text;index:idx_post_history_status_cid"`     // ポストのCID
	Uri       typedef.AtUri         `json:"uri" gorm:"type:text;index:idx_post_history_status_uri"`     // ポストのuri
	BlockedBy bool                  `json:"blocked_by" gorm:"type:boolean"`                             // ブロックされている
	Muted     bool                  `json:"muted" gorm:"type:boolean"`                                  // ミュートされている
}

type PostHistories []*PostHistory

type PostHistory struct {
	Id           typedef.PostHistroyId `json:"id" gorm:"type:text;primary_key"`                     // id
	Owner        typedef.OwnerId       `json:"owner" gorm:"type:text;index:idx_post_history_owner"` // ポスト履歴の所持者(所持者のDID
	Did          typedef.Did           `json:"did" gorm:"type:text;index:idx_post_history_did"`     // ポストの投稿者のDID
	Cid          typedef.Cid           `json:"cid" gorm:"type:text;index:idx_post_history_cid"`     // ポストのCID
	Uri          typedef.AtUri         `json:"uri" gorm:"type:text;index:idx_post_history_uri"`     // ポストのuri
	Tid          typedef.Tid           `json:"tid" gorm:"type:text;index:idx_post_history_tid"`
	PostFeedType enum.PostFeedType     `json:"post_feed_type" gorm:"type:integer"`                                // ポストがfeedに出てきた理由(repost)
	CreatedAt    *time.Time            `json:"created_at" gorm:"type:datetime;index:idx_post_history_created_at"` // 作成日時
	IndexedAt    *time.Time            `json:"indexed_at" gorm:"type:datetime;index:idx_post_history_indexed_at"` // インデックス日時
	Text         string                `json:"text" gorm:"type:text"`                                             // テキスト
}
