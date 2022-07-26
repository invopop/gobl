package l10n

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/jsonschema"
)

// CountryCode defines an ISO 3166-2 country code.
type CountryCode Code

// List of all ISO 3166-2 country codes that we know about.
const (
	AF CountryCode = "AF"
	AX CountryCode = "AX"
	AL CountryCode = "AL"
	DZ CountryCode = "DZ"
	AS CountryCode = "AS"
	AD CountryCode = "AD"
	AO CountryCode = "AO"
	AI CountryCode = "AI"
	AQ CountryCode = "AQ"
	AG CountryCode = "AG"
	AR CountryCode = "AR"
	AM CountryCode = "AM"
	AW CountryCode = "AW"
	AU CountryCode = "AU"
	AT CountryCode = "AT"
	AZ CountryCode = "AZ"
	BS CountryCode = "BS"
	BH CountryCode = "BH"
	BD CountryCode = "BD"
	BB CountryCode = "BB"
	BY CountryCode = "BY"
	BE CountryCode = "BE"
	BZ CountryCode = "BZ"
	BJ CountryCode = "BJ"
	BM CountryCode = "BM"
	BT CountryCode = "BT"
	BO CountryCode = "BO"
	BQ CountryCode = "BQ"
	BA CountryCode = "BA"
	BW CountryCode = "BW"
	BV CountryCode = "BV"
	BR CountryCode = "BR"
	IO CountryCode = "IO"
	BN CountryCode = "BN"
	BG CountryCode = "BG"
	BF CountryCode = "BF"
	BI CountryCode = "BI"
	CV CountryCode = "CV"
	KH CountryCode = "KH"
	CM CountryCode = "CM"
	CA CountryCode = "CA"
	KY CountryCode = "KY"
	CF CountryCode = "CF"
	TD CountryCode = "TD"
	CL CountryCode = "CL"
	CN CountryCode = "CN"
	CX CountryCode = "CX"
	CC CountryCode = "CC"
	CO CountryCode = "CO"
	KM CountryCode = "KM"
	CG CountryCode = "CG"
	CD CountryCode = "CD"
	CK CountryCode = "CK"
	CR CountryCode = "CR"
	CI CountryCode = "CI"
	HR CountryCode = "HR"
	CU CountryCode = "CU"
	CW CountryCode = "CW"
	CY CountryCode = "CY"
	CZ CountryCode = "CZ"
	DK CountryCode = "DK"
	DJ CountryCode = "DJ"
	DM CountryCode = "DM"
	DO CountryCode = "DO"
	EC CountryCode = "EC"
	EG CountryCode = "EG"
	SV CountryCode = "SV"
	GQ CountryCode = "GQ"
	ER CountryCode = "ER"
	EE CountryCode = "EE"
	SZ CountryCode = "SZ"
	ET CountryCode = "ET"
	FK CountryCode = "FK"
	FO CountryCode = "FO"
	FJ CountryCode = "FJ"
	FI CountryCode = "FI"
	FR CountryCode = "FR"
	GF CountryCode = "GF"
	PF CountryCode = "PF"
	TF CountryCode = "TF"
	GA CountryCode = "GA"
	GM CountryCode = "GM"
	GE CountryCode = "GE"
	DE CountryCode = "DE"
	GH CountryCode = "GH"
	GI CountryCode = "GI"
	GR CountryCode = "GR"
	GL CountryCode = "GL"
	GD CountryCode = "GD"
	GP CountryCode = "GP"
	GU CountryCode = "GU"
	GT CountryCode = "GT"
	GG CountryCode = "GG"
	GN CountryCode = "GN"
	GW CountryCode = "GW"
	GY CountryCode = "GY"
	HT CountryCode = "HT"
	HM CountryCode = "HM"
	VA CountryCode = "VA"
	HN CountryCode = "HN"
	HK CountryCode = "HK"
	HU CountryCode = "HU"
	IS CountryCode = "IS"
	IN CountryCode = "IN"
	ID CountryCode = "ID"
	IR CountryCode = "IR"
	IQ CountryCode = "IQ"
	IE CountryCode = "IE"
	IM CountryCode = "IM"
	IL CountryCode = "IL"
	IT CountryCode = "IT"
	JM CountryCode = "JM"
	JP CountryCode = "JP"
	JE CountryCode = "JE"
	JO CountryCode = "JO"
	KZ CountryCode = "KZ"
	KE CountryCode = "KE"
	KI CountryCode = "KI"
	KP CountryCode = "KP"
	KR CountryCode = "KR"
	KW CountryCode = "KW"
	KG CountryCode = "KG"
	LA CountryCode = "LA"
	LV CountryCode = "LV"
	LB CountryCode = "LB"
	LS CountryCode = "LS"
	LR CountryCode = "LR"
	LY CountryCode = "LY"
	LI CountryCode = "LI"
	LT CountryCode = "LT"
	LU CountryCode = "LU"
	MO CountryCode = "MO"
	MG CountryCode = "MG"
	MW CountryCode = "MW"
	MY CountryCode = "MY"
	MV CountryCode = "MV"
	ML CountryCode = "ML"
	MT CountryCode = "MT"
	MH CountryCode = "MH"
	MQ CountryCode = "MQ"
	MR CountryCode = "MR"
	MU CountryCode = "MU"
	YT CountryCode = "YT"
	MX CountryCode = "MX"
	FM CountryCode = "FM"
	MD CountryCode = "MD"
	MC CountryCode = "MC"
	MN CountryCode = "MN"
	ME CountryCode = "ME"
	MS CountryCode = "MS"
	MA CountryCode = "MA"
	MZ CountryCode = "MZ"
	MM CountryCode = "MM"
	NA CountryCode = "NA"
	NR CountryCode = "NR"
	NP CountryCode = "NP"
	NL CountryCode = "NL"
	NC CountryCode = "NC"
	NZ CountryCode = "NZ"
	NI CountryCode = "NI"
	NE CountryCode = "NE"
	NG CountryCode = "NG"
	NU CountryCode = "NU"
	NF CountryCode = "NF"
	MK CountryCode = "MK"
	MP CountryCode = "MP"
	NO CountryCode = "NO"
	OM CountryCode = "OM"
	PK CountryCode = "PK"
	PW CountryCode = "PW"
	PS CountryCode = "PS"
	PA CountryCode = "PA"
	PG CountryCode = "PG"
	PY CountryCode = "PY"
	PE CountryCode = "PE"
	PH CountryCode = "PH"
	PN CountryCode = "PN"
	PL CountryCode = "PL"
	PT CountryCode = "PT"
	PR CountryCode = "PR"
	QA CountryCode = "QA"
	RE CountryCode = "RE"
	RO CountryCode = "RO"
	RU CountryCode = "RU"
	RW CountryCode = "RW"
	BL CountryCode = "BL"
	SH CountryCode = "SH"
	KN CountryCode = "KN"
	LC CountryCode = "LC"
	MF CountryCode = "MF"
	PM CountryCode = "PM"
	VC CountryCode = "VC"
	WS CountryCode = "WS"
	SM CountryCode = "SM"
	ST CountryCode = "ST"
	SA CountryCode = "SA"
	SN CountryCode = "SN"
	RS CountryCode = "RS"
	SC CountryCode = "SC"
	SL CountryCode = "SL"
	SG CountryCode = "SG"
	SX CountryCode = "SX"
	SK CountryCode = "SK"
	SI CountryCode = "SI"
	SB CountryCode = "SB"
	SO CountryCode = "SO"
	ZA CountryCode = "ZA"
	GS CountryCode = "GS"
	SS CountryCode = "SS"
	ES CountryCode = "ES"
	LK CountryCode = "LK"
	SD CountryCode = "SD"
	SR CountryCode = "SR"
	SJ CountryCode = "SJ"
	SE CountryCode = "SE"
	CH CountryCode = "CH"
	SY CountryCode = "SY"
	TW CountryCode = "TW"
	TJ CountryCode = "TJ"
	TZ CountryCode = "TZ"
	TH CountryCode = "TH"
	TL CountryCode = "TL"
	TG CountryCode = "TG"
	TK CountryCode = "TK"
	TO CountryCode = "TO"
	TT CountryCode = "TT"
	TN CountryCode = "TN"
	TR CountryCode = "TR"
	TM CountryCode = "TM"
	TC CountryCode = "TC"
	TV CountryCode = "TV"
	UG CountryCode = "UG"
	UA CountryCode = "UA"
	AE CountryCode = "AE"
	GB CountryCode = "GB" // Great Britain and Northern Ireland
	US CountryCode = "US" // United States
	UM CountryCode = "UM"
	UY CountryCode = "UY"
	UZ CountryCode = "UZ"
	VU CountryCode = "VU"
	VE CountryCode = "VE"
	VN CountryCode = "VN"
	VG CountryCode = "VG"
	VI CountryCode = "VI"
	WF CountryCode = "WF"
	EH CountryCode = "EH"
	YE CountryCode = "YE"
	ZM CountryCode = "ZM"
	ZW CountryCode = "ZW"
)

