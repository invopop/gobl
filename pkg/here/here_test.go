// Copyright (c) 2014-2022 TSUYUSATO Kitsune
// This software is released under the MIT License.
// http://opensource.org/licenses/mit-license.php
// Source: https://github.com/makenowjust/heredoc

package here_test

import (
	"bytes"
	"testing"

	"github.com/invopop/gobl/internal/here"
)

type testCase struct {
	raw, expect string
}

var tests = []testCase{
	{"", ""},
	{`
		Foo
		Bar
		`,
		"Foo\nBar\n"},
	{`Foo
		Bar`,
		"Foo\nBar"},
	{`Foo
			
		Bar
		`,
		"Foo\n\t\nBar\n"}, // Second line contains two tabs.
	{`
		Foo
			Bar
				Hoge
					`,
		"Foo\n\tBar\n\t\tHoge\n\t\t\t"},
	{`Foo Bar`, "Foo Bar"},
	{
		`
		Foo
		Bar
	`, "Foo\nBar\n"},
	{"\n\u3000zenkaku space", "\x80\x80zenkaku space"},
}

func TestDoc(t *testing.T) {
	for i, test := range tests {
		result := here.Doc(test.raw)
		if result != test.expect {
			t.Errorf("tests[%d] failed: expected=> %#v, result=> %#v", i, test.expect, result)
		}
	}
}

func TestBytes(t *testing.T) {
	for i, test := range tests {
		result := here.Bytes(test.raw)
		if !bytes.Equal(result, []byte(test.expect)) {
			t.Errorf("tests[%d] failed: expected=> %#v, result=> %#v", i, test.expect, result)
		}
	}
}

func TestDocf(t *testing.T) {
	tc := `
		int: %3d
		string: %s
	`
	i := 42
	s := "Hello"
	expect := "int:  42\nstring: Hello\n"

	result := here.Docf(tc, i, s)
	if result != expect {
		t.Errorf("test failed: expected=> %#v, result=> %#v", expect, result)
	}
}
