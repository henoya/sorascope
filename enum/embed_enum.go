package enum

//go:generate go run github.com/alvaroloes/enumer -type EmbedType -json
type EmbedType int // defined type
const (
	EmbedUnknown EmbedType = iota - 1
	EmbedNone
	EmbedImages
	EmbedExternal
	EmbedRecord
	EmbedRecordWithMedia
	EmbedRecordNotFound
	EmbedRecordBlocked
	EmbedFeedGenerator
	EmbedGraphListView
)
