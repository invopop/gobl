package currency

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/invopop/gobl/num"
	"github.com/invopop/yaml"
)

var definitions *defs

// defs is used internally to load and access currency
// definitions.
type defs struct {
	byCode     map[Code]*Def
	byPriority []*Def
}

const (
	// DefaultCurrencyTemplate defines how to output currencies for most
	// common use cases.
	DefaultCurrencyTemplate = "%u%n"
)

// Def helps define how to format a currency as is based on the
// [Ruby Money Gem's](https://rubymoney.github.io/money/) Currency model.
type Def struct {
	// Priority is an arbitrary number used for ordering of currencies
	// roughly based on popularity.
	Priority int `json:"priority"`
	// Standard ISO 4217 code
	ISOCode Code `json:"iso_code"`
	// ISO numeric code
	ISONumeric string `json:"iso_numeric"`
	// English name of the currency
	Name string `json:"name"`
	// Symbol representation
	Symbol string `json:"symbol"`
	// When presented alongside other currency's with potentially
	// the same symbol, use this representation instead.
	DisambiguateSymbol string `json:"disambiguate_symbol"`
	// Alternative presentation symbols
	AlternateSymbols []string `json:"alternate_symbols"`
	// Name of the currency subunit
	SubunitName string `json:"subunit_name"`
	// Conversion amount to subunit
	Subunits uint32 `json:"subunits"`
	// Template determines how to layout the units and number
	Template string `json:"template"`
	// Decimal mark normally expected in output
	DecimalMark string `json:"decimal_mark"`
	// Thousands separator normally expected in output
	ThousandsSeparator string `json:"thousands_separator"`
	// Smallest acceptable amount of the currency
	SmallestDenomination int `json:"smallest_denomination"`
	// NumeralSystem defines how numbers should be printed out, by default this
	// is 'western'.
	NumeralSystem num.NumeralSystem `json:"numeral_system"`
}

// FormatOption defines how to configure the formatter for common
// use cases and custom options.
type FormatOption func(*Def, num.Formatter) num.Formatter

// WithDisambiguateSymbol will override the default symbol to use with one that
// is unique for the context. Lots of countries for example use "$" as their
// main currency symbol, using this option will ensure that `US$` is used
// in output instead.
func WithDisambiguateSymbol() FormatOption {
	return func(d *Def, f num.Formatter) num.Formatter {
		f.Unit = d.DisambiguateSymbol
		if f.Unit == "" {
			// fall back to symbol
			f.Unit = d.Symbol
		}
		return f
	}
}

// WithNumeralSystem will override the default numeral system used to output
// numbers.
func WithNumeralSystem(ns num.NumeralSystem) FormatOption {
	return func(_ *Def, f num.Formatter) num.Formatter {
		f.NumeralSystem = ns
		return f
	}
}

// Formatter provides a number formatter for the currency definition.
func (d *Def) Formatter(opts ...FormatOption) num.Formatter {
	f := num.Formatter{
		DecimalMark:        d.DecimalMark,
		ThousandsSeparator: d.ThousandsSeparator,
		Unit:               d.Symbol,
		Template:           d.Template,
	}
	if d.Template == "" {
		f.Template = DefaultCurrencyTemplate
	}
	for _, opt := range opts {
		f = opt(d, f)
	}
	return f
}

// FormatAmount takes the provided amount and formats it according
// to the default rules of the currency definition.
func (d *Def) FormatAmount(amount num.Amount) string {
	return d.Formatter().Amount(amount)
}

// FormatPercentage takes the provided percentage and formats it
// according to the decimal and thousands rules of the currency
// definition.
func (d *Def) FormatPercentage(percentage num.Percentage) string {
	return d.Formatter().Percentage(percentage)
}

// Zero provides the currency's zero amount which is pre-set with the
// minimum precision for the currency.
func (d *Def) Zero() num.Amount {
	return num.MakeAmount(0, d.Subunits)
}

// Rescale takes the provided amount and ensures its scale matches
// that of the currency.
func (d *Def) Rescale(a num.Amount) num.Amount {
	return a.Rescale(d.Subunits)
}

// RescaleUp ensures tha the amount has *at least* the same
// precision as the currency.
func (d *Def) RescaleUp(a num.Amount) num.Amount {
	return a.RescaleUp(d.Subunits)
}

// Definitions provides an array of all currency definitions
// ordered by priority.
func Definitions() []*Def {
	return definitions.all()
}

// Get provides the code's currency definition, or
// nil if no match found.
func Get(c Code) *Def {
	return definitions.get(c)
}

// ByISONumber tries to find the currency definition by it's assigned ISO
// number.
func ByISONumber(n string) *Def {
	for _, def := range definitions.all() {
		if def.ISONumeric == n {
			return def
		}
	}
	return nil
}

func (ds *defs) all() []*Def {
	return ds.byPriority
}

func (ds *defs) get(c Code) *Def {
	if def, ok := ds.byCode[c]; ok {
		return def
	}
	return nil
}

// load attempts to ready all of the currency JSON definition
// files from the provided source and load them into memory.
func (ds *defs) load(src fs.FS, root string) error {
	ds.byPriority = make([]*Def, 0)
	ds.byCode = make(map[Code]*Def)

	err := fs.WalkDir(src, root, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking directory: %w", err)
		}

		switch filepath.Ext(path) {
		case ".yaml", ".yml", ".json":
			// good
		default:
			return nil
		}

		data, err := fs.ReadFile(src, path)
		if err != nil {
			return fmt.Errorf("reading file '%s': %w", path, err)
		}

		list := make([]*Def, 0)
		if err := yaml.Unmarshal(data, &list); err != nil {
			return fmt.Errorf("unmarshalling file '%s': %w", path, err)
		}

		for _, def := range list {
			if _, ok := ds.byCode[def.ISOCode]; ok {
				return fmt.Errorf("duplicate currency code: %s", def.ISOCode)
			}
			ds.byCode[def.ISOCode] = def
			ds.byPriority = append(ds.byPriority, def)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Sort the byPriority list by the priority field
	// in ascending order.
	slices.SortStableFunc(ds.byPriority, func(def1, def2 *Def) int {
		return def1.Priority - def2.Priority
	})

	return nil
}
