package legal

// Chapter represents a chapter inside a document.
type Chapter struct {
	// Unique anchor for this chapter inside the document
	Anchor string `json:"$anchor,omitempty" jsonschema:"title=Anchor"`
	// Index of this chapter inside its context.
	Index int32 `json:"idx" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Link to the source material for this chapter
	Ref string `json:"$ref,omitempty" jsonschema:"title=Ref"`
	// Chapter title
	Title string `json:"title" jsonschema:"title=Title"`
	// Additional sub-title for the chapter.
	SubTitle string `json:"sub_title,omitempty" jsonschema:"title=Sub Title"`
	// Sections containing the chapter's contents.
	Sections []*Section `json:"sections,omitempty" jsonschema:"title=Sections"`
}
