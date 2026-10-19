package legal

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

func init() {
	norm.Register(
		norm.For(normalizeAnalysis),
		norm.For(normalizeAnnotation),
		norm.For(normalizeAssent),
		norm.For(normalizeContract),
		norm.For(normalizeChapter),
		norm.For(normalizeSection),
		norm.For(normalizeParty),
		norm.For(normalizeSignatory),
		norm.For(normalizeEffect),
		norm.For(normalizeRecital),
		norm.For(normalizeDefinition),
		norm.For(normalizeGoverningLaw),
		norm.For(normalizeDisputeResolution),
		norm.For(normalizeExecution),
		norm.For(normalizeReference),
		norm.For(normalizeResource),
	)
}

func normalizeContract(c *Contract) {
	uuid.Normalize(&c.UUID)
	if c.Agreement.IsZero() && !c.UUID.IsZero() {
		c.Agreement = c.UUID
	}
	uuid.Normalize(&c.Agreement)
	if c.Kind == "" {
		c.Kind = ContractKindAgreement
	}
	c.Code = cbc.NormalizeCode(c.Code)
	c.Version = cbc.NormalizeCode(c.Version)
	c.Title = cbc.NormalizeString(c.Title)
	c.SubTitle = cbc.NormalizeString(c.SubTitle)
	c.Description = cbc.NormalizeString(c.Description)
	for i, chapter := range c.Chapters {
		chapter.Index = int32(i + 1)
	}
}

func normalizeParty(p *Party) {
	p.DefinedAs = cbc.NormalizeString(p.DefinedAs)
	p.Description = cbc.NormalizeString(p.Description)
}

func normalizeSignatory(s *Signatory) {
	s.Capacity = cbc.NormalizeString(s.Capacity)
}

func normalizeEffect(e *Effect) {
	e.Trigger = cbc.NormalizeString(e.Trigger)
}

func normalizeRecital(r *Recital) {
	r.Anchor = cbc.NormalizeString(r.Anchor)
	r.Content = cbc.NormalizeString(r.Content)
}

func normalizeDefinition(d *Definition) {
	d.Anchor = cbc.NormalizeString(d.Anchor)
	d.Term = cbc.NormalizeString(d.Term)
	d.Meaning = cbc.NormalizeString(d.Meaning)
}

func normalizeGoverningLaw(l *GoverningLaw) {
	l.Jurisdiction = cbc.NormalizeString(l.Jurisdiction)
	l.Instrument = cbc.NormalizeString(l.Instrument)
}

func normalizeDisputeResolution(d *DisputeResolution) {
	d.Forum = cbc.NormalizeString(d.Forum)
	d.Seat = cbc.NormalizeString(d.Seat)
	d.Rules = cbc.NormalizeString(d.Rules)
}

func normalizeExecution(e *Execution) {
	e.Statement = cbc.NormalizeString(e.Statement)
}

func normalizeReference(r *Reference) {
	uuid.Normalize(&r.UUID)
	r.Code = cbc.NormalizeCode(r.Code)
	r.Title = cbc.NormalizeString(r.Title)
	r.URL = cbc.NormalizeString(r.URL)
}

func normalizeResource(r *Resource) {
	r.Title = cbc.NormalizeString(r.Title)
	r.URL = cbc.NormalizeString(r.URL)
	r.MIME = cbc.NormalizeString(r.MIME)
}

func normalizeAssent(a *Assent) {
	uuid.Normalize(&a.UUID)
	a.Capacity = cbc.NormalizeString(a.Capacity)
	a.Statement = cbc.NormalizeString(a.Statement)
}

func normalizeAnalysis(a *Analysis) {
	uuid.Normalize(&a.UUID)
	a.Method = cbc.NormalizeCode(a.Method)
}

func normalizeAnnotation(a *Annotation) {
	a.Anchor = cbc.NormalizeString(a.Anchor)
	a.Summary = cbc.NormalizeString(a.Summary)
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
