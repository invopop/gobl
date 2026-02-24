package no

import "testing"

func TestNorwayVATRates(t *testing.T) {
	if len(taxCategories) != 1 {
		t.Fatalf("expected 1 tax category, got %d", len(taxCategories))
	}
	cat := taxCategories[0]
	if len(cat.Rates) != 3 {
		t.Fatalf("expected 3 vat rates, got %d", len(cat.Rates))
	}

	if cat.Rates[0].Values[0].Percent.String() != "25.0%" {
		t.Fatalf("expected standard rate 25.0%%, got %s", cat.Rates[0].Values[0].Percent.String())
	}
	if cat.Rates[1].Values[0].Percent.String() != "15.0%" {
		t.Fatalf("expected reduced rate 15.0%%, got %s", cat.Rates[1].Values[0].Percent.String())
	}
	if cat.Rates[2].Values[0].Percent.String() != "12.0%" {
		t.Fatalf("expected reduced rate 12.0%%, got %s", cat.Rates[2].Values[0].Percent.String())
	}
}
