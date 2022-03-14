package legal

// Section represents a part of a chapter or other section.
type Section struct {
	// Unique anchor to be able to refer to this section inside a document
	Anchor string `json:"$anchor,omitempty" jsonschema:"title=Anchor"`
	// Position inside the sections context
	Index int32 `json:"idx" jsonschema:"title=Index"`
	// Link to any reference material for this section.
	Ref string `json:"$ref,omitempty" jsonschema:"title=Ref"`
	// Additional comments regarding this section not meant to be included in output.
	Comment string `json:"$comment,omitempty" jsonschema:"title=Comment"`
	// Title text.
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Text contents for the section with Markdown formatting.
	Content string `json:"content" jsonschema:"title=Content"`
	// Sub-sections show after this sections contents.
	Sections []*Section `json:"sections,omitempty" jsonschema:"title=Sub Sections"`
}

// Calculate runs through the sections sublings and makes sure they are
// correctly indexed.
func (s *Section) Calculate() error {
	for i, ss := range s.Sections {
		if err := ss.Calculate(); err != nil {
			return err
		}
		ss.Index = int32(i + 1) // count from 1
	}
	return nil
}
