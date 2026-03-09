package rules_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestFaults(t *testing.T) {
	t.Run("empty faults should not return error", func(t *testing.T) {
		var faults rules.Faults
		assert.NoError(t, faults)
	})
}
