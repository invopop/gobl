package no

import (
	"testing"

	"github.com/invopop/gobl/cbc"
)

func TestCleanNorwayOrgNr(t *testing.T) {
	d, ok := cleanNorwayOrgNr(cbc.Code("974 760 673"))
	if !ok || d != "974760673" {
		t.Fatalf("expected 974760673,true got %s,%v", d, ok)
	}

	_, ok = cleanNorwayOrgNr(cbc.Code("123"))
	if ok {
		t.Fatalf("expected ok=false for short code")
	}
}

func TestCleanNorwayTaxCode(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"974760673", "974760673", true},
		{"974760673MVA", "974760673", true},
		{"NO974760673MVA", "974760673", true},
		{" no 974 760 673 mva ", "974760673", true},
		{"NOMVA", "", false},
	}

	for _, tc := range cases {
		got, ok := cleanNorwayTaxCode(cbc.Code(tc.in))
		if ok != tc.ok || (ok && got != tc.want) {
			t.Fatalf("input=%q expected %q,%v got %q,%v", tc.in, tc.want, tc.ok, got, ok)
		}
	}
}
