// Copyright (c) 2014-2017 TSUYUSATO Kitsune
// This software is released under the MIT License.
// http://opensource.org/licenses/mit-license.php
// Source: https://github.com/makenowjust/heredoc

// Package here provides creation of here-documents from raw strings.
//
// Golang supports raw-string syntax.
//
//	doc := `
//		Foo
//		Bar
//	`
//
// But raw-string cannot recognize indentation. Thus such content is an indented string, equivalent to
//
//	"\n\tFoo\n\tBar\n"
//
// This doesn't look good when incorporating texts into markdown or other formats.
//
// `here` solves this problem by removing unnecessary indentation based on the minimum indentation detected:
//
//	doc := here.Doc(`
//		Foo
//		Bar
//	`)
//
// Is equivalent to
//
//	"Foo\nBar\n"
//
// Given that the backtick in Go cannot be re-used in blocks, `here` will automatically
// replace all tildes (~) with backticks (`), unless escaped.
//
//	"This is a ~code~ example with \~ tildes"
//
// Will be converted to:
//
//	"This is a `code` example with ~ tildes"
package here

import (
	"fmt"
	"strings"
	"unicode"
)

const maxInt = int(^uint(0) >> 1)

// Doc returns un-indented string as here-document.
func Doc(raw string) string {
	skipFirstLine := false
	if len(raw) > 0 && raw[0] == '\n' {
		raw = raw[1:]
	} else {
		skipFirstLine = true
	}

	lines := strings.Split(raw, "\n")

	minIndentSize := getMinIndent(lines, skipFirstLine)
	lines = removeIndentation(lines, minIndentSize, skipFirstLine)
	lines = removeEmptyTailLines(lines)
	lines = replaceTildesWithBackticks(lines)

	return strings.Join(lines, "\n")
}

// Bytes returns un-indented byte array from the provided here-document.
func Bytes(raw string) []byte {
	return []byte(Doc(raw))
}

// getMinIndent calculates the minimum indentation in lines, excluding empty lines.
func getMinIndent(lines []string, skipFirstLine bool) int {
	minIndentSize := maxInt

	for i, line := range lines {
		if i == 0 && skipFirstLine {
			continue
		}

		indentSize := 0
		for _, r := range line {
			if unicode.IsSpace(r) {
				indentSize++
			} else {
				break
			}
		}

		if len(line) == indentSize {
			if i == len(lines)-1 && indentSize < minIndentSize {
				lines[i] = ""
			}
		} else if indentSize < minIndentSize {
			minIndentSize = indentSize
		}
	}
	return minIndentSize
}

// removeIndentation removes n characters from the front of each line in lines.
// Skips first line if skipFirstLine is true, skips empty lines.
func removeIndentation(lines []string, n int, skipFirstLine bool) []string {
	for i, line := range lines {
		if i == 0 && skipFirstLine {
			continue
		}

		if len(lines[i]) >= n {
			lines[i] = line[n:]
		}
	}
	return lines
}

// removeEmptyTailLines removes empty lines from the end of the lines array.
func removeEmptyTailLines(lines []string) []string {
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			break
		}
		lines = lines[:i]
	}
	return lines
}

func replaceTildesWithBackticks(lines []string) []string {
	for i := range lines {
		lines[i] = replaceTildes(lines[i], '`')
	}
	return lines
}

func replaceTildes(input string, s rune) string {
	var sb strings.Builder
	escaped := false
	for _, ch := range input {
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '~' {
			if escaped {
				sb.WriteRune('~') // Remove escape character
			} else {
				sb.WriteRune(s) // Replace ~ with rune
			}
		} else {
			if escaped {
				sb.WriteRune('\\')
			}
			sb.WriteRune(ch)
		}
		escaped = false
	}

	return sb.String()
}

// Docf returns unindented and formatted string as here-document.
// Formatting is done as for fmt.Printf().
func Docf(raw string, args ...interface{}) string {
	return fmt.Sprintf(Doc(raw), args...)
}
