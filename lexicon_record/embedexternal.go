package lexiconRecord
import (
	"github.com/bluesky-social/indigo/lex/util"
)
type EmbedExternal struct {
	LexiconTypeID string                  
	External      *EmbedExternal_External 
}
type EmbedExternal_External struct {
	Description string        
	Thumb       *util.LexBlob 
	Title       string        
	Uri         string        
}
type EmbedExternal_View struct {
	LexiconTypeID string                      
	External      *EmbedExternal_ViewExternal 
}
type EmbedExternal_ViewExternal struct {
	Description string  
	Thumb       *string 
	Title       string  
	Uri         string  
}
