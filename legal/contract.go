package legal

// Contract represents a basic legal document between a set of parties.
type Contract struct {
	// Title of the document
	Title string `json:"title" jsonschema:"title=Title"`
	// Sub-title
	SubTitle string `json:"sub_title,omitempty" jsonschema:"title=Sub Title"`
	// Brief summary reflecting on the contents.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Set of chapters that make up the contract.
	Chapters []*Chapter `json:"chapters,omitempty" jsonschema:"title=Chapters"`
}

// Calculate runs through the chapters and makes sure they are
// correctly indexed.
func (c *Contract) Calculate() error {
	for i, sc := range c.Chapters {
		if err := sc.Calculate(); err != nil {
			return err
		}
		sc.Index = int32(i + 1) // count from 1
	}
	return nil
}
