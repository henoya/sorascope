package enum

//go:generate go run github.com/alvaroloes/enumer -type PostType -json
type PostType int // defined type
const (
	PostTypeUnknown PostType = iota - 1
	PostTypePost
	PostTypeReply
)

//go:generate enumer -type PostFeedType -json
type PostFeedType int // defined type
const (
	PostFeedTypeUnknown PostFeedType = iota - 1
	PostFeedTypePostView
	PostFeedTypePostViewerState
	PostFeedTypeFeedViewPost
	PostFeedTypeReplyRef
	PostFeedTypeReasonRepost
	PostFeedTypeThreadViewPost
	PostFeedTypeNotFoundPost
	PostFeedTypeBlockedPost
	PostFeedTypeBlockedAuthor
	PostFeedTypeGeneratorView
	PostFeedTypeGeneratorViewerState
	PostFeedTypeSkeletonFeedPost
	PostFeedTypeSkeletonReasonRepost
)
