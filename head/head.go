// Package head defines the contents to be used in envelope
// headers.
package head

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("head"),
		Header{},
		Stamp{},
		Link{},
	)
	rules.Register(
		"head",
		rules.GOBL.Add("HEAD"),
		headerRules(),
		stampRules(),
		linkRules(),
	)
}
