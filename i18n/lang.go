package i18n

import (
	"errors"
)

// To simplify language management GoBL does not support full localization
// and instead focusses on simple multi-language based on the ISO 639-1 set
// of two letter codes. For business documents, this is sufficient as they
// are generally issued in a given country context.

// Lang represents the two letter language code.
type Lang string

// ISO 639-1 two-letter codes source: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
const (
	AB Lang = "ab" // Abkhazian
	AA Lang = "aa" // Afar
	AF Lang = "af" // Afrikaans
	AK Lang = "ak" // Akan
	SQ Lang = "sq" // Albanian
	AM Lang = "am" // Amharic
	AR Lang = "ar" // Arabic
	AN Lang = "an" // Aragonese
	HY Lang = "hy" // Armenian
	AS Lang = "as" // Assamese
	AV Lang = "av" // Avaric
	AE Lang = "ae" // Avestan
	AY Lang = "ay" // Aymara
	AZ Lang = "az" // Azerbaijani
	BM Lang = "bm" // Bambara
	BA Lang = "ba" // Bashkir
	EU Lang = "eu" // Basque
	BE Lang = "be" // Belarusian
	BN Lang = "bn" // Bengali
	BH Lang = "bh" // Bihari Languages
	BI Lang = "bi" // Bislama
	BS Lang = "bs" // Bosnian
	BR Lang = "br" // Breton
	BG Lang = "bg" // Bulgarian
	MY Lang = "my" // Burmese
	CA Lang = "ca" // Catalan, Valencian
	CH Lang = "ch" // Chamorro
	CE Lang = "ce" // Chechen
	NY Lang = "ny" // Chichewa, Chewa, Nyanja
	ZH Lang = "zh" // Chinese
	CV Lang = "cv" // Chuvash
	KW Lang = "kw" // Cornish
	CO Lang = "co" // Corsican
	CR Lang = "cr" // Cree
	HR Lang = "hr" // Croation
	CS Lang = "cs" // Czech
	DA Lang = "da" // Danish
	DV Lang = "dv" // Divehi, Dhivei, Maldivian
	NL Lang = "nl" // Dutch, Flemish
	DZ Lang = "dz" // Dzongkha
	EN Lang = "en" // English
	EO Lang = "eo" // Esperanto
	ET Lang = "et" // Estonian
	EE Lang = "ee" // Ewe
	FO Lang = "fo" // Faroese
	FJ Lang = "fj" // Fijian
	FI Lang = "fi" // Finnish
	FR Lang = "fr" // Frence
	FF Lang = "ff" // Fulah
	GL Lang = "gl" // Galician
	KA Lang = "ka" // Georgian
	DE Lang = "de" // German
	EL Lang = "el" // Greek
	GN Lang = "gn" // Guarani
	GU Lang = "gu" // Gujarati
	HT Lang = "ht" // Haitian
	HA Lang = "ha" // Hausa
	HE Lang = "he" // Hebrew
	HZ Lang = "hz" // Herero
	HI Lang = "hi" // Hindi
	HO Lang = "ho" // Hiri Motu
	HU Lang = "hu" // Hungarian
	IA Lang = "ia" // Interlingua
	ID Lang = "id" // Indonesian
	IE Lang = "ie" // Interligue
	GA Lang = "ga" // Irish
	IG Lang = "ig" // Igbo
	IK Lang = "ik" // Inupiaq
	IO Lang = "io" // Ido
	IS Lang = "is" // Icelandic
	IT Lang = "it" // Italian
	IU Lang = "iu" // Inuktitut
	JA Lang = "ja" // Japanese
	JV Lang = "jv" // Javanese
	KL Lang = "kl" // Kalaallisut, Greenlandic
	KN Lang = "kn" // Kannada
	KR Lang = "kr" // Kanuri
	KS Lang = "ks" // Kashmiri
	KK Lang = "kk" // Kazakh
	KM Lang = "km" // Central Khmer
	KI Lang = "ki" // Kikuyu, Gikuyu
	RW Lang = "rw" // Kinyarwanda
	KY Lang = "ky" // Kirighiz, Kyrgyz
	KV Lang = "kv" // Komi
	KG Lang = "kg" // Kongo
	KO Lang = "ko" // Korean
	KU Lang = "ku" // Kurdish
	KJ Lang = "kj" // Kuanyama, Kwanyama
	LA Lang = "la" // Latin
	LB Lang = "lb" // Luxemburgish, Letzeburgesch
	LG Lang = "lg" // Ganda
	LI Lang = "li" // Limburgan
	LN Lang = "ln" // Lingala
	LO Lang = "lo" // Lao
	LT Lang = "lt" // Lithuanian
	LU Lang = "lu" // Luba-Katanga
	LV Lang = "lv" // Latvian
	GV Lang = "gv" // Manx
	MK Lang = "mk" // Macedonian
	MG Lang = "mg" // Malagasy
	MS Lang = "ms" // Malay
	ML Lang = "ml" // Malayalam
	MT Lang = "mt" // Maltese
	MI Lang = "mi" // Mãori
	MR Lang = "mr" // Marathi
	MH Lang = "mh" // Marshallese
	MN Lang = "mn" // Mongolian
	NA Lang = "na" // Nauru
	NV Lang = "nv" // Navajo
	ND Lang = "nd" // North Ndebele
	NE Lang = "ne" // Nepali
	NG Lang = "ng" // Ndonga
	NB Lang = "nb" // Norwegian Bokmål
	NN Lang = "nn" // Norwegian Nynorsk
	NO Lang = "no" // Norwegian
	II Lang = "ii" // Sichuan Yi, Nuosu
	NR Lang = "nr" // South Ndebele
	OC Lang = "oc" // Occitan
	OJ Lang = "oj" // Ojibwa
	CU Lang = "cu" // Church Slavic
	OM Lang = "om" // Oromo
	OR Lang = "or" // Oriya
	OS Lang = "os" // Ossetian, Ossetic
	PA Lang = "pa" // Punjabi, Panjabi
	PI Lang = "pi" // Pali
	FA Lang = "fa" // Persian
	PL Lang = "pl" // Polish
	PS Lang = "ps" // Pashto, Pushto
	PT Lang = "pt" // Portuguese
	QU Lang = "qu" // Quechua
	RM Lang = "rm" // Romansh
	RN Lang = "rn" // Rundi
	RO Lang = "ro" // Romanian, Moldavian, Moldovan
	RU Lang = "ru" // Russian
	SA Lang = "sa" // Sanskrit
	SC Lang = "sc" // Sardinian
	SD Lang = "sd" // Sindhi
	SE Lang = "se" // Northen Sami
	SM Lang = "sm" // Samoan
	SG Lang = "sg" // Sango
	SR Lang = "sr" // Serbian
	GD Lang = "gd" // Gaelic, Scottish Gaelic
	SN Lang = "sn" // Shona
	SI Lang = "si" // Sinhala, Singalese
	SK Lang = "sk" // Slovak
	SL Lang = "sl" // Slovenian
	SO Lang = "so" // Somali
	ST Lang = "st" // Southern Sotho
	ES Lang = "es" // Spanish, Castilian
	SU Lang = "su" // Sundanese
	SW Lang = "sw" // Swahili
	SS Lang = "ss" // Swati
	SV Lang = "sv" // Swedish
	TA Lang = "ta" // Tamil
	TE Lang = "te" // Teluga
	TG Lang = "tg" // Tajik
	TH Lang = "th" // Thai
	TI Lang = "ti" // Tigrinya
	BO Lang = "bo" // Tibetan
	TK Lang = "tk" // Turkmen
	TL Lang = "tl" // Tagalog
	TN Lang = "tn" // Tswana
	TO Lang = "to" // Tonga
	TR Lang = "tr" // Turkish
	TS Lang = "ts" // Tsonga
	TT Lang = "tt" // Tatar
	TW Lang = "tw" // Twi
	TY Lang = "ty" // Tahitian
	UG Lang = "ug" // Uighur Uyghur
	UK Lang = "uk" // Ukranian
	UR Lang = "ur" // Urdu
	UZ Lang = "uz" // Uzbek
	VE Lang = "ve" // Venda
	VI Lang = "vi" // Viatnamese
	VO Lang = "vo" // Volapük
	WA Lang = "wa" // Walloon
	CY Lang = "cy" // Welsh
	WO Lang = "wo" // Wolof
	FY Lang = "fy" // Western Frisian
	XH Lang = "xh" // Xhosa
	YI Lang = "yi" // Yiddish
	YO Lang = "yo" // Yoruba
	ZA Lang = "za" // Zhuang, Chuang
	ZU Lang = "zu" // Zulu
)

