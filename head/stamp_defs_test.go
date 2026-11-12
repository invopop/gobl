package head_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
)

func TestStampDefs(t *testing.T) {
	sd := &cbc.Definition{
		Key:   "test-stamp-code",
		Name:  i18n.NewString("Test Code"),
		Usage: &cbc.Usage{Keys: []cbc.Key{"test-addon"}},
	}
	head.RegisterStamp(sd)

	assert.Equal(t, sd, head.StampDefFor("test-stamp-code"))
	assert.Nil(t, head.StampDefFor("test-stamp-missing"))
	assert.Contains(t, head.AllStampDefs(), sd)
	assert.Contains(t, head.FilterStampDefs([]cbc.Key{"test-addon"}, nil), sd)
	assert.NotContains(t, head.FilterStampDefs([]cbc.Key{"other-addon"}, nil), sd)
}
