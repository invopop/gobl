package legal

// Contract represents a basic legal document between a set of parties.
type Contract struct {
	Title    string     `json:"title"`
	SubTitle string     `json:"sub_title,omitempty"`
	Chapters []*Chapter `json:"chapters,omitempty"`
}

// Calculate runs through the chapters and makes sure they are
// correctly indexed.
func (c *Contract) Calculate() error {
	for i, sc := range c.Chapters {
		if err := sc.Calculate(); err != nil {
			return err
		}
		sc.Idx = int32(i + 1) // count from 1
	}
	return nil
}
