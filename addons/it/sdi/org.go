package sdi

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func orgAddressRules() *rules.Set {
	return rules.For(new(org.Address),
		rules.Assert("01", "address either street or post office box must be set",
			is.Func("street or postbox", addressHasStreetOrPostBox),
		),
		rules.Field("street",
			rules.Assert("02", "address street must use Latin-1 characters",
				is.FuncError("latin1", validateLatin1String),
			),
		),
		rules.Field("po_box",
			rules.Assert("03", "address post office box must use Latin-1 characters",
				is.FuncError("latin1", validateLatin1String),
			),
		),
		rules.Field("country",
			rules.Assert("04", "address country is required", is.Present),
		),
		rules.Field("locality",
			rules.Assert("05", "address locality is required", is.Present),
			rules.Assert("06", "address locality must use Latin-1 characters",
				is.FuncError("latin1", validateLatin1String),
			),
		),
		rules.When(is.Func("Italian address", addressIsItalian),
			rules.Field("code",
				rules.Assert("07", "Italian address code is required", is.Present),
				rules.Assert("08", "Italian address code must be 5 digits", is.Matches(`^\d{5}$`)),
			),
		),
	)
}

func addressHasStreetOrPostBox(val any) bool {
	a, ok := val.(*org.Address)
	if !ok || a == nil {
		return true
	}
	return a.Street != "" || a.PostOfficeBox != ""
}

func addressIsItalian(val any) bool {
	a, ok := val.(*org.Address)
	if !ok || a == nil {
		return false
	}
	return a.Country.In("IT")
}
