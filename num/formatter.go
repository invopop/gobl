package num

import "strings"

const (
	// DefaultFormatterTemplate is the default template used to
	// format numbers for output.
	DefaultFormatterTemplate = "%n%u" // e.g. "12%"
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
}

// Format takes the provided amount and formats it according to
// the rules of the formatter.
func (f Formatter) Format(amount Amount) string {
	n := f.formatNumber(amount.String())
	t := f.Template
	if t == "" {
		t = DefaultFormatterTemplate
	}
	t = strings.Replace(t, "%u", f.Unit, 1)
	t = strings.Replace(t, "%n", n, 1)
	t = strings.TrimSpace(t)
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
	return n
}
