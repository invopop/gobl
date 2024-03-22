package currency

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

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
	// Format determines how to layout the units and number
	Format string `json:"format"`
	// HTML entity code for the symbol
	HTMLEntity string `json:"html_entity"`
	// Decimal mark normally expected in output
	DecimalMark string `json:"decimal_mark"`
	// Thousands separator normally expected in output
	ThousandsSeparator string `json:"thousands_separator"`
	// Smallest acceptable amount of the currency
	SmallestDenomination int `json:"smallest_denomination"`
}

// Amount takes the provided amount and formats it according
// to the rules of the currency definition.
func (d *Def) Amount(amount num.Amount) string {
	n := d.formatNumber(amount.String())
	f := d.Format
	if f == "" {
		f = "%u%n"
	}
	f = strings.Replace(f, "%u", d.Symbol, 1)
	f = strings.Replace(f, "%n", n, 1)
	return f
}

// Percentage formats a percentage according to the rules
// of the currency, but without a currency symbol.
func (d *Def) Percentage(percent num.Percentage) string {
	n := d.formatNumber(percent.StringWithoutSymbol())
	return n + "%"
}

func (d *Def) formatNumber(n string) string {
	p := strings.Split(n, ".")
	n = p[0]
	// split the main part with thousands separator
	for i := len(n) - 3; i > 0; i = i - 3 {
		n = n[:i] + d.ThousandsSeparator + n[i:]
	}
	if len(p) == 2 {
		n = n + d.DecimalMark + p[1]
	}
	return n
}

// Zero provides the currency's zero amount which is pre-set with the
// minimum precision for the currency.
func (d *Def) Zero() num.Amount {
	return num.MakeAmount(0, d.Subunits)
}

// Defs provides an array of all currency definitions
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

	err := fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
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
	sort.Slice(ds.byPriority, func(i, j int) bool {
		return ds.byPriority[i].Priority < ds.byPriority[j].Priority
	})

	return nil
}
