package org

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// init registers the intrinsic normalizers for the org package's types. Each
// function performs only the type's own normalization; the norm engine handles
// recursion into nested values and the application of regime/addon normalizers.
func init() {
	norm.Register("org",
		norm.For(normalizeParty),
		norm.For(normalizePerson),
		norm.For(normalizeItem),
		norm.For(normalizeAddress),
		norm.For(normalizeDocumentRef),
		norm.For(normalizeEmail),
		norm.For(normalizeIdentity),
		norm.For(normalizeInbox),
		norm.For(normalizeName),
		norm.For(normalizeNote),
		norm.For(normalizeRegistration),
		norm.For(normalizeTelephone),
		norm.For(normalizeWebsite),
		norm.For(normalizeAttachment),
	)
}

func normalizeParty(p *Party) {
	if p == nil {
		return
	}
	uuid.Normalize(&p.UUID)
	p.Label = cbc.NormalizeString(p.Label)
	p.Name = cbc.NormalizeString(p.Name)
	p.Alias = cbc.NormalizeString(p.Alias)
	p.Ext = p.Ext.Clean()
	if p.TaxID != nil {
		// tax ids are normalized only by their own tax regime, if any
		p.TaxID.Normalize()
	}
}

func normalizePerson(p *Person) {
	if p == nil {
		return
	}
	uuid.Normalize(&p.UUID)
	p.Label = cbc.NormalizeString(p.Label)
	p.Role = cbc.NormalizeString(p.Role)
}

func normalizeItem(i *Item) {
	if i == nil {
		return
	}
	i.Name = cbc.NormalizeString(i.Name)
	i.Description = cbc.NormalizeString(i.Description)
	i.Ref = cbc.NormalizeCode(i.Ref)
	i.Ext = i.Ext.Clean()
}

func normalizeAddress(a *Address) {
	if a == nil {
		return
	}
	uuid.Normalize(&a.UUID)
	a.PostOfficeBox = cbc.NormalizeString(a.PostOfficeBox)
	a.Number = cbc.NormalizeString(a.Number)
	a.Floor = cbc.NormalizeString(a.Floor)
	a.Block = cbc.NormalizeString(a.Block)
	a.Door = cbc.NormalizeString(a.Door)
	a.Street = cbc.NormalizeString(a.Street)
	a.StreetExtra = cbc.NormalizeString(a.StreetExtra)
	a.Locality = cbc.NormalizeString(a.Locality)
	a.Region = cbc.NormalizeString(a.Region)
	a.State = cbc.NormalizeAlphanumericalCode(a.State)
	a.Code = cbc.NormalizeCode(a.Code)
}

func normalizeDocumentRef(dr *DocumentRef) {
	if dr == nil {
		return
	}
	uuid.Normalize(&dr.UUID)
	dr.Series = cbc.NormalizeCode(dr.Series)
	dr.Code = cbc.NormalizeCode(dr.Code)
	dr.Reason = cbc.NormalizeString(dr.Reason)
	dr.URL = cbc.NormalizeString(dr.URL)
	dr.Ext = dr.Ext.Clean()
}

func normalizeEmail(e *Email) {
	if e == nil {
		return
	}
	uuid.Normalize(&e.UUID)
	e.Label = cbc.NormalizeString(e.Label)
	e.Address = cbc.NormalizeString(e.Address)
}

func normalizeIdentity(i *Identity) {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	i.Label = cbc.NormalizeString(i.Label)
	i.Type = cbc.NormalizeCode(i.Type)
	i.Code = cbc.NormalizeCode(i.Code)
	i.Description = cbc.NormalizeString(i.Description)
	i.Ext = i.Ext.Clean()
}

func normalizeInbox(i *Inbox) {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	code := i.Code.String()
	if is.EmailFormat.Check(code) {
		i.Email = code
		i.Code = ""
	} else if is.URL.Check(code) {
		i.URL = code
		i.Code = ""
	}
	i.Label = cbc.NormalizeString(i.Label)
	i.Scheme = cbc.NormalizeUpperCode(i.Scheme)
	i.Code = cbc.NormalizeCode(i.Code)

	// Custom normalizations
	switch i.Key {
	case InboxKeyPeppol:
		if i.Scheme == "" {
			if len(i.Code) >= 5 && i.Code[4] == ':' {
				numbers := i.Code[:4]
				i.Scheme = numbers
				i.Code = i.Code[5:]
			}
		}
	}
}

func normalizeName(n *Name) {
	if n == nil {
		return
	}
	uuid.Normalize(&n.UUID)
	n.Alias = cbc.NormalizeString(n.Alias)
	n.Prefix = cbc.NormalizeString(n.Prefix)
	n.Given = cbc.NormalizeString(n.Given)
	n.Middle = cbc.NormalizeString(n.Middle)
	n.Surname = cbc.NormalizeString(n.Surname)
	n.Surname2 = cbc.NormalizeString(n.Surname2)
	n.Suffix = cbc.NormalizeString(n.Suffix)
}

func normalizeNote(n *Note) {
	if n == nil {
		return
	}
	uuid.Normalize(&n.UUID)
	n.Code = cbc.NormalizeCode(n.Code)
	n.Text = cbc.NormalizeString(n.Text)
	n.Ext = n.Ext.Clean()
}

func normalizeRegistration(r *Registration) {
	if r == nil {
		return
	}
	uuid.Normalize(&r.UUID)
	r.Label = cbc.NormalizeString(r.Label)
	r.Office = cbc.NormalizeString(r.Office)
	r.Book = cbc.NormalizeString(r.Book)
	r.Volume = cbc.NormalizeString(r.Volume)
	r.Sheet = cbc.NormalizeString(r.Sheet)
	r.Section = cbc.NormalizeString(r.Section)
	r.Page = cbc.NormalizeString(r.Page)
	r.Entry = cbc.NormalizeString(r.Entry)
	r.Other = cbc.NormalizeString(r.Other)
}

func normalizeTelephone(t *Telephone) {
	if t == nil {
		return
	}
	uuid.Normalize(&t.UUID)
	t.Label = cbc.NormalizeString(t.Label)
	t.Number = strings.TrimSpace(t.Number)
}

func normalizeWebsite(w *Website) {
	if w == nil {
		return
	}
	uuid.Normalize(&w.UUID)
	w.Label = cbc.NormalizeString(w.Label)
	w.Title = cbc.NormalizeString(w.Title)
	w.URL = cbc.NormalizeString(w.URL)
}

func normalizeAttachment(a *Attachment) {
	if a == nil {
		return
	}
	uuid.Normalize(&a.UUID)
	a.Code = cbc.NormalizeCode(a.Code)
	a.Name = cbc.NormalizeString(a.Name)
	a.Description = cbc.NormalizeString(a.Description)
	a.URL = cbc.NormalizeString(a.URL)
	a.MIME = cbc.NormalizeString(a.MIME)
}
