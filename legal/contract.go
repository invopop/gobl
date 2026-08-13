package legal

import (
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

// Contract represents the structured body of a legal document.
type Contract struct {
	uuid.Identify

	// Title of the document
	Title string `json:"title" jsonschema:"title=Title"`
	// Sub-title
	SubTitle string `json:"sub_title,omitempty" jsonschema:"title=Sub Title"`
	// Brief summary reflecting on the contents.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Set of chapters that make up the contract.
	Chapters []*Chapter `json:"chapters,omitempty" jsonschema:"title=Chapters"`
}

// Calculate prepares the contract for use as a GOBL document payload.
// Intrinsic normalization is registered by type; this method is only the
// standard document-level integration hook used by schema.Object.
func (c *Contract) Calculate() error {
	norm.Normalize(c)
	return nil
}
