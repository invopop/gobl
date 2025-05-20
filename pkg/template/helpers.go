package template

import (
	"strings"
)

// Indent takes the text, finds all matching `\n`, and adds
// *two* spaces immediately after for each of the provided counts.
// This is useful for indenting variables as blocks of text to
// be correctly presented in YAML files.
//
// Example YAML block:
//
//	rsa_key: |-
//	  {{ .Key | indent 1 }}
func Indent(count int, text string) string {
	spaces := ""
	for i := 0; i < count; i++ {
		spaces = spaces + "  "
	}
	return strings.ReplaceAll(text, "\n", "\n"+spaces)
}

// Optional is useful when outputting strings to ensure that
// nil values are outputted correctly.
func Optional(in any) string {
	if in == nil {
		return ""
	}
	if s, ok := in.(string); ok {
		return s
	}
	return ""
}
