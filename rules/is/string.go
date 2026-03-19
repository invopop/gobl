package is

import (
	"regexp"
	"unicode"

	"github.com/asaskevich/govalidator"
)

type StringTest struct {
	desc string
	test func(string) bool
}

func StringFunc(desc string, test func(string) bool) StringTest {
	return StringTest{
		desc: desc,
		test: test,
	}
}

func (t StringTest) Check(value any) bool {
	isString, str, _, _ := StringOrBytes(value)
	if !isString {
		return false
	}
	return t.test(str)
}

func (t StringTest) String() string {
	return t.desc
}

var (
	// EmailFormat tests that a string is a valid email address format.
	// Note that it does NOT check if the MX record exists or not.
	EmailFormat = StringFunc("email-format", govalidator.IsEmail)
	// URL tests that a string is a valid URL.
	URL = StringFunc("url", govalidator.IsURL)
	// RequestURL tests that a string is a valid request URL.
	RequestURL = StringFunc("request-url", govalidator.IsRequestURL)
	// RequestURI tests that a string is a valid request URI.
	RequestURI = StringFunc("request-uri", govalidator.IsRequestURI)
	// Alpha tests that a string contains English letters only (a-zA-Z).
	Alpha = StringFunc("alpha", govalidator.IsAlpha)
	// Digit tests that a string contains digits only (0-9).
	Digit = StringFunc("digit", isDigit)
	// Alphanumeric tests that a string contains English letters and digits only (a-zA-Z0-9).
	Alphanumeric = StringFunc("alphanumeric", govalidator.IsAlphanumeric)
	// UTFLetter tests that a string contains unicode letters only.
	UTFLetter = StringFunc("utf-letter", govalidator.IsUTFLetter)
	// UTFDigit tests that a string contains unicode decimal digits only.
	UTFDigit = StringFunc("utf-digit", govalidator.IsUTFDigit)
	// UTFLetterNumeric tests that a string contains unicode letters and numbers only.
	UTFLetterNumeric = StringFunc("utf-letter-numeric", govalidator.IsUTFLetterNumeric)
	// UTFNumeric tests that a string contains unicode number characters (category N) only.
	UTFNumeric = StringFunc("utf-numeric", isUTFNumeric)
	// LowerCase tests that a string contains lower case unicode letters only.
	LowerCase = StringFunc("lower-case", govalidator.IsLowerCase)
	// UpperCase tests that a string contains upper case unicode letters only.
	UpperCase = StringFunc("upper-case", govalidator.IsUpperCase)
	// Hexadecimal tests that a string is a valid hexadecimal number.
	Hexadecimal = StringFunc("hexadecimal", govalidator.IsHexadecimal)
	// HexColor tests that a string is a valid hexadecimal color code.
	HexColor = StringFunc("hex-color", govalidator.IsHexcolor)
	// RGBColor tests that a string is a valid RGB color in the form of rgb(R, G, B).
	RGBColor = StringFunc("rgb-color", govalidator.IsRGBcolor)
	// Int tests that a string is a valid integer number.
	Int = StringFunc("int", govalidator.IsInt)
	// Float tests that a string is a floating point number.
	Float = StringFunc("float", govalidator.IsFloat)
	// UUIDv3 tests that a string is a valid version 3 UUID.
	UUIDv3 = StringFunc("uuid-v3", govalidator.IsUUIDv3)
	// UUIDv4 tests that a string is a valid version 4 UUID.
	UUIDv4 = StringFunc("uuid-v4", govalidator.IsUUIDv4)
	// UUIDv5 tests that a string is a valid version 5 UUID.
	UUIDv5 = StringFunc("uuid-v5", govalidator.IsUUIDv5)
	// UUID tests that a string is a valid UUID (any version).
	UUID = StringFunc("uuid", govalidator.IsUUID)
	// CreditCard tests that a string is a valid credit card number.
	CreditCard = StringFunc("credit-card", govalidator.IsCreditCard)
	// ISBN10 tests that a string is a valid ISBN version 10.
	ISBN10 = StringFunc("isbn-10", govalidator.IsISBN10)
	// ISBN13 tests that a string is a valid ISBN version 13.
	ISBN13 = StringFunc("isbn-13", govalidator.IsISBN13)
	// ISBN tests that a string is a valid ISBN (either version 10 or 13).
	ISBN = StringFunc("isbn", isISBN)
	// JSON tests that a string is in valid JSON format.
	JSON = StringFunc("json", govalidator.IsJSON)
	// ASCII tests that a string contains ASCII characters only.
	ASCII = StringFunc("ascii", govalidator.IsASCII)
	// PrintableASCII tests that a string contains printable ASCII characters only.
	PrintableASCII = StringFunc("printable-ascii", govalidator.IsPrintableASCII)
	// Multibyte tests that a string contains multibyte characters.
	Multibyte = StringFunc("multibyte", govalidator.IsMultibyte)
	// FullWidth tests that a string contains full-width characters.
	FullWidth = StringFunc("full-width", govalidator.IsFullWidth)
	// HalfWidth tests that a string contains half-width characters.
	HalfWidth = StringFunc("half-width", govalidator.IsHalfWidth)
	// VariableWidth tests that a string contains both full-width and half-width characters.
	VariableWidth = StringFunc("variable-width", govalidator.IsVariableWidth)
	// Base64 tests that a string is encoded in Base64.
	Base64 = StringFunc("base64", govalidator.IsBase64)
	// DataURI tests that a string is a valid base64-encoded data URI.
	DataURI = StringFunc("data-uri", govalidator.IsDataURI)
	// E164 tests that a string is a valid E.164 telephone number.
	E164 = StringFunc("e164", isE164Number)
	// CountryCode2 tests that a string is a valid ISO 3166-1 alpha-2 country code.
	CountryCode2 = StringFunc("country-code-2", govalidator.IsISO3166Alpha2)
	// CountryCode3 tests that a string is a valid ISO 3166-1 alpha-3 country code.
	CountryCode3 = StringFunc("country-code-3", govalidator.IsISO3166Alpha3)
	// CurrencyCode tests that a string is a valid ISO 4217 currency code.
	CurrencyCode = StringFunc("currency-code", govalidator.IsISO4217)
	// DialString tests that a string is a valid dial string that can be passed to Dial().
	DialString = StringFunc("dial-string", govalidator.IsDialString)
	// MAC tests that a string is a valid MAC address.
	MAC = StringFunc("mac", govalidator.IsMAC)
	// IP tests that a string is a valid IP address (either version 4 or 6).
	IP = StringFunc("ip", govalidator.IsIP)
	// IPv4 tests that a string is a valid version 4 IP address.
	IPv4 = StringFunc("ipv4", govalidator.IsIPv4)
	// IPv6 tests that a string is a valid version 6 IP address.
	IPv6 = StringFunc("ipv6", govalidator.IsIPv6)
	// Subdomain tests that a string is a valid subdomain.
	Subdomain = StringFunc("subdomain", isSubdomain)
	// Domain tests that a string is a valid domain name.
	Domain = StringFunc("domain", isDomain)
	// DNSName tests that a string is a valid DNS name.
	DNSName = StringFunc("dns-name", govalidator.IsDNSName)
	// Host tests that a string is a valid IP address or DNS name.
	Host = StringFunc("host", govalidator.IsHost)
	// Port tests that a string is a valid port number.
	Port = StringFunc("port", govalidator.IsPort)
	// Latitude tests that a string is a valid latitude coordinate.
	Latitude = StringFunc("latitude", govalidator.IsLatitude)
	// Longitude tests that a string is a valid longitude coordinate.
	Longitude = StringFunc("longitude", govalidator.IsLongitude)
	// SSN tests that a string is a valid US Social Security Number.
	SSN = StringFunc("ssn", govalidator.IsSSN)
	// Semver tests that a string is a valid semantic version.
	Semver = StringFunc("semver", govalidator.IsSemver)
)

var (
	reDigit = regexp.MustCompile(`^[0-9]+$`)
	// reSubdomain source: https://stackoverflow.com/a/7933253
	reSubdomain = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?$`)
	// reE164 source: https://stackoverflow.com/a/23299989
	reE164 = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	// reDomain source: https://stackoverflow.com/a/7933253
	// Slightly modified: removed 255 max length validation since Go regex does not
	// support lookarounds. The length check is handled separately in isDomain.
	reDomain = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-zA-Z]{1,63}|xn--[a-z0-9]{1,59})$`)
)

func isDigit(value string) bool {
	return reDigit.MatchString(value)
}

func isE164Number(value string) bool {
	return reE164.MatchString(value)
}

func isSubdomain(value string) bool {
	return reSubdomain.MatchString(value)
}

func isDomain(value string) bool {
	if len(value) > 255 {
		return false
	}
	return reDomain.MatchString(value)
}

func isUTFNumeric(value string) bool {
	for _, c := range value {
		if !unicode.IsNumber(c) {
			return false
		}
	}
	return true
}

func isISBN(value string) bool {
	return govalidator.IsISBN(value, 10) || govalidator.IsISBN(value, 13)
}
