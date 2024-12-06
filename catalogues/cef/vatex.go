package cef

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyVATEX is used for the CEF VATEX exemption codes.
	ExtKeyVATEX cbc.Key = "cef-vatex"
)

var extVATEX = &cbc.Definition{
	Key:  ExtKeyVATEX,
	Name: i18n.NewString("CET VATEX - VAT exemption reason codes"),
	Desc: i18n.NewString(here.Doc(`
		Codes for the reasons for VAT exemption as defined by the Connecting Europe Facility (CEF).
	`)),
	Values: []*cbc.Definition{
		{
			Code: "VATEX-EU-79-C",
			Name: i18n.NewString("Exempt based on article 79, point c of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132",
			Name: i18n.NewString("Exempt based on article 132 of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1A",
			Name: i18n.NewString("Exempt based on article 132, section 1 (a) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1B",
			Name: i18n.NewString("Exempt based on article 132, section 1 (b) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1C",
			Name: i18n.NewString("Exempt based on article 132, section 1 (c) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1D",
			Name: i18n.NewString("Exempt based on article 132, section 1 (d) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1E",
			Name: i18n.NewString("Exempt based on article 132, section 1 (e) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1F",
			Name: i18n.NewString("Exempt based on article 132, section 1 (f) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1G",
			Name: i18n.NewString("Exempt based on article 132, section 1 (g) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1H",
			Name: i18n.NewString("Exempt based on article 132, section 1 (h) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1I",
			Name: i18n.NewString("Exempt based on article 132, section 1 (i) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1J",
			Name: i18n.NewString("Exempt based on article 132, section 1 (j) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1K",
			Name: i18n.NewString("Exempt based on article 132, section 1 (k) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1L",
			Name: i18n.NewString("Exempt based on article 132, section 1 (l) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1M",
			Name: i18n.NewString("Exempt based on article 132, section 1 (m) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1N",
			Name: i18n.NewString("Exempt based on article 132, section 1 (n) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1O",
			Name: i18n.NewString("Exempt based on article 132, section 1 (o) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1P",
			Name: i18n.NewString("Exempt based on article 132, section 1 (p) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-132-1Q",
			Name: i18n.NewString("Exempt based on article 132, section 1 (q) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143",
			Name: i18n.NewString("Exempt based on article 143 of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1A",
			Name: i18n.NewString("Exempt based on article 143, section 1 (a) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1B",
			Name: i18n.NewString("Exempt based on article 143, section 1 (b) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1C",
			Name: i18n.NewString("Exempt based on article 143, section 1 (c) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1D",
			Name: i18n.NewString("Exempt based on article 143, section 1 (d) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1E",
			Name: i18n.NewString("Exempt based on article 143, section 1 (e) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1F",
			Name: i18n.NewString("Exempt based on article 143, section 1 (f) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1FA",
			Name: i18n.NewString("Exempt based on article 143, section 1 (fa) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1G",
			Name: i18n.NewString("Exempt based on article 143, section 1 (g) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1H",
			Name: i18n.NewString("Exempt based on article 143, section 1 (h) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1I",
			Name: i18n.NewString("Exempt based on article 143, section 1 (i) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1J",
			Name: i18n.NewString("Exempt based on article 143, section 1 (j) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1K",
			Name: i18n.NewString("Exempt based on article 143, section 1 (k) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-143-1L",
			Name: i18n.NewString("Exempt based on article 143, section 1 (l) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148",
			Name: i18n.NewString("Exempt based on article 148 of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-A",
			Name: i18n.NewString("Exempt based on article 148, section (a) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-B",
			Name: i18n.NewString("Exempt based on article 148, section (b) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-C",
			Name: i18n.NewString("Exempt based on article 148, section (c) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-D",
			Name: i18n.NewString("Exempt based on article 148, section (d) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-E",
			Name: i18n.NewString("Exempt based on article 148, section (e) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-F",
			Name: i18n.NewString("Exempt based on article 148, section (f) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-148-G",
			Name: i18n.NewString("Exempt based on article 148, section (g) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151",
			Name: i18n.NewString("Exempt based on article 151 of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1A",
			Name: i18n.NewString("Exempt based on article 151, section 1 (a) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1AA",
			Name: i18n.NewString("Exempt based on article 151, section 1 (aa) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1B",
			Name: i18n.NewString("Exempt based on article 151, section 1 (b) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1C",
			Name: i18n.NewString("Exempt based on article 151, section 1 (c) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1D",
			Name: i18n.NewString("Exempt based on article 151, section 1 (d) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-151-1E",
			Name: i18n.NewString("Exempt based on article 151, section 1 (e) of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-309",
			Name: i18n.NewString("Exempt based on article 309 of Council Directive 2006/112/EC"),
		},
		{
			Code: "VATEX-EU-AE",
			Name: i18n.NewString("Reverse charge"),
		},
		{
			Code: "VATEX-EU-D",
			Name: i18n.NewString("Travel agents VAT scheme."),
		},
		{
			Code: "VATEX-EU-F",
			Name: i18n.NewString("Second hand goods VAT scheme."),
		},
		{
			Code: "VATEX-EU-G",
			Name: i18n.NewString("Export outside the EU"),
		},
		{
			Code: "VATEX-EU-I",
			Name: i18n.NewString("Works of art VAT scheme."),
		},
		{
			Code: "VATEX-EU-IC",
			Name: i18n.NewString("Intra-community supply"),
		},
		{
			Code: "VATEX-EU-J",
			Name: i18n.NewString("Collectors items and antiques VAT scheme."),
		},
		{
			Code: "VATEX-EU-O",
			Name: i18n.NewString("Not subject to VAT"),
		},
		{
			Code: "VATEX-FR-FRANCHISE",
			Name: i18n.NewString("France domestic VAT franchise in base"),
		},
		{
			Code: "VATEX-FR-CNWVAT",
			Name: i18n.NewString("France domestic Credit Notes without VAT, due to supplier forfeit of VAT for discount"),
		},
	},
}
