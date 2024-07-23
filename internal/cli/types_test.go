package cli

import (
	"reflect"
	"testing"

	"github.com/invopop/gobl/schema"
)

func TestFindType(t *testing.T) {
	const (
		idInvoice = "https://gobl.org/draft-0/foo/invoice"
	)
	type Invoice struct{}
	r := map[reflect.Type]schema.ID{
		reflect.TypeOf(Invoice{}): idInvoice,
	}

	t.Run("exact schema match", func(t *testing.T) {
		got := findType(r, idInvoice)
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("exact type match", func(t *testing.T) {
		got := findType(r, "Invoice")
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("exact type match with package", func(t *testing.T) {
		got := findType(r, "cli.Invoice")
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("wrong package", func(t *testing.T) {
		got := findType(r, "wrongpkg.Invoice")
		if got != "" {
			t.Errorf("Unexpected result: %v", got)
		}
	})
	t.Run("implicit schema", func(t *testing.T) {
		got := findType(r, "foo.Invoice")
		if got != idInvoice {
			t.Errorf("Unexpected result: %v", got)
		}
	})
}

func Test_toSchema(t *testing.T) {
	tests := map[string]string{
		"bill.Invoice":          "/bill/invoice",
		"bill.SomeOtherInvoice": "/bill/some-other-invoice",
		"bill.NASAInvoice":      "/bill/nasa-invoice",
		"https://full/schema":   "https://full/schema",
		"bill.123Invoice":       "/bill/123-invoice",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got := toSchema(input)
			if got != want {
				t.Errorf("Unexpected result: %q", got)
			}
		})
	}
}
