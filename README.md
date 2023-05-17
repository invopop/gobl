# GOBL

<img src="https://github.com/invopop/gobl/blob/main/gobl_logo_black_rgb.svg?raw=true" width="181" height="219" alt="GOBL Logo">

Go Business Language. Core library and Schemas.

Released under the Apache 2.0 [LICENSE](https://github.com/invopop/gobl/blob/main/LICENSE), Copyright 2021,2022 [Invopop Ltd.](https://invopop.com).

[![Lint](https://github.com/invopop/gobl/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl)](https://goreportcard.com/report/github.com/invopop/gobl)
[![GoDoc](https://godoc.org/github.com/invopop/gobl?status.svg)](https://godoc.org/github.com/invopop/gobl)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl)

[Official GOBL documentation site](https://docs.gobl.org).

## Introduction

GOBL, the Go Business Language library and tools, aim to:

- Help developers build electronic business documents, especially invoices, anywhere in the world.
- Define a a set of open [JSON Schema](https://json-schema.org/) with the flexibility to be used and shared.
- Build a global database of local tax categories and, whenever practical to do so, provide current and historical tax rates in code.
- Validate business documents according to local requirements, including tax ID validation.
- Define the algorithms used to make tax calculations while avoiding rounding errors.
- Provide built-in support for signing documents using [JSON Web Signatures](https://en.wikipedia.org/wiki/JSON_Web_Signature).
- Output simple and easy to read JSON documents that emphasize the use of keys instead of abstract codes, like `credit-transfer` instead of `30` (UNTDID4461 code for sender initiated bank or wire transfer).
- Be flexible enough to support extreme local complexity but produce output that is easily legible in other countries.
- Build a global community of contributors tired of the complexity of current standards based on XML or EDI.

## Companion Projects

GOBL makes it easy to create business documents, like invoices, but checkout some of the companion projects that help create, use, and convert into other formats:

- [CLI](https://github.com/invopop/gobl.cli) - the official GOBL command line tool, including WASM release for streaming in browsers.
- [Builder](https://github.com/invopop/gobl.builder) - Available to try at [build.gobl.org](https://build.gobl.org), this tool makes it easy to build, test, and discover the features of GOBL.
- [Generator](https://github.com/invopop/gobl.generator) - Ruby project to convert GOBL JSON Schema into libraries for other languages or documentation.
- [Docs](https://github.com/invopop/gobl.docs) - Content of the official GOBL Documentation Site [docs.gobl.org](https://docs.gobl.org).
- [Ruby](https://github.com/invopop/gobl.ruby) - Easily build or read GOBL documents in Ruby.
- [FacturaE](https://github.com/invopop/gobl.facturae) - convert into the [Spanish FacturaE](https://www.facturae.gob.es/Paginas/Index.aspx) format.
- [FatturaPA](https://github.com/invopop/gobl.fatturapa) - convert into the [Italian FatturaPA](https://www.fatturapa.gov.it/it/index.html) format.

## Usage

GOBL is a Go library, so the following instructions assume you'd like to build documents from your own Go applications. See some of the links above if you'd like to develop in another language or use a CLI.

### Installation

Run the following command to install the package:

```
go get github.com/invopop/gobl
```

### Building an Invoice

There are lots of different ways to get data into GOBL but for the following example we're going to try and build an invoice in several steps.

First define a minimal or "partial" GOBL Invoice Document:

```go
inv := &bill.Invoice{
	Series:    "F23",
	Code:      "00010",
	IssueDate: cal.MakeDate(2023, time.May, 11),
	Supplier: &org.Party{
		TaxID: &tax.Identity{
			Country: l10n.US,
		},
		Name:  "Provider One Inc.",
		Alias: "Provider One",
		Emails: []*org.Email{
			{
				Address: "billing@provideone.com",
			},
		},
		Addresses: []*org.Address{
			{
				Number:   "16",
				Street:   "Jessie Street",
				Locality: "San Francisco",
				Region:   "CA",
				Code:     "94105",
				Country:  l10n.US,
			},
		},
	},
	Customer: &org.Party{
		Name: "Sample Customer",
		Emails: []*org.Email{
			{
				Address: "email@sample.com",
			},
		},
	},
	Lines: []*bill.Line{
		{
			Quantity: num.MakeAmount(20, 0),
			Item: &org.Item{
				Name:  "A stylish mug",
				Price: num.MakeAmount(2000, 2),
				Unit:  org.UnitHour,
			},
			Taxes: []*tax.Combo{
				{
					Category: common.TaxCategoryST,
					Percent:  num.NewPercentage(85, 3),
				},
			},
		},
	},
}
```

Notice that the are no sums or calculations yet. The next step involves "inserting" the invoice document into an "envelope". In GOBL, we use the concept of an envelope to hold data and provide functionality to guarantee that no modifications have been made to the payload.

Insert our previous Invoice into an envelope as follows:

```go
// Prepare an "Envelope"
env := gobl.NewEnvelope()
if err := env.Insert(inv); err != nil {
	panic(err)
}
```

