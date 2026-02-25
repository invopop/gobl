package no

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

func TestNormalizeTaxIdentityCases(t *testing.T) {
	tests := []struct {
		name string
		in   cbc.Code
		want cbc.Code
	}{
		{
			name: "strips_no_prefix_spaces_and_mva_suffix",
			in:   cbc.Code(" NO 974 760 673 mva "),
			want: cbc.Code("974760673"),
		},
		{
			name: "strips_no_prefix_and_mva_suffix",
			in:   cbc.Code("NO974760673MVA"),
			want: cbc.Code("974760673"),
		},
		{
			name: "digits_only_kept",
			in:   cbc.Code("974760673"),
			want: cbc.Code("974760673"),
		},
		{
			name: "unknown_format_unchanged",
			in:   cbc.Code("SOMETHING"),
			want: cbc.Code("SOMETHING"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    tc.in,
			}

			normalizeTaxIdentity(id)

			if id.Code != tc.want {
				t.Fatalf("expected normalized code to be %s, got %s", tc.want, id.Code)
			}
		})
	}
}