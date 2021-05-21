package i18n

import "errors"

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
	"ab", // Abkhazian
	"aa", // Afar
	"af", // Afrikaans
	"ak", // Akan
	"sq", // Albanian
	"am", // Amharic
	"ar", // Arabic
	"an", // Aragonese
	"hy", // Armenian
	"as", // Assamese
	"av", // Avaric
	"ae", // Avestan
	"ay", // Aymara
	"az", // Azerbaijani
	"bm", // Bambara
	"ba", // Bashkir
	"eu", // Basque
	"be", // Belarusian
	"bn", // Bengali
	"bh", // Bihari Languages
	"bi", // Bislama
	"bs", // Bosnian
	"br", // Breton
	"bg", // Bulgarian
	"my", // Burmese
	"ca", // Catalan, Valencian
	"ch", // Chamorro
	"ce", // Chechen
	"ny", // Chichewa, Chewa, Nyanja
	"zh", // Chinese
	"cv", // Chuvash
	"kw", // Cornish
	"co", // Corsican
	"cr", // Cree
	"hr", // Croation
	"cs", // Czech
	"da", // Danish
	"dv", // Divehi, Dhivei, Maldivian
	"nl", // Dutch, Flemish
	"dz", // Dzongkha
	"en", // English
	"eo", // Esperanto
	"et", // Estonian
	"ee", // Ewe
	"fo", // Faroese
	"fj", // Fijian
	"fi", // Finnish
	"fr", // Frence
	"ff", // Fulah
	"gl", // Galician
	"ka", // Georgian
	"de", // German
	"el", // Greek
	"gn", // Guarani
	"gu", // Gujarati
	"ht", // Haitian
	"ha", // Hausa
	"he", // Hebrew
	"hz", // Herero
	"hi", // Hindi
	"ho", // Hiri Motu
	"hu", // Hungarian
	"ia", // Interlingua
	"id", // Indonesian
	"ie", // Interligue
	"ga", // Irish
	"ig", // Igbo
	"ik", // Inupiaq
	"io", // Ido
	"is", // Icelandic
	"it", // Italian
	"iu", // Inuktitut
	"ja", // Japanese
	"jv", // Javanese
	"kl", // Kalaallisut, Greenlandic
	"kn", // Kannada
	"kr", // Kanuri
	"ks", // Kashmiri
	"kk", // Kazakh
	"km", // Central Khmer
	"ki", // Kikuyu, Gikuyu
	"rw", // Kinyarwanda
	"ky", // Kirighiz, Kyrgyz
	"kv", // Komi
	"kg", // Kongo
	"ko", // Korean
	"ku", // Kurdish
	"kj", // Kuanyama, Kwanyama
	"la", // Latin
	"lb", // Luxemburgish, Letzeburgesch
	"lg", // Ganda
	"li", // Limburgan
	"ln", // Lingala
	"lo", // Lao
	"lt", // Lithuanian
	"lu", // Luba-Katanga
	"lv", // Latvian
	"gv", // Manx
	"mk", // Macedonian
	"mg", // Malagasy
	"ms", // Malay
	"ml", // Malayalam
	"mt", // Maltese
	"mi", // Mãori
	"mr", // Marathi
	"mh", // Marshallese
	"mn", // Mongolian
	"na", // Nauru
	"nv", // Navajo
	"nd", // North Ndebele
	"ne", // Nepali
	"ng", // Ndonga
	"nb", // Norwegian Bokmål
	"nn", // Norwegian Nynorsk
	"no", // Norwegian
	"ii", // Sichuan Yi, Nuosu
	"nr", // South Ndebele
	"oc", // Occitan
	"oj", // Ojibwa
	"cu", // Church Slavic
	"om", // Oromo
	"or", // Oriya
	"os", // Ossetian, Ossetic
	"pa", // Punjabi, Panjabi
	"pi", // Pali
	"fa", // Persian
	"pl", // Polish
	"ps", // Pashto, Pushto
	"pt", // Portuguese
	"qu", // Quechua
	"rm", // Romansh
	"rn", // Rundi
	"ro", // Romanian, Moldavian, Moldovan
	"ru", // Russian
	"sa", // Sanskrit
	"sc", // Sardinian
	"sd", // Sindhi
	"se", // Northen Sami
	"sm", // Samoan
	"sg", // Sango
	"sr", // Serbian
	"gd", // Gaelic, Scottish Gaelic
	"sn", // Shona
	"si", // Sinhala, Singalese
	"sk", // Slovak
	"sl", // Slovenian
	"so", // Somali
	"st", // Southern Sotho
	"es", // Spanish, Castilian
	"su", // Sundanese
	"sw", // Swahili
	"ss", // Swati
	"sv", // Swedish
	"ta", // Tamil
	"te", // Teluga
	"tg", // Tajik
	"th", // Thai
	"ti", // Tigrinya
	"bo", // Tibetan
	"tk", // Turkmen
	"tl", // Tagalog
	"tn", // Tswana
	"to", // Tonga
	"tr", // Turkish
	"ts", // Tsonga
	"tt", // Tatar
	"tw", // Twi
	"ty", // Tahitian
	"ug", // Uighur Uyghur
	"uk", // Ukranian
	"ur", // Urdu
	"uz", // Uzbek
	"ve", // Venda
	"vi", // Viatnamese
	"vo", // Volapük
	"wa", // Walloon
	"cy", // Welsh
	"wo", // Wolof
	"fy", // Western Frisian
	"xh", // Xhosa
	"yi", // Yiddish
	"yo", // Yoruba
	"za", // Zhuang, Chuang
	"zu", // Zulu
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
