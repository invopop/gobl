package l10n

import "errors"

// Country represents the ISO 3166-2 country code.
type Country string

// List of all ISO countries that we know about.
const (
	AF Country = "AF"
	AX Country = "AX"
	AL Country = "AL"
	DZ Country = "DZ"
	AS Country = "AS"
	AD Country = "AD"
	AO Country = "AO"
	AI Country = "AI"
	AQ Country = "AQ"
	AG Country = "AG"
	AR Country = "AR"
	AM Country = "AM"
	AW Country = "AW"
	AU Country = "AU"
	AT Country = "AT"
	AZ Country = "AZ"
	BS Country = "BS"
	BH Country = "BH"
	BD Country = "BD"
	BB Country = "BB"
	BY Country = "BY"
	BE Country = "BE"
	BZ Country = "BZ"
	BJ Country = "BJ"
	BM Country = "BM"
	BT Country = "BT"
	BO Country = "BO"
	BQ Country = "BQ"
	BA Country = "BA"
	BW Country = "BW"
	BV Country = "BV"
	BR Country = "BR"
	IO Country = "IO"
	BN Country = "BN"
	BG Country = "BG"
	BF Country = "BF"
	BI Country = "BI"
	CV Country = "CV"
	KH Country = "KH"
	CM Country = "CM"
	CA Country = "CA"
	KY Country = "KY"
	CF Country = "CF"
	TD Country = "TD"
	CL Country = "CL"
	CN Country = "CN"
	CX Country = "CX"
	CC Country = "CC"
	CO Country = "CO"
	KM Country = "KM"
	CG Country = "CG"
	CD Country = "CD"
	CK Country = "CK"
	CR Country = "CR"
	CI Country = "CI"
	HR Country = "HR"
	CU Country = "CU"
	CW Country = "CW"
	CY Country = "CY"
	CZ Country = "CZ"
	DK Country = "DK"
	DJ Country = "DJ"
	DM Country = "DM"
	DO Country = "DO"
	EC Country = "EC"
	EG Country = "EG"
	SV Country = "SV"
	GQ Country = "GQ"
	ER Country = "ER"
	EE Country = "EE"
	SZ Country = "SZ"
	ET Country = "ET"
	FK Country = "FK"
	FO Country = "FO"
	FJ Country = "FJ"
	FI Country = "FI"
	FR Country = "FR"
	GF Country = "GF"
	PF Country = "PF"
	TF Country = "TF"
	GA Country = "GA"
	GM Country = "GM"
	GE Country = "GE"
	DE Country = "DE"
	GH Country = "GH"
	GI Country = "GI"
	GR Country = "GR"
	GL Country = "GL"
	GD Country = "GD"
	GP Country = "GP"
	GU Country = "GU"
	GT Country = "GT"
	GG Country = "GG"
	GN Country = "GN"
	GW Country = "GW"
	GY Country = "GY"
	HT Country = "HT"
	HM Country = "HM"
	VA Country = "VA"
	HN Country = "HN"
	HK Country = "HK"
	HU Country = "HU"
	IS Country = "IS"
	IN Country = "IN"
	ID Country = "ID"
	IR Country = "IR"
	IQ Country = "IQ"
	IE Country = "IE"
	IM Country = "IM"
	IL Country = "IL"
	IT Country = "IT"
	JM Country = "JM"
	JP Country = "JP"
	JE Country = "JE"
	JO Country = "JO"
	KZ Country = "KZ"
	KE Country = "KE"
	KI Country = "KI"
	KP Country = "KP"
	KR Country = "KR"
	KW Country = "KW"
	KG Country = "KG"
	LA Country = "LA"
	LV Country = "LV"
	LB Country = "LB"
	LS Country = "LS"
	LR Country = "LR"
	LY Country = "LY"
	LI Country = "LI"
	LT Country = "LT"
	LU Country = "LU"
	MO Country = "MO"
	MG Country = "MG"
	MW Country = "MW"
	MY Country = "MY"
	MV Country = "MV"
	ML Country = "ML"
	MT Country = "MT"
	MH Country = "MH"
	MQ Country = "MQ"
	MR Country = "MR"
	MU Country = "MU"
	YT Country = "YT"
	MX Country = "MX"
	FM Country = "FM"
	MD Country = "MD"
	MC Country = "MC"
	MN Country = "MN"
	ME Country = "ME"
	MS Country = "MS"
	MA Country = "MA"
	MZ Country = "MZ"
	MM Country = "MM"
	NA Country = "NA"
	NR Country = "NR"
	NP Country = "NP"
	NL Country = "NL"
	NC Country = "NC"
	NZ Country = "NZ"
	NI Country = "NI"
	NE Country = "NE"
	NG Country = "NG"
	NU Country = "NU"
	NF Country = "NF"
	MK Country = "MK"
	MP Country = "MP"
	NO Country = "NO"
	OM Country = "OM"
	PK Country = "PK"
	PW Country = "PW"
	PS Country = "PS"
	PA Country = "PA"
	PG Country = "PG"
	PY Country = "PY"
	PE Country = "PE"
	PH Country = "PH"
	PN Country = "PN"
	PL Country = "PL"
	PT Country = "PT"
	PR Country = "PR"
	QA Country = "QA"
	RE Country = "RE"
	RO Country = "RO"
	RU Country = "RU"
	RW Country = "RW"
	BL Country = "BL"
	SH Country = "SH"
	KN Country = "KN"
	LC Country = "LC"
	MF Country = "MF"
	PM Country = "PM"
	VC Country = "VC"
	WS Country = "WS"
	SM Country = "SM"
	ST Country = "ST"
	SA Country = "SA"
	SN Country = "SN"
	RS Country = "RS"
	SC Country = "SC"
	SL Country = "SL"
	SG Country = "SG"
	SX Country = "SX"
	SK Country = "SK"
	SI Country = "SI"
	SB Country = "SB"
	SO Country = "SO"
	ZA Country = "ZA"
	GS Country = "GS"
	SS Country = "SS"
	ES Country = "ES"
	LK Country = "LK"
	SD Country = "SD"
	SR Country = "SR"
	SJ Country = "SJ"
	SE Country = "SE"
	CH Country = "CH"
	SY Country = "SY"
	TW Country = "TW"
	TJ Country = "TJ"
	TZ Country = "TZ"
	TH Country = "TH"
	TL Country = "TL"
	TG Country = "TG"
	TK Country = "TK"
	TO Country = "TO"
	TT Country = "TT"
	TN Country = "TN"
	TR Country = "TR"
	TM Country = "TM"
	TC Country = "TC"
	TV Country = "TV"
	UG Country = "UG"
	UA Country = "UA"
	AE Country = "AE"
	GB Country = "GB"
	US Country = "US"
	UM Country = "UM"
	UY Country = "UY"
	UZ Country = "UZ"
	VU Country = "VU"
	VE Country = "VE"
	VN Country = "VN"
	VG Country = "VG"
	VI Country = "VI"
	WF Country = "WF"
	EH Country = "EH"
	YE Country = "YE"
	ZM Country = "ZM"
	ZW Country = "ZW"
)

