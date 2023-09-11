package lexiconRecord

import (
	comatprototypes "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
)

type EmbedRecord struct {
	LexiconTypeID string
	Record        *comatprototypes.RepoStrongRef
}
type EmbedRecord_View struct {
	LexiconTypeID string
	Record        *EmbedRecord_View_Record
}
type EmbedRecord_ViewBlocked struct {
	LexiconTypeID string
	Author        *bsky.FeedDefs_BlockedAuthor
	Blocked       bool
	Uri           string
}
type EmbedRecord_ViewNotFound struct {
	LexiconTypeID string
	NotFound      bool
	Uri           string
}
type EmbedRecord_ViewRecord struct {
	LexiconTypeID string
	Author        *bsky.ActorDefs_ProfileViewBasic
	Cid           string
	Embeds        []*EmbedRecord_ViewRecord_Embeds_Elem
	IndexedAt     string
	Labels        []*comatprototypes.LabelDefs_Label
	Uri           string
	Value         *util.LexiconTypeDecoder
}
type EmbedRecord_ViewRecord_Embeds_Elem struct {
	EmbedImages_View          *EmbedImages_View
	EmbedExternal_View        *EmbedExternal_View
	EmbedRecord_View          *EmbedRecord_View
	EmbedRecordWithMedia_View *EmbedRecordWithMedia_View
}
type EmbedRecord_View_Record struct {
	EmbedRecord_ViewRecord   *EmbedRecord_ViewRecord
	EmbedRecord_ViewNotFound *EmbedRecord_ViewNotFound
	EmbedRecord_ViewBlocked  *EmbedRecord_ViewBlocked
	FeedDefs_GeneratorView   *bsky.FeedDefs_GeneratorView
	GraphDefs_ListView       *bsky.GraphDefs_ListView
}
