package lexiconRecord
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/bluesky-social/indigo/lex/util"
	cbg "github.com/whyrusleeping/cbor-gen"
)
type EmbedRecordWithMedia struct {
	LexiconTypeID string                      
	Media         *EmbedRecordWithMedia_Media 
	Record        *EmbedRecord                
}
type EmbedRecordWithMedia_Media struct {
	EmbedImages   *EmbedImages
	EmbedExternal *EmbedExternal
}
type EmbedRecordWithMedia_View struct {
	LexiconTypeID string                           
	Media         *EmbedRecordWithMedia_View_Media 
	Record        *EmbedRecord_View                
}
type EmbedRecordWithMedia_View_Media struct {
	EmbedImages_View   *EmbedImages_View
	EmbedExternal_View *EmbedExternal_View
}
