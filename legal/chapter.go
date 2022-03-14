package legal

// Chapter represents a chapter inside a document.
type Chapter struct {
	// Unique anchor for this chapter inside the document
	Anchor string `json:"$anchor,omitempty" jsonschema:"title=Anchor"`
	// Index of this chapter inside it's context
	Index int32 `json:"idx" jsonschema:"title=Index"`
	// Link to the source material for this chapter
	Ref string `json:"$ref,omitempty" jsonschema:"title=Ref"`
	// Chapter title
	Title string `json:"title" jsonschema:"title=Title"`
	// Additional sub-title for the chapter.
	SubTitle string `json:"sub_title,omitempty" jsonschema:"title=Sub Title"`
	// Sections of contents of the chatper.
	Sections []*Section `json:"sections" jsonschema:"title=Sections"`
}

// Calculate runs through the chapters sections and makes sure they are
// correctly indexed.
func (c *Chapter) Calculate() error {
	for i, ss := range c.Sections {
		if err := ss.Calculate(); err != nil {
			return err
		}
		ss.Index = int32(i + 1) // count from 1
	}
	return nil
}
