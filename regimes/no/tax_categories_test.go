package no

import (
	"testing"

	"github.com/invopop/gobl/num"
)

func TestNorwayVATRates(t *testing.T) {
	if len(taxCategories) != 1 {
		t.Fatalf("expected 1 tax category, got %d", len(taxCategories))
	}
	cat := taxCategories[0]
	if len(cat.Rates) != 3 {
		t.Fatalf("expected 3 vat rates, got %d", len(cat.Rates))
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "general",
			got:  cat.Rates[0].Values[0].Percent.String(),
			want: num.MakePercentage(250, 3).String(),
		},
		{
			name: "reduced_food",
			got:  cat.Rates[1].Values[0].Percent.String(),
			want: num.MakePercentage(150, 3).String(),
		},
		{
			name: "reduced_transport",
			got:  cat.Rates[2].Values[0].Percent.String(),
			want: num.MakePercentage(120, 3).String(),
		},
	}

	for _, tc := range tests {
		if tc.got != tc.want {
			t.Fatalf("expected %s rate %s, got %s", tc.name, tc.want, tc.got)
		}
	}
}