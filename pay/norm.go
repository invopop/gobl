package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

// init registers the intrinsic normalizers for the pay package's types. The
// norm engine handles recursion, the global cleaning of codes and extensions,
// and the application of regime/addon normalizers.
func init() {
	norm.Register(
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
}

func normalizeInstructions(i *Instructions) {
	if i == nil {
		return
	}
	i.Detail = cbc.NormalizeString(i.Detail)
	i.Notes = cbc.NormalizeString(i.Notes)
}

func normalizeTerms(t *Terms) {
	if t == nil {
		return
	}
	t.Notes = cbc.NormalizeString(t.Notes)
}
