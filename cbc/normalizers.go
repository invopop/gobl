package cbc

import (
	"strings"
)

// NormalizeString will attempt to clean a string by removing any potentially invalid UTF-8
// characters (replaced with ?), trimming whitespace, and removing nil characters.
func NormalizeString(in string) string {
	out := strings.ReplaceAll(in, "\u0000", "")
	out = strings.ToValidUTF8(out, "?")
	return strings.TrimSpace(out)
}
