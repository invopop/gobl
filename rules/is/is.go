// Package is provides common tests for using inside rule assertions. Most of these are wrappers around
// govalidator functions, but some are custom implementations for specific use cases.
// Heavily inspired by github.com/invopop/validation.
package is

import (
	"regexp"
	"unicode"

	"github.com/asaskevich/govalidator"
	"github.com/invopop/gobl/rules"
)

var (
	// EmailFormat tests that a string is a valid email address format.
	// Note that it does NOT check if the MX record exists or not.
	EmailFormat = rules.ByString("email-format", govalidator.IsEmail)
	// URL tests that a string is a valid URL.
	URL = rules.ByString("url", govalidator.IsURL)
	// RequestURL tests that a string is a valid request URL.
	RequestURL = rules.ByString("request-url", govalidator.IsRequestURL)
	// RequestURI tests that a string is a valid request URI.
	RequestURI = rules.ByString("request-uri", govalidator.IsRequestURI)
	// Alpha tests that a string contains English letters only (a-zA-Z).
	Alpha = rules.ByString("alpha", govalidator.IsAlpha)
	// Digit tests that a string contains digits only (0-9).
	Digit = rules.ByString("digit", isDigit)
	// Alphanumeric tests that a string contains English letters and digits only (a-zA-Z0-9).
	Alphanumeric = rules.ByString("alphanumeric", govalidator.IsAlphanumeric)
	// UTFLetter tests that a string contains unicode letters only.
	UTFLetter = rules.ByString("utf-letter", govalidator.IsUTFLetter)
	// UTFDigit tests that a string contains unicode decimal digits only.
	UTFDigit = rules.ByString("utf-digit", govalidator.IsUTFDigit)
	// UTFLetterNumeric tests that a string contains unicode letters and numbers only.
	UTFLetterNumeric = rules.ByString("utf-letter-numeric", govalidator.IsUTFLetterNumeric)
	// UTFNumeric tests that a string contains unicode number characters (category N) only.
	UTFNumeric = rules.ByString("utf-numeric", isUTFNumeric)
	// LowerCase tests that a string contains lower case unicode letters only.
	LowerCase = rules.ByString("lower-case", govalidator.IsLowerCase)
	// UpperCase tests that a string contains upper case unicode letters only.
	UpperCase = rules.ByString("upper-case", govalidator.IsUpperCase)
	// Hexadecimal tests that a string is a valid hexadecimal number.
	Hexadecimal = rules.ByString("hexadecimal", govalidator.IsHexadecimal)
	// HexColor tests that a string is a valid hexadecimal color code.
	HexColor = rules.ByString("hex-color", govalidator.IsHexcolor)
	// RGBColor tests that a string is a valid RGB color in the form of rgb(R, G, B).
	RGBColor = rules.ByString("rgb-color", govalidator.IsRGBcolor)
	// Int tests that a string is a valid integer number.
	Int = rules.ByString("int", govalidator.IsInt)
	// Float tests that a string is a floating point number.
	Float = rules.ByString("float", govalidator.IsFloat)
	// UUIDv3 tests that a string is a valid version 3 UUID.
	UUIDv3 = rules.ByString("uuid-v3", govalidator.IsUUIDv3)
	// UUIDv4 tests that a string is a valid version 4 UUID.
	UUIDv4 = rules.ByString("uuid-v4", govalidator.IsUUIDv4)
	// UUIDv5 tests that a string is a valid version 5 UUID.
	UUIDv5 = rules.ByString("uuid-v5", govalidator.IsUUIDv5)
	// UUID tests that a string is a valid UUID (any version).
	UUID = rules.ByString("uuid", govalidator.IsUUID)
	// CreditCard tests that a string is a valid credit card number.
	CreditCard = rules.ByString("credit-card", govalidator.IsCreditCard)
	// ISBN10 tests that a string is a valid ISBN version 10.
	ISBN10 = rules.ByString("isbn-10", govalidator.IsISBN10)
	// ISBN13 tests that a string is a valid ISBN version 13.
	ISBN13 = rules.ByString("isbn-13", govalidator.IsISBN13)
	// ISBN tests that a string is a valid ISBN (either version 10 or 13).
	ISBN = rules.ByString("isbn", isISBN)
	// JSON tests that a string is in valid JSON format.
	JSON = rules.ByString("json", govalidator.IsJSON)
	// ASCII tests that a string contains ASCII characters only.
	ASCII = rules.ByString("ascii", govalidator.IsASCII)
	// PrintableASCII tests that a string contains printable ASCII characters only.
	PrintableASCII = rules.ByString("printable-ascii", govalidator.IsPrintableASCII)
	// Multibyte tests that a string contains multibyte characters.
	Multibyte = rules.ByString("multibyte", govalidator.IsMultibyte)
	// FullWidth tests that a string contains full-width characters.
	FullWidth = rules.ByString("full-width", govalidator.IsFullWidth)
	// HalfWidth tests that a string contains half-width characters.
	HalfWidth = rules.ByString("half-width", govalidator.IsHalfWidth)
	// VariableWidth tests that a string contains both full-width and half-width characters.
	VariableWidth = rules.ByString("variable-width", govalidator.IsVariableWidth)
	// Base64 tests that a string is encoded in Base64.
	Base64 = rules.ByString("base64", govalidator.IsBase64)
	// DataURI tests that a string is a valid base64-encoded data URI.
	DataURI = rules.ByString("data-uri", govalidator.IsDataURI)
	// E164 tests that a string is a valid E.164 telephone number.
	E164 = rules.ByString("e164", isE164Number)
	// CountryCode2 tests that a string is a valid ISO 3166-1 alpha-2 country code.
	CountryCode2 = rules.ByString("country-code-2", govalidator.IsISO3166Alpha2)
	// CountryCode3 tests that a string is a valid ISO 3166-1 alpha-3 country code.
	CountryCode3 = rules.ByString("country-code-3", govalidator.IsISO3166Alpha3)
	// CurrencyCode tests that a string is a valid ISO 4217 currency code.
	CurrencyCode = rules.ByString("currency-code", govalidator.IsISO4217)
	// DialString tests that a string is a valid dial string that can be passed to Dial().
	DialString = rules.ByString("dial-string", govalidator.IsDialString)
	// MAC tests that a string is a valid MAC address.
	MAC = rules.ByString("mac", govalidator.IsMAC)
	// IP tests that a string is a valid IP address (either version 4 or 6).
	IP = rules.ByString("ip", govalidator.IsIP)
	// IPv4 tests that a string is a valid version 4 IP address.
	IPv4 = rules.ByString("ipv4", govalidator.IsIPv4)
	// IPv6 tests that a string is a valid version 6 IP address.
	IPv6 = rules.ByString("ipv6", govalidator.IsIPv6)
	// Subdomain tests that a string is a valid subdomain.
	Subdomain = rules.ByString("subdomain", isSubdomain)
	// Domain tests that a string is a valid domain name.
	Domain = rules.ByString("domain", isDomain)
	// DNSName tests that a string is a valid DNS name.
	DNSName = rules.ByString("dns-name", govalidator.IsDNSName)
	// Host tests that a string is a valid IP address or DNS name.
	Host = rules.ByString("host", govalidator.IsHost)
	// Port tests that a string is a valid port number.
	Port = rules.ByString("port", govalidator.IsPort)
	// Latitude tests that a string is a valid latitude coordinate.
	Latitude = rules.ByString("latitude", govalidator.IsLatitude)
	// Longitude tests that a string is a valid longitude coordinate.
	Longitude = rules.ByString("longitude", govalidator.IsLongitude)
	// SSN tests that a string is a valid US Social Security Number.
	SSN = rules.ByString("ssn", govalidator.IsSSN)
	// Semver tests that a string is a valid semantic version.
	Semver = rules.ByString("semver", govalidator.IsSemver)
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