var iso639_1 = []Lang{
	AB, // Abkhazian
	AA, // Afar
	AF, // Afrikaans
	AK, // Akan
	SQ, // Albanian
	AM, // Amharic
	AR, // Arabic
	AN, // Aragonese
	HY, // Armenian
	AS, // Assamese
	AV, // Avaric
	AE, // Avestan
	AY, // Aymara
	AZ, // Azerbaijani
	BM, // Bambara
	BA, // Bashkir
	EU, // Basque
	BE, // Belarusian
	BN, // Bengali
	BH, // Bihari Languages
	BI, // Bislama
	BS, // Bosnian
	BR, // Breton
	BG, // Bulgarian
	MY, // Burmese
	CA, // Catalan, Valencian
	CH, // Chamorro
	CE, // Chechen
	NY, // Chichewa, Chewa, Nyanja
	ZH, // Chinese
	CV, // Chuvash
	KW, // Cornish
	CO, // Corsican
	CR, // Cree
	HR, // Croation
	CS, // Czech
	DA, // Danish
	DV, // Divehi, Dhivei, Maldivian
	NL, // Dutch, Flemish
	DZ, // Dzongkha
	EN, // English
	EO, // Esperanto
	ET, // Estonian
	EE, // Ewe
	FO, // Faroese
	FJ, // Fijian
	FI, // Finnish
	FR, // Frence
	FF, // Fulah
	GL, // Galician
	KA, // Georgian
	DE, // German
	EL, // Greek
	GN, // Guarani
	GU, // Gujarati
	HT, // Haitian
	HA, // Hausa
	HE, // Hebrew
	HZ, // Herero
	HI, // Hindi
	HO, // Hiri Motu
	HU, // Hungarian
	IA, // Interlingua
	ID, // Indonesian
	IE, // Interligue
	GA, // Irish
	IG, // Igbo
	IK, // Inupiaq
	IO, // Ido
	IS, // Icelandic
	IT, // Italian
	IU, // Inuktitut
	JA, // Japanese
	JV, // Javanese
	KL, // Kalaallisut, Greenlandic
	KN, // Kannada
	KR, // Kanuri
	KS, // Kashmiri
	KK, // Kazakh
	KM, // Central Khmer
	KI, // Kikuyu, Gikuyu
	RW, // Kinyarwanda
	KY, // Kirighiz, Kyrgyz
	KV, // Komi
	KG, // Kongo
	KO, // Korean
	KU, // Kurdish
	KJ, // Kuanyama, Kwanyama
	LA, // Latin
	LB, // Luxemburgish, Letzeburgesch
	LG, // Ganda
	LI, // Limburgan
	LN, // Lingala
	LO, // Lao
	LT, // Lithuanian
	LU, // Luba-Katanga
	LV, // Latvian
	GV, // Manx
	MK, // Macedonian
	MG, // Malagasy
	MS, // Malay
	ML, // Malayalam
	MT, // Maltese
	MI, // Mãori
	MR, // Marathi
	MH, // Marshallese
	MN, // Mongolian
	NA, // Nauru
	NV, // Navajo
	ND, // North Ndebele
	NE, // Nepali
	NG, // Ndonga
	NB, // Norwegian Bokmål
	NN, // Norwegian Nynorsk
	NO, // Norwegian
	II, // Sichuan Yi, Nuosu
	NR, // South Ndebele
	OC, // Occitan
	OJ, // Ojibwa
	CU, // Church Slavic
	OM, // Oromo
	OR, // Oriya
	OS, // Ossetian, Ossetic
	PA, // Punjabi, Panjabi
	PI, // Pali
	FA, // Persian
	PL, // Polish
	PS, // Pashto, Pushto
	PT, // Portuguese
	QU, // Quechua
	RM, // Romansh
	RN, // Rundi
	RO, // Romanian, Moldavian, Moldovan
	RU, // Russian
	SA, // Sanskrit
	SC, // Sardinian
	SD, // Sindhi
	SE, // Northen Sami
	SM, // Samoan
	SG, // Sango
	SR, // Serbian
	GD, // Gaelic, Scottish Gaelic
	SN, // Shona
	SI, // Sinhala, Singalese
	SK, // Slovak
	SL, // Slovenian
	SO, // Somali
	ST, // Southern Sotho
	ES, // Spanish, Castilian
	SU, // Sundanese
	SW, // Swahili
	SS, // Swati
	SV, // Swedish
	TA, // Tamil
	TE, // Teluga
	TG, // Tajik
	TH, // Thai
	TI, // Tigrinya
	BO, // Tibetan
	TK, // Turkmen
	TL, // Tagalog
	TN, // Tswana
	TO, // Tonga
	TR, // Turkish
	TS, // Tsonga
	TT, // Tatar
	TW, // Twi
	TY, // Tahitian
	UG, // Uighur Uyghur
	UK, // Ukranian
	UR, // Urdu
	UZ, // Uzbek
	VE, // Venda
	VI, // Viatnamese
	VO, // Volapük
	WA, // Walloon
	CY, // Welsh
	WO, // Wolof
	FY, // Western Frisian
	XH, // Xhosa
	YI, // Yiddish
	YO, // Yoruba
	ZA, // Zhuang, Chuang
	ZU, // Zulu
}

// Validate ensures the language code is valid according
// to the ISO 639-1 two-letter list.
func (l Lang) Validate() error {
	if string(l) == "" {
		return nil
	}
	for _, lc := range iso639_1 {
		if l == lc {
			return nil
		}
	}
	return errors.New("invalid language code")
}
