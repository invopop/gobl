package l10n

import "github.com/invopop/gobl/cal"

// Union codes
const (
	EU   Code = "EU"
	SEPA Code = "SEPA"
)

// sepaLaunch is the date the SEPA Credit Transfer scheme went live. The EPC
// (see SEPA union below) does not publish individual entry dates for the
// founding cohort — only that the early non-EEA members were resolved "between
// March 2006 and December 2013" — so this launch date is used as their Since.
// Members added later carry their documented effective date instead.
var sepaLaunch = cal.MakeDate(2008, 1, 28)

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
		// SEPA geographical scope per the EPC "List of SEPA Scheme Countries"
		// (EPC409-09, v8.0, 24 December 2025). Dates are the documented entry
		// dates where the EPC gives them; the founding cohort uses the scheme
		// launch (see sepaLaunch) as the EPC does not date them individually.
		// EU sub-territories (Åland, Azores, Canary Islands, the French overseas
		// departments/collectivities, etc.) are covered by their parent country
		// code and are not listed separately.
		Code: SEPA,
		Name: "Single Euro Payments Area",
		Members: []*UnionMember{
			// EU member states (in at, or via EU accession; founding cohort
			// uses the SEPA launch date).
			{Code: AT, Since: sepaLaunch},
			{Code: BE, Since: sepaLaunch},
			{Code: BG, Since: sepaLaunch},
			{Code: HR, Since: cal.MakeDate(2013, 7, 1)}, // EU accession; not in SEPA at 2008 launch
			{Code: CY, Since: sepaLaunch},
			{Code: CZ, Since: sepaLaunch},
			{Code: DK, Since: sepaLaunch},
			{Code: EE, Since: sepaLaunch},
			{Code: FI, Since: sepaLaunch},
			{Code: FR, Since: sepaLaunch},
			{Code: DE, Since: sepaLaunch},
			{Code: GR, AltCode: EL, Since: sepaLaunch}, // tax code EL, ISO GR
			{Code: HU, Since: sepaLaunch},
			{Code: IE, Since: sepaLaunch},
			{Code: IT, Since: sepaLaunch},
			{Code: LV, Since: sepaLaunch},
			{Code: LT, Since: sepaLaunch},
			{Code: LU, Since: sepaLaunch},
			{Code: MT, Since: sepaLaunch},
			{Code: NL, Since: sepaLaunch},
			{Code: PL, Since: sepaLaunch},
			{Code: PT, Since: sepaLaunch},
			{Code: RO, Since: sepaLaunch},
			{Code: SK, Since: sepaLaunch},
			{Code: SI, Since: sepaLaunch},
			{Code: ES, Since: sepaLaunch},
			{Code: SE, Since: sepaLaunch},
			// EEA, non-EU.
			{Code: IS, Since: sepaLaunch},
			{Code: LI, Since: sepaLaunch},
			{Code: NO, Since: sepaLaunch},
			// Non-EEA founding members (EPC: resolved March 2006 - December 2013).
			{Code: CH, Since: sepaLaunch},
			{Code: MC, Since: sepaLaunch},
			{Code: SM, Since: sepaLaunch},
			// United Kingdom: in since launch; post-Brexit it remains in scope
			// (EPC effective date 1 February 2020).
			{Code: GB, Since: sepaLaunch},
			{Code: GI, Since: sepaLaunch}, // Gibraltar
			// British Crown Dependencies, from 1 May 2016.
			{Code: JE, Since: cal.MakeDate(2016, 5, 1)},
			{Code: GG, Since: cal.MakeDate(2016, 5, 1)},
			{Code: IM, Since: cal.MakeDate(2016, 5, 1)},
			// Andorra and Vatican City, from 1 March 2019.
			{Code: AD, Since: cal.MakeDate(2019, 3, 1)},
			{Code: VA, Since: cal.MakeDate(2019, 3, 1)},
			// Albania, Montenegro, North Macedonia and Moldova, operational
			// from 5 October 2025.
			{Code: AL, Since: cal.MakeDate(2025, 10, 5)},
			{Code: ME, Since: cal.MakeDate(2025, 10, 5)},
			{Code: MK, Since: cal.MakeDate(2025, 10, 5)},
			{Code: MD, Since: cal.MakeDate(2025, 10, 5)},
			// Serbia, operational from May 2026.
			{Code: RS, Since: cal.MakeDate(2026, 5, 1)},
		},
	},
}
