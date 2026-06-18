package l10n

import "github.com/invopop/gobl/cal"

// Union codes
const (
	EU   Code = "EU"
	SEPA Code = "SEPA"
)

// sepaSince is the SEPA scheme inception date (launch of the SEPA Credit
// Transfer scheme). Per-country onboarding happened on later, varying dates,
// but those are not tracked here: membership is used to answer "is this country
// in SEPA?" for current documents, so a single founding date is sufficient.
// Refine individual members with their accession date if historical precision
// is ever required.
var sepaSince = cal.MakeDate(2008, 1, 28)

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
	{
		Code: SEPA,
		Name: "Single Euro Payments Area",
		Members: []*UnionMember{
			// EU member states
			{Code: AT, Since: sepaSince},
			{Code: BE, Since: sepaSince},
			{Code: BG, Since: sepaSince},
			{Code: HR, Since: sepaSince},
			{Code: CY, Since: sepaSince},
			{Code: CZ, Since: sepaSince},
			{Code: DK, Since: sepaSince},
			{Code: EE, Since: sepaSince},
			{Code: FI, Since: sepaSince},
			{Code: FR, Since: sepaSince},
			{Code: DE, Since: sepaSince},
			{Code: GR, AltCode: EL, Since: sepaSince}, // tax code EL, ISO GR
			{Code: HU, Since: sepaSince},
			{Code: IE, Since: sepaSince},
			{Code: IT, Since: sepaSince},
			{Code: LV, Since: sepaSince},
			{Code: LT, Since: sepaSince},
			{Code: LU, Since: sepaSince},
			{Code: MT, Since: sepaSince},
			{Code: NL, Since: sepaSince},
			{Code: PL, Since: sepaSince},
			{Code: PT, Since: sepaSince},
			{Code: RO, Since: sepaSince},
			{Code: SK, Since: sepaSince},
			{Code: SI, Since: sepaSince},
			{Code: ES, Since: sepaSince},
			{Code: SE, Since: sepaSince},
			// Non-EU SEPA participants: EEA, plus other states and
			// territories that have joined the SEPA schemes.
			{Code: IS, Since: sepaSince},
			{Code: LI, Since: sepaSince},
			{Code: NO, Since: sepaSince},
			{Code: CH, Since: sepaSince},
			{Code: GB, Since: sepaSince},
			{Code: MC, Since: sepaSince},
			{Code: SM, Since: sepaSince},
			{Code: AD, Since: sepaSince},
			{Code: VA, Since: sepaSince},
			{Code: JE, Since: sepaSince},
			{Code: GG, Since: sepaSince},
			{Code: IM, Since: sepaSince},
			{Code: GI, Since: sepaSince},
		},
	},
}
