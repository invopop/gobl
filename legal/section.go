package legal

// Section represents a part of a chapter or other section.
type Section struct {
	ID       string     `json:"id,omitempty" jsonschema:"title=ID,description=Text ID to be able to refer to this section."`
	Idx      int32      `json:"idx" jsonschema:"title=Index,description=Position inside a context."`
	Ref      string     `json:"ref,omitempty" jsonschema:"title=Ref,description=Link to the source material for this section."`
	Title    string     `json:"title,omitempty" jsonschema:"title=Title,description=Title text for this section."`
	Content  string     `json:"content" jsonschema:"title=Content,description=Actual content of the section."`
	Sections []*Section `json:"sections,omitempty" jsonschema:"title=Sub Sections,description=Additional sections to be show after this section's contents."`
}

// Calculate runs through the sections sublings and makes sure they are
// correctly indexed.
func (s *Section) Calculate() error {
	for i, ss := range s.Sections {
		if err := ss.Calculate(); err != nil {
			return err
		}
		ss.Idx = int32(i + 1) // count from 1
	}
	return nil
}
