package legal

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

func init() {
	norm.Register(
		norm.For(normalizeContract),
		norm.For(normalizeChapter),
		norm.For(normalizeSection),
	)
}

func normalizeContract(c *Contract) {
	uuid.Normalize(&c.UUID)
	c.Title = cbc.NormalizeString(c.Title)
	c.SubTitle = cbc.NormalizeString(c.SubTitle)
	c.Description = cbc.NormalizeString(c.Description)
	for i, chapter := range c.Chapters {
		chapter.Index = int32(i + 1)
	}
}

func normalizeChapter(c *Chapter) {
	c.Anchor = cbc.NormalizeString(c.Anchor)
	c.Ref = cbc.NormalizeString(c.Ref)
	c.Title = cbc.NormalizeString(c.Title)
	c.SubTitle = cbc.NormalizeString(c.SubTitle)
	for i, section := range c.Sections {
		section.Index = int32(i + 1)
	}
}

func normalizeSection(s *Section) {
	s.Anchor = cbc.NormalizeString(s.Anchor)
	s.Ref = cbc.NormalizeString(s.Ref)
	s.Comment = cbc.NormalizeString(s.Comment)
	s.Title = cbc.NormalizeString(s.Title)
	s.Content = cbc.NormalizeString(s.Content)
	for i, section := range s.Sections {
		section.Index = int32(i + 1)
	}
}