var iso3166_2 = []Country{
	AF,
	AX,
	AL,
	DZ,
	AS,
	AD,
	AO,
	AI,
	AQ,
	AG,
	AR,
	AM,
	AW,
	AU,
	AT,
	AZ,
	BS,
	BH,
	BD,
	BB,
	BY,
	BE,
	BZ,
	BJ,
	BM,
	BT,
	BO,
	BQ,
	BA,
	BW,
	BV,
	BR,
	IO,
	BN,
	BG,
	BF,
	BI,
	CV,
	KH,
	CM,
	CA,
	KY,
	CF,
	TD,
	CL,
	CN,
	CX,
	CC,
	CO,
	KM,
	CG,
	CD,
	CK,
	CR,
	CI,
	HR,
	CU,
	CW,
	CY,
	CZ,
	DK,
	DJ,
	DM,
	DO,
	EC,
	EG,
	SV,
	GQ,
	ER,
	EE,
	SZ,
	ET,
	FK,
	FO,
	FJ,
	FI,
	FR,
	GF,
	PF,
	TF,
	GA,
	GM,
	GE,
	DE,
	GH,
	GI,
	GR,
	GL,
	GD,
	GP,
	GU,
	GT,
	GG,
	GN,
	GW,
	GY,
	HT,
	HM,
	VA,
	HN,
	HK,
	HU,
	IS,
	IN,
	ID,
	IR,
	IQ,
	IE,
	IM,
	IL,
	IT,
	JM,
	JP,
	JE,
	JO,
	KZ,
	KE,
	KI,
	KP,
	KR,
	KW,
	KG,
	LA,
	LV,
	LB,
	LS,
	LR,
	LY,
	LI,
	LT,
	LU,
	MO,
	MG,
	MW,
	MY,
	MV,
	ML,
	MT,
	MH,
	MQ,
	MR,
	MU,
	YT,
	MX,
	FM,
	MD,
	MC,
	MN,
	ME,
	MS,
	MA,
	MZ,
	MM,
	NA,
	NR,
	NP,
	NL,
	NC,
	NZ,
	NI,
	NE,
	NG,
	NU,
	NF,
	MK,
	MP,
	NO,
	OM,
	PK,
	PW,
	PS,
	PA,
	PG,
	PY,
	PE,
	PH,
	PN,
	PL,
	PT,
	PR,
	QA,
	RE,
	RO,
	RU,
	RW,
	BL,
	SH,
	KN,
	LC,
	MF,
	PM,
	VC,
	WS,
	SM,
	ST,
	SA,
	SN,
	RS,
	SC,
	SL,
	SG,
	SX,
	SK,
	SI,
	SB,
	SO,
	ZA,
	GS,
	SS,
	ES,
	LK,
	SD,
	SR,
	SJ,
	SE,
	CH,
	SY,
	TW,
	TJ,
	TZ,
	TH,
	TL,
	TG,
	TK,
	TO,
	TT,
	TN,
	TR,
	TM,
	TC,
	TV,
	UG,
	UA,
	AE,
	GB,
	US,
	UM,
	UY,
	UZ,
	VU,
	VE,
	VN,
	VG,
	VI,
	WF,
	EH,
	YE,
	ZM,
	ZW,
}

// Validate ensures the country code is valid according
// to the ISO 3166-2 list.
func (c Country) Validate() error {
	if string(c) == "" {
		return nil
	}
	for _, cc := range iso3166_2 {
		if c == cc {
			return nil
		}
	}
	return errors.New("invalid country code")
}
