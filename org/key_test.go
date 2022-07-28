package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	var checks = []struct {
		key org.Key
		err bool
		val string
	}{
		{key: "h"},
		{key: "test"},
		{key: "1a"},
		{key: "ack1"},
		{key: org.Key("1a").With("foo"), val: "1a+foo"},
		{key: "-a", err: true},
		{key: "a-", err: true},
		{key: "1", err: true},
		{key: "-", err: true},
		{key: "+", err: true},
		{key: "a+", err: true},
	}
	for _, check := range checks {
		t.Run(check.key.String(), func(t *testing.T) {
			err := check.key.Validate()
			if check.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if check.val != "" {
				assert.Equal(t, check.key.String(), check.val)
			}
		})
	}

}

func TestKeyIn(t *testing.T) {
	c := org.Key("standard")

	assert.True(t, c.In("pro", "reduced+eqs", "standard"))
	assert.False(t, c.In("pro", "reduced"))
}
