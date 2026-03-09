package rules

import (
	"github.com/asaskevich/govalidator"
	"github.com/expr-lang/expr"
)

var isHelpers = []expr.Option{
	// Email
	expr.Function(
		"isEmailFormat",
		func(params ...any) (any, error) {
			return govalidator.IsEmail(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// URL / URI
	expr.Function(
		"isURL",
		func(params ...any) (any, error) {
			return govalidator.IsURL(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isRequestURL",
		func(params ...any) (any, error) {
			return govalidator.IsRequestURL(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Character set
	expr.Function(
		"isAlpha",
		func(params ...any) (any, error) {
			return govalidator.IsAlpha(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isAlphanumeric",
		func(params ...any) (any, error) {
			return govalidator.IsAlphanumeric(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isNumeric",
		func(params ...any) (any, error) {
			return govalidator.IsNumeric(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isInt",
		func(params ...any) (any, error) {
			return govalidator.IsInt(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isFloat",
		func(params ...any) (any, error) {
			return govalidator.IsFloat(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isASCII",
		func(params ...any) (any, error) {
			return govalidator.IsASCII(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isPrintableASCII",
		func(params ...any) (any, error) {
			return govalidator.IsPrintableASCII(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Case
	expr.Function(
		"isLowerCase",
		func(params ...any) (any, error) {
			return govalidator.IsLowerCase(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isUpperCase",
		func(params ...any) (any, error) {
			return govalidator.IsUpperCase(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Hex / color
	expr.Function(
		"isHexadecimal",
		func(params ...any) (any, error) {
			return govalidator.IsHexadecimal(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isHexColor",
		func(params ...any) (any, error) {
			return govalidator.IsHexcolor(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// UUID / ULID
	expr.Function(
		"isUUID",
		func(params ...any) (any, error) {
			return govalidator.IsUUID(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isUUIDv3",
		func(params ...any) (any, error) {
			return govalidator.IsUUIDv3(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isUUIDv4",
		func(params ...any) (any, error) {
			return govalidator.IsUUIDv4(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isUUIDv5",
		func(params ...any) (any, error) {
			return govalidator.IsUUIDv5(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Encoding
	expr.Function(
		"isBase64",
		func(params ...any) (any, error) {
			return govalidator.IsBase64(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isJSON",
		func(params ...any) (any, error) {
			return govalidator.IsJSON(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Network
	expr.Function(
		"isIP",
		func(params ...any) (any, error) {
			return govalidator.IsIP(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isIPv4",
		func(params ...any) (any, error) {
			return govalidator.IsIPv4(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isIPv6",
		func(params ...any) (any, error) {
			return govalidator.IsIPv6(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isDNSName",
		func(params ...any) (any, error) {
			return govalidator.IsDNSName(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isHost",
		func(params ...any) (any, error) {
			return govalidator.IsHost(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isPort",
		func(params ...any) (any, error) {
			return govalidator.IsPort(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isMAC",
		func(params ...any) (any, error) {
			return govalidator.IsMAC(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// International codes
	expr.Function(
		"isCountryCode2",
		func(params ...any) (any, error) {
			return govalidator.IsISO3166Alpha2(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isCountryCode3",
		func(params ...any) (any, error) {
			return govalidator.IsISO3166Alpha3(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isCurrencyCode",
		func(params ...any) (any, error) {
			return govalidator.IsISO4217(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isE164",
		func(params ...any) (any, error) {
			return govalidator.IsE164(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Geographic coordinates
	expr.Function(
		"isLatitude",
		func(params ...any) (any, error) {
			return govalidator.IsLatitude(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isLongitude",
		func(params ...any) (any, error) {
			return govalidator.IsLongitude(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Versioning
	expr.Function(
		"isSemver",
		func(params ...any) (any, error) {
			return govalidator.IsSemver(params[0].(string)), nil
		},
		new(func(string) bool),
	),

	// Credit card / ISBN
	expr.Function(
		"isCreditCard",
		func(params ...any) (any, error) {
			return govalidator.IsCreditCard(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isISBN10",
		func(params ...any) (any, error) {
			return govalidator.IsISBN10(params[0].(string)), nil
		},
		new(func(string) bool),
	),
	expr.Function(
		"isISBN13",
		func(params ...any) (any, error) {
			return govalidator.IsISBN13(params[0].(string)), nil
		},
		new(func(string) bool),
	),
}
