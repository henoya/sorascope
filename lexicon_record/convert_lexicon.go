//go:generate -command lexicon $PWD/../gen-script/gen-test.sh
package lexiconRecord

//go:generate lexicon api/bsky/embedexternal.go
//go:generate lexicon api/bsky/embedimages.go
//go:generate lexicon api/bsky/embedrecord.go
//go:generate lexicon api/bsky/embedrecordWithMedia.go
