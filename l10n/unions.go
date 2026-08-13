package l10n

import "github.com/invopop/gobl/cal"

// Union codes
const (
	EU Code = "EU"
)

// Unions provides of list of significant political and economic
// unions that countries may be members of.
func Unions() UnionDefs {
	return unions
}

// Union is a convenience function to return the Union definition
// for a given code.
func Union(code Code) *UnionDef {
	return unions.Code(code)
}

var unions = UnionDefs{
	{
		Code: EU,
		Name: "European Union",
		Members: []*UnionMember{
			{
				Code:  AT,
				Since: cal.MakeDate(1995, 1, 1),
			},
			{
				Code:  BE,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  BG,
				Since: cal.MakeDate(2007, 1, 1),
			},
			{
				Code:  CY,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  CZ,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  DE,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  DK,
				Since: cal.MakeDate(1973, 1, 1),
			},
			{
				Code:  EE,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  ES,
				Since: cal.MakeDate(1986, 1, 1),
			},
			{
				Code:  FI,
				Since: cal.MakeDate(1995, 1, 1),
			},
			{
				Code:  FR,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  GB,
				Since: cal.MakeDate(1973, 1, 1),
				Until: cal.MakeDate(2020, 1, 31),
			},
			{
				Code:    GR,
				AltCode: EL,
				Since:   cal.MakeDate(1981, 1, 1),
			},
			{
				Code:  HR,
				Since: cal.MakeDate(2013, 1, 1),
			},
			{
				Code:  HU,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  IE,
				Since: cal.MakeDate(1973, 1, 1),
			},
			{
				Code:  IT,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  LT,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  LU,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  LV,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  MT,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  NL,
				Since: cal.MakeDate(1958, 1, 1),
			},
			{
				Code:  PL,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  PT,
				Since: cal.MakeDate(1986, 1, 1),
			},
			{
				Code:  RO,
				Since: cal.MakeDate(2007, 1, 1),
			},
			{
				Code:  SE,
				Since: cal.MakeDate(1995, 1, 1),
			},
			{
				Code:  SI,
				Since: cal.MakeDate(2004, 1, 1),
			},
			{
				Code:  SK,
				Since: cal.MakeDate(2004, 1, 1),
			},
		},
	},
}
