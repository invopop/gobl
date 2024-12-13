package en16931_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/stretchr/testify/assert"
)

func TestScenarios(t *testing.T) {
	t.Run("provides list of scenarios", func(t *testing.T) {
		scenarios := en16931.Scenarios()
		assert.NotEmpty(t, scenarios)
	})
}
