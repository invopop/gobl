package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

// init registers the intrinsic normalizers for the pay package's types. The
// norm engine handles recursion into nested values and the application of
// regime/addon normalizers.
func init() {
	norm.Register("pay",
		norm.For(normalizeRecord),
		norm.For(normalizeInstructions),
		norm.For(normalizeTerms),
	)
}

func normalizeRecord(r *Record) {
	if r == nil {
		return
	}
	uuid.Normalize(&r.UUID)
	r.Ref = cbc.NormalizeString(r.Ref)
	r.Description = cbc.NormalizeString(r.Description)
	r.Ext = r.Ext.Clean()
}

func normalizeInstructions(i *Instructions) {
	if i == nil {
		return
	}
	i.Ref = cbc.NormalizeCode(i.Ref)
	i.Detail = cbc.NormalizeString(i.Detail)
	i.Notes = cbc.NormalizeString(i.Notes)
	i.Ext = i.Ext.Clean()
}

func normalizeTerms(t *Terms) {
	if t == nil {
		return
	}
	t.Notes = cbc.NormalizeString(t.Notes)
	t.Ext = t.Ext.Clean()
}
