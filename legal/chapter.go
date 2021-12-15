package legal

// Chapter represents a chapter inside a document.
type Chapter struct {
	ID       string     `json:"id,omitempty" jsonschema:"title=ID,description=Unique ID for this chapter inside the document."`
	Idx      int32      `json:"idx" jsonschema:"title=Index,description=Index of this chapter inside the document."`
	Ref      string     `json:"ref,omitempty" jsonschema:"title=Ref,description=Link to the source material for this chapter."`
	Title    string     `json:"title" jsonschema:"title=Title,description=Chapter Title."`
	SubTitle string     `json:"sub_title,omitempty" jsonschema:"title=Sub Title,description=Additional sub-title for this chapter."`
	Sections []*Section `json:"sections" jsonschema:"title=Sections,description=Contents of this chapter."`
}

// Calculate runs through the chapters sections and makes sure they are
// correctly indexed.
func (c *Chapter) Calculate() error {
	for i, ss := range c.Sections {
		if err := ss.Calculate(); err != nil {
			return err
		}
		ss.Idx = int32(i + 1) // count from 1
	}
	return nil
}
