package l10n

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// List of all ISO 3166-2 country codes that we know about.
const (
	AF Code = "AF"
	AX Code = "AX"
	AL Code = "AL"
	DZ Code = "DZ"
	AS Code = "AS"
	AD Code = "AD"
	AO Code = "AO"
	AI Code = "AI"
	AQ Code = "AQ"
	AG Code = "AG"
	AR Code = "AR"
	AM Code = "AM"
	AW Code = "AW"
	AU Code = "AU"
	AT Code = "AT"
	AZ Code = "AZ"
	BS Code = "BS"
	BH Code = "BH"
	BD Code = "BD"
	BB Code = "BB"
	BY Code = "BY"
	BE Code = "BE"
	BZ Code = "BZ"
	BJ Code = "BJ"
	BM Code = "BM"
	BT Code = "BT"
	BO Code = "BO"
	BQ Code = "BQ"
	BA Code = "BA"
	BW Code = "BW"
	BV Code = "BV"
	BR Code = "BR"
	IO Code = "IO"
	BN Code = "BN"
	BG Code = "BG"
	BF Code = "BF"
	BI Code = "BI"
	CV Code = "CV"
	KH Code = "KH"
	CM Code = "CM"
	CA Code = "CA"
	KY Code = "KY"
	CF Code = "CF"
	TD Code = "TD"
	CL Code = "CL"
	CN Code = "CN"
	CX Code = "CX"
	CC Code = "CC"
	CO Code = "CO"
	KM Code = "KM"
	CG Code = "CG"
	CD Code = "CD"
	CK Code = "CK"
	CR Code = "CR"
	CI Code = "CI"
	HR Code = "HR"
	CU Code = "CU"
	CW Code = "CW"
	CY Code = "CY"
	CZ Code = "CZ"
	DK Code = "DK"
	DJ Code = "DJ"
	DM Code = "DM"
	DO Code = "DO"
	EC Code = "EC"
	EG Code = "EG"
	SV Code = "SV"
	GQ Code = "GQ"
	ER Code = "ER"
	EE Code = "EE"
	SZ Code = "SZ"
	ET Code = "ET"
	FK Code = "FK"
	FO Code = "FO"
	FJ Code = "FJ"
	FI Code = "FI"
	FR Code = "FR"
	GF Code = "GF"
	PF Code = "PF"
	TF Code = "TF"
	GA Code = "GA"
	GM Code = "GM"
	GE Code = "GE"
	DE Code = "DE"
	GH Code = "GH"
	GI Code = "GI"
	GR Code = "GR"
	GL Code = "GL"
	GD Code = "GD"
	GP Code = "GP"
	GU Code = "GU"
	GT Code = "GT"
	GG Code = "GG"
	GN Code = "GN"
	GW Code = "GW"
	GY Code = "GY"
	HT Code = "HT"
	HM Code = "HM"
	VA Code = "VA"
	HN Code = "HN"
	HK Code = "HK"
	HU Code = "HU"
	IS Code = "IS"
	IN Code = "IN"
	ID Code = "ID"
	IR Code = "IR"
	IQ Code = "IQ"
	IE Code = "IE"
	IM Code = "IM"
	IL Code = "IL"
	IT Code = "IT"
	JM Code = "JM"
	JP Code = "JP"
	JE Code = "JE"
	JO Code = "JO"
	KZ Code = "KZ"
	KE Code = "KE"
	KI Code = "KI"
	KP Code = "KP"
	KR Code = "KR"
	KW Code = "KW"
	KG Code = "KG"
	LA Code = "LA"
	LV Code = "LV"
	LB Code = "LB"
	LS Code = "LS"
	LR Code = "LR"
	LY Code = "LY"
	LI Code = "LI"
	LT Code = "LT"
	LU Code = "LU"
	MO Code = "MO"
	MG Code = "MG"
	MW Code = "MW"
	MY Code = "MY"
	MV Code = "MV"
	ML Code = "ML"
	MT Code = "MT"
	MH Code = "MH"
	MQ Code = "MQ"
	MR Code = "MR"
	MU Code = "MU"
	YT Code = "YT"
	MX Code = "MX"
	FM Code = "FM"
	MD Code = "MD"
	MC Code = "MC"
	MN Code = "MN"
	ME Code = "ME"
	MS Code = "MS"
	MA Code = "MA"
	MZ Code = "MZ"
	MM Code = "MM"
	NA Code = "NA"
	NR Code = "NR"
	NP Code = "NP"
	NL Code = "NL"
	NC Code = "NC"
	NZ Code = "NZ"
	NI Code = "NI"
	NE Code = "NE"
	NG Code = "NG"
	NU Code = "NU"
	NF Code = "NF"
	MK Code = "MK"
	MP Code = "MP"
	NO Code = "NO"
	OM Code = "OM"
	PK Code = "PK"
	PW Code = "PW"
	PS Code = "PS"
	PA Code = "PA"
	PG Code = "PG"
	PY Code = "PY"
	PE Code = "PE"
	PH Code = "PH"
	PN Code = "PN"
	PL Code = "PL"
	PT Code = "PT"
	PR Code = "PR"
	QA Code = "QA"
	RE Code = "RE"
	RO Code = "RO"
	RU Code = "RU"
	RW Code = "RW"
	BL Code = "BL"
	SH Code = "SH"
	KN Code = "KN"
	LC Code = "LC"
	MF Code = "MF"
	PM Code = "PM"
	VC Code = "VC"
	WS Code = "WS"
	SM Code = "SM"
	ST Code = "ST"
	SA Code = "SA"
	SN Code = "SN"
	RS Code = "RS"
	SC Code = "SC"
	SL Code = "SL"
	SG Code = "SG"
	SX Code = "SX"
	SK Code = "SK"
	SI Code = "SI"
	SB Code = "SB"
	SO Code = "SO"
	ZA Code = "ZA"
	GS Code = "GS"
	SS Code = "SS"
	ES Code = "ES"
	LK Code = "LK"
	SD Code = "SD"
	SR Code = "SR"
	SJ Code = "SJ"
	SE Code = "SE"
	CH Code = "CH"
	SY Code = "SY"
	TW Code = "TW"
	TJ Code = "TJ"
	TZ Code = "TZ"
	TH Code = "TH"
	TL Code = "TL"
	TG Code = "TG"
	TK Code = "TK"
	TO Code = "TO"
	TT Code = "TT"
	TN Code = "TN"
	TR Code = "TR"
	TM Code = "TM"
	TC Code = "TC"
	TV Code = "TV"
	UG Code = "UG"
	UA Code = "UA"
	AE Code = "AE"
	GB Code = "GB" // Great Britain and Nothern Ireland
	US Code = "US" // United States
	UM Code = "UM"
	UY Code = "UY"
	UZ Code = "UZ"
	VU Code = "VU"
	VE Code = "VE"
	VN Code = "VN"
	VG Code = "VG"
	VI Code = "VI"
	WF Code = "WF"
	EH Code = "EH"
	YE Code = "YE"
	ZM Code = "ZM"
	ZW Code = "ZW"
)

var iso3166_2 = []interface{}{
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

var (
	// IsCountry validates that the code is a valid country code.
	IsCountry = validation.In(iso3166_2...)
)
