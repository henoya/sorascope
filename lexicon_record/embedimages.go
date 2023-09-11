package lexiconRecord
import (
	"github.com/bluesky-social/indigo/lex/util"
)
type EmbedImages struct {
	LexiconTypeID string               
	Images        []*EmbedImages_Image 
}
type EmbedImages_Image struct {
	Alt   string        
	Image *util.LexBlob 
}
type EmbedImages_View struct {
	LexiconTypeID string                   
	Images        []*EmbedImages_ViewImage 
}
type EmbedImages_ViewImage struct {
	Alt      string 
	Fullsize string 
	Thumb    string 
}
