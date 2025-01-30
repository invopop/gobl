// Copyright (c) 2014-2022 TSUYUSATO Kitsune
// This software is released under the MIT License.
// http://opensource.org/licenses/mit-license.php
// Source: https://github.com/makenowjust/heredoc

package here_test

import (
	"bytes"
	"testing"

	"github.com/invopop/gobl/pkg/here"
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
		"Foo\nBar"},
	{`Foo
		Bar`,
		"Foo\nBar"},
	{`Foo
			
		Bar
		`,
		"Foo\n\t\nBar"}, // Second line contains two tabs.
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
	`, "Foo\nBar"},
	{`
		This is a code block:

		~~~go
		var x = 1
		~~~
	`, "This is a code block:\n\n```go\nvar x = 1\n```"},
	{
		`With ~code~ blocks`,
		"With `code` blocks",
	},
	{
		`With \~code\~ allowed tildes`,
		"With ~code~ allowed tildes",
	},
	{`
		Code block with tildes:

		~~~
		x := ~foo~
		md := "Random test \~1"
		~~~
	`, "Code block with tildes:\n\n```\nx := `foo`\nmd := \"Random test ~1\"\n```"},
	{
		`This is a ~code~ example with \~ tildes`,
		"This is a `code` example with ~ tildes",
	},
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
	expect := "int:  42\nstring: Hello"

	result := here.Docf(tc, i, s)
	if result != expect {
		t.Errorf("test failed: expected=> %#v, result=> %#v", expect, result)
	}
}