// CountryDef provides the structure use to define a Country Code
// definition.
type CountryDef struct {
	// ISO 3166-2 Country code
	Code CountryCode `json:"code" jsonschema:"ISO Country Code"`
	// English name of the country
	Name string `json:"name" jsonschema:"Name"`
	// Internet Top-Level-Domain
	TLD string `json:"tld" jsonschema:"Top level domain"`
}

// CountryDefinitions provides and array of country definitions including
// the official ISO country code, the name in English, and the countries
// top-level-domain name.
var CountryDefinitions = []CountryDef{
	{AF, "Afghanistan", "af"},
	{AX, "Åland Islands", "ax"},
	{AL, "Albania", "al"},
	{DZ, "Algeria", "dz"},
	{AS, "American Samoa", "as"},
	{AD, "Andorra", "ad"},
	{AO, "Angola", "ao"},
	{AI, "Anguilla", "ai"},
	{AQ, "Antarctica", "aq"},
	{AG, "Antigua and Barbuda", "ag"},
	{AR, "Argentina", "ar"},
	{AM, "Armenia", "am"},
	{AW, "Aruba", "aw"},
	{AU, "Australia ", "au"},
	{AT, "Austria", "at"},
	{AZ, "Azerbaijan", "az"},
	{BS, "Bahamas (the)", "bs"},
	{BH, "Bahrain", "bh"},
	{BD, "Bangladesh", "bd"},
	{BB, "Barbados", "bb"},
	{BY, "Belarus", "by"},
	{BE, "Belgium", "be"},
	{BZ, "Belize", "bz"},
	{BJ, "Benin", "bj"},
	{BM, "Bermuda", "bm"},
	{BT, "Bhutan", "bt"},
	{BO, "Bolivia (Plurinational State of)", "bo"},
	{BQ, "Bonaire, Sint Eustatius and Saba", "bq"},
	{BA, "Bosnia and Herzegovina", "ba"},
	{BW, "Botswana", "bw"},
	{BV, "Bouvet Island", "bv"},
	{BR, "Brazil", "br"},
	{IO, "British Indian Ocean Territory (the)", "io"},
	{BN, "Brunei Darussalam", "bn"},
	{BG, "Bulgaria", "bg"},
	{BF, "Burkina Faso", "bf"},
	{BI, "Burundi", "bi"},
	{CV, "Cabo Verde", "cv"},
	{KH, "Cambodia", "kh"},
	{CM, "Cameroon", "cm"},
	{CA, "Canada", "ca"},
	{KY, "Cayman Islands (the)", "ky"},
	{CF, "Central African Republic (the)", "cf"},
	{TD, "Chad", "td"},
	{CL, "Chile", "cl"},
	{CN, "China", "cn"},
	{CX, "Christmas Island", "cx"},
	{CC, "Cocos (Keeling) Islands (the)", "cc"},
	{CO, "Colombia", "co"},
	{KM, "Comoros (the)", "km"},
	{CG, "Congo (the Democratic Republic of the)", "cd"},
	{CD, "Congo (the)", "cd"},
	{CK, "Cook Islands (the)", "ck"},
	{CR, "Costa Rica", "cr"},
	{CI, "Côte d'Ivoire", "ci"},
	{HR, "Croatia", "hr"},
	{CU, "Cuba", "cu"},
	{CW, "Curaçao", "cw"},
	{CY, "Cyprus", "cy"},
	{CZ, "Czechia", "cz"},
	{DK, "Denmark", "dk"},
	{DJ, "Djibouti", "dj"},
	{DM, "Dominica", "dm"},
	{DO, "Dominican Republic (the)", "do"},
	{EC, "Ecuador", "ec"},
	{EG, "Egypt", "eg"},
	{SV, "El Salvador", "sv"},
	{GQ, "Equatorial Guinea", "gq"},
	{ER, "Eritrea", "er"},
	{EE, "Estonia", "ee"},
	{SZ, "Eswatini", "sz"},
	{ET, "Ethiopia", "et"},
	{FK, "Falkland Islands (the)", "fk"},
	{FO, "Faroe Islands (the)", "fo"},
	{FJ, "Fiji", "fj"},
	{FI, "Finland", "fi"},
	{FR, "France", "fr"},
	{GF, "French Guiana", "gf"},
	{PF, "French Polynesia", "pf"},
	{TF, "French Southern Territories (the) ", "tf"},
	{GA, "Gabon", "ga"},
	{GM, "Gambia (the)", "gm"},
	{GE, "Georgia", "ge"},
	{DE, "Germany", "de"},
	{GH, "Ghana", "gh"},
	{GI, "Gibraltar", "gi"},
	{GR, "Greece", "gr"},
	{GL, "Greenland", "gl"},
	{GD, "Grenada", "gd"},
	{GP, "Guadeloupe", "gp"},
	{GU, "Guam", "gu"},
	{GT, "Guatemala", "gt"},
	{GG, "Guernsey", "gg"},
	{GN, "Guinea", "gn"},
	{GW, "Guinea-Bissau", "gw"},
	{GY, "Guyana", "gy"},
	{HT, "Haiti", "ht"},
	{HM, "Heard Island and McDonald Islands", "hm"},
	{VA, "Holy See (the)", "va"},
	{HN, "Honduras", "hn"},
	{HK, "Hong Kong", "hk"},
	{HU, "Hungary", "hu"},
	{IS, "Iceland", "is"},
	{IN, "India", "in"},
	{ID, "Indonesia", "id"},
	{IR, "Iran (Islamic Republic of)", "ir"},
	{IQ, "Iraq", "iq"},
	{IE, "Ireland", "ie"},
	{IM, "Isle of Man", "im"},
	{IL, "Israel", "il"},
	{IT, "Italy", "it"},
	{JM, "Jamaica", "jm"},
	{JP, "Japan", "jp"},
	{JE, "Jersey", "je"},
	{JO, "Jordan", "jo"},
	{KZ, "Kazakhstan", "kz"},
	{KE, "Kenya", "ke"},
	{KI, "Kiribati", "ki"},
	{KP, "Korea (the Democratic People's Republic of)", "kp"},
	{KR, "Korea (the Republic of)", "kr"},
	{KW, "Kuwait", "kw"},
	{KG, "Kyrgyzstan", "kg"},
	{LA, "Lao People's Democratic Republic (the)", "la"},
	{LV, "Latvia", "lv"},
	{LB, "Lebanon", "lb"},
	{LS, "Lesotho", "ls"},
	{LR, "Liberia", "lr"},
	{LY, "Libya", "ly"},
	{LI, "Liechtenstein", "li"},
	{LT, "Lithuania", "lt"},
	{LU, "Luxembourg", "lu"},
	{MO, "Macao", "mo"},
	{MK, "North Macedonia", "mk"},
	{MG, "Madagascar", "mg"},
	{MW, "Malawi", "mw"},
	{MY, "Malaysia", "my"},
	{MV, "Maldives", "mv"},
	{ML, "Mali", "ml"},
	{MT, "Malta", "mt"},
	{MH, "Marshall Islands (the)", "mh"},
	{MQ, "Martinique", "mq"},
	{MR, "Mauritania", "mr"},
	{MU, "Mauritius", "mu"},
	{YT, "Mayotte", "yt"},
	{MX, "Mexico", "mx"},
	{FM, "Micronesia (Federated States of)", "fm"},
	{MD, "Moldova (the Republic of)", "md"},
	{MC, "Monaco", "mc"},
	{MN, "Mongolia", "mn"},
	{ME, "Montenegro", "me"},
	{MS, "Montserrat", "ms"},
	{MA, "Morocco", "ma"},
	{MZ, "Mozambique", "mz"},
	{MM, "Myanmar", "mm"},
	{NA, "Namibia", "na"},
	{NR, "Nauru", "nr"},
	{NP, "Nepal", "np"},
	{NL, "Netherlands (the)", "nl"},
	{NC, "New Caledonia", "nc"},
	{NZ, "New Zealand", "nz"},
	{NI, "Nicaragua", "ni"},
	{NE, "Niger (the)", "ne"},
	{NG, "Nigeria", "ng"},
	{NU, "Niue", "nu"},
	{NF, "Norfolk Island", "nf"},
	{MP, "Northern Mariana Islands (the)", "mp"},
	{NO, "Norway", "no"},
	{OM, "Oman", "om"},
	{PK, "Pakistan", "pk"},
	{PW, "Palau", "pw"},
	{PS, "Palestine, State of", "ps"},
	{PA, "Panama", "pa"},
	{PG, "Papua New Guinea", "pg"},
	{PY, "Paraguay", "py"},
	{PE, "Peru", "pe"},
	{PH, "Philippines (the)", "ph"},
	{PN, "Pitcairn", "pn"},
	{PL, "Poland", "pl"},
	{PT, "Portugal", "pt"},
	{PR, "Puerto Rico", "pr"},
	{QA, "Qatar", "qa"},
	{RE, "Réunion", "re"},
	{RO, "Romania", "ro"},
	{RU, "Russian Federation (the)", "ru"},
	{RW, "Rwanda", "rw"},
	{BL, "Saint Barthélemy", "bl"},
	{SH, "Saint Helena, Ascension and Tristan da Cunha", "sh"},
	{KN, "Saint Kitts and Nevis", "kn"},
	{LC, "Saint Lucia", "lc"},
	{MF, "Saint Martin (French part)", "mf"},
	{PM, "Saint Pierre and Miquelon", "pm"},
	{VC, "Saint Vincent and the Grenadines", "vc"},
	{WS, "Samoa", "ws"},
	{SM, "San Marino", "sm"},
	{ST, "Sao Tome and Principe", "st"},
	{SA, "Saudi Arabia", "sa"},
	{SN, "Senegal", "sn"},
	{RS, "Serbia", "rs"},
	{SC, "Seychelles", "sc"},
	{SL, "Sierra Leone", "sl"},
	{SG, "Singapore", "sg"},
	{SX, "Sint Maarten (Dutch part)", "sx"},
	{SK, "Slovakia", "sk"},
	{SI, "Slovenia", "si"},
	{SB, "Solomon Islands", "sb"},
	{SO, "Somalia", "so"},
	{ZA, "South Africa", "za"},
	{GS, "South Georgia and the South Sandwich Islands", "gs"},
	{SS, "South Sudan", "ss"},
	{ES, "Spain", "es"},
	{LK, "Sri Lanka", "lk"},
	{SD, "Sudan (the)", "sd"},
	{SR, "Suriname", "sr"},
	{SJ, "Svalbard and Jan Mayen", ""},
	{SE, "Sweden", "se"},
	{CH, "Switzerland", "ch"},
	{SY, "Syrian Arab Republic (the)", "sy"},
	{TW, "Taiwan (Province of China)", "tw"},
	{TJ, "Tajikistan", "tj"},
	{TZ, "Tanzania, the United Republic of", "tz"},
	{TH, "Thailand", "th"},
	{TL, "Timor-Leste ", "tl"},
	{TG, "Togo", "tg"},
	{TK, "Tokelau", "tk"},
	{TO, "Tonga", "to"},
	{TT, "Trinidad and Tobago", "tt"},
	{TN, "Tunisia", "tn"},
	{TR, "Türkiye", "tr"},
	{TM, "Turkmenistan", "tm"},
	{TC, "Turks and Caicos Islands (the)", "tc"},
	{TV, "Tuvalu", "tv"},
	{UG, "Uganda", "ug"},
	{UA, "Ukraine", "ua"},
	{AE, "United Arab Emirates (the)", "ae"},
	{GB, "United Kingdom of Great Britain and Northern Ireland (the)", "uk"},
	{US, "United States of America (the)", "us"},
	{UM, "United States Minor Outlying Islands (the)", ""},
	{UY, "Uruguay", "uy"},
	{UZ, "Uzbekistan", "uz"},
	{VU, "Vanuatu", "vu"},
	{VE, "Venezuela (Bolivarian Republic of)", "ve"},
	{VN, "Viet Nam", "vn"},
	{VG, "Virgin Islands (British)", "vg"},
	{VI, "Virgin Islands (U.S.)", "vi"},
	{WF, "Wallis and Futuna", "wf"},
	{EH, "Western Sahara", ""},
	{YE, "Yemen", "ye"},
	{ZM, "Zambia", "zm"},
	{ZW, "Zimbabwe", "zw"},
}

func validCountryCodes() []interface{} {
	list := make([]interface{}, len(CountryDefinitions))
	for i, v := range CountryDefinitions {
		list[i] = string(v.Code)
	}
	return list
}

var (
	isCountry = validation.In(validCountryCodes()...)
)

// Validate ensures the country code is inside the known and valid
// list of countries.
func (c CountryCode) Validate() error {
	return validation.Validate(string(c), isCountry)
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (CountryCode) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Country Code",
		OneOf:       make([]*jsonschema.Schema, len(CountryDefinitions)),
		Description: "",
	}
	for i, v := range CountryDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       v.Code,
			Description: v.Name,
		}
	}
	return s
}

// In returns true if the country code is contained inside the provided set
func (c CountryCode) In(set ...CountryCode) bool {
	for _, x := range set {
		if c == x {
			return true
		}
	}
	return false
}

// String provides string representation of the country code
func (c CountryCode) String() string {
	return string(c)
}

// Name provides the Country Name for the code
func (c CountryCode) Name() string {
	for _, v := range CountryDefinitions {
		if v.Code == c {
			return v.Name
		}
	}
	return ""
}
