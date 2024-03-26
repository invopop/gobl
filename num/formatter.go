package num

import "strings"

const (
	// DefaultFormatterTemplate is the default template used to
	// format numbers for output.
	DefaultFormatterTemplate = "%n%u" // e.g. "12%"
)

// NumeralSystem describes how to output numbers.
type NumeralSystem string

const (
	// NumeralWestern defines the Western Arabic numeral
	// system, and is the default.
	NumeralWestern NumeralSystem = "western"
	// NumeralArabic can be used to output numbers using the
	// Arabic numeral system.
	NumeralArabic NumeralSystem = "arabic"
)

// Formatter is used to define how an amount should be formatted
// alongside a unit if necessary.
type Formatter struct {
	// DecimalMark is the character used to separate the whole
	// number from the decimal part.
	DecimalMark string
	// ThousandsSeparator is the character used to separate
	// thousands in the whole number.
	ThousandsSeparator string
	// Unit is the string representation or symbol of the unit.
	Unit string
	// Template is the string used to present the number and unit
	// together with two simple placeholders, `%n` for the number and
	// `%u` for the unit.
	Template string
	// NumeralSystem determines how numbers should be output. By default
	// this is 'western'.
	NumeralSystem NumeralSystem
}

var (
	formatArabicReplacer = strings.NewReplacer(
		"0", "٠",
		"1", "١",
		"2", "٢",
		"3", "٣",
		"4", "٤",
		"5", "٥",
		"6", "٦",
		"7", "٧",
		"8", "٨",
		"9", "٩",
	)
)

// MakeFormatter prepares a new formatter with the two main configuration
// options, decimal and thousands separators.
func MakeFormatter(decimalMark, thousandsSeparator string) Formatter {
	return Formatter{
		DecimalMark:        decimalMark,
		ThousandsSeparator: thousandsSeparator,
	}
}

// WithUnit providers a formatter with a unit set.
func (f Formatter) WithUnit(unit string) Formatter {
	f.Unit = unit
	return f
}

// WithoutUnit provides a formatter without a unit set.
func (f Formatter) WithoutUnit() Formatter {
	f.Unit = ""
	return f
}

// WithTemplate sets the template for use with formatting with
// units.
func (f Formatter) WithTemplate(template string) Formatter {
	f.Template = template
	return f
}

// WithNumeralSystem overrides the default western numeral system
// with that defined.
func (f Formatter) WithNumeralSystem(ns NumeralSystem) Formatter {
	f.NumeralSystem = ns
	return f
}

// Amount takes the provided amount and formats it according to
// the rules of the formatter.
func (f Formatter) Amount(amount Amount) string {
	n := f.formatNumber(amount.String())
	return f.formatWithUnits(n)
}

// Percentage tries to format the percentage value according to the
// rules of the formatter, but replacing the unit with a percentage
// symbol, and using the default template.
func (f Formatter) Percentage(percent Percentage) string {
	n := f.formatNumber(percent.StringWithoutSymbol())
	return f.WithUnit("%").WithTemplate("").formatWithUnits(n)
}

func (f Formatter) formatWithUnits(n string) string {
	if f.Unit == "" {
		return n
	}
	t := f.Template
	if t == "" {
		t = DefaultFormatterTemplate
	}
	t = strings.Replace(t, "%u", f.Unit, 1)
	t = strings.Replace(t, "%n", n, 1)
	return t
}

func (f Formatter) formatNumber(n string) string {
	p := strings.Split(n, ".")
	n = p[0]
	// split the main part with thousands separator
	for i := len(n) - 3; i > 0; i = i - 3 {
		n = n[:i] + f.ThousandsSeparator + n[i:]
	}
	if len(p) == 2 {
		n = n + f.DecimalMark + p[1]
	}
	n = f.formatNumberForSystem(n)
	return n
}

func (f Formatter) formatNumberForSystem(n string) string {
	switch f.NumeralSystem {
	case NumeralArabic:
		// Support for Arabic numbers is somewhat experimental, feedback welcome!
		return formatArabicReplacer.Replace(n)
	default:
		return n
	}
}
