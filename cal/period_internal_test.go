package cal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The rules engine never hands a nil pointer or a foreign type to a test
// function, so the defensive guards can only be exercised directly.
func TestPeriodHasBound(t *testing.T) {
	assert.True(t, periodHasBound(nil))
	assert.True(t, periodHasBound((*Period)(nil)))
	assert.True(t, periodHasBound("not a period"))

	assert.False(t, periodHasBound(&Period{}))
	assert.False(t, periodHasBound(&Period{Label: "Q3"}))
	assert.True(t, periodHasBound(&Period{Start: MakeDate(2022, 1, 25)}))
	assert.True(t, periodHasBound(&Period{End: MakeDate(2022, 2, 28)}))
	assert.True(t, periodHasBound(&Period{
		Start: MakeDate(2022, 1, 25),
		End:   MakeDate(2022, 2, 28),
	}))
}

func TestPeriodEndNotBeforeStart(t *testing.T) {
	assert.True(t, periodEndNotBeforeStart(nil))
	assert.True(t, periodEndNotBeforeStart((*Period)(nil)))
	assert.True(t, periodEndNotBeforeStart("not a period"))

	// One-sided periods have nothing to compare.
	assert.True(t, periodEndNotBeforeStart(&Period{Start: MakeDate(2022, 1, 25)}))
	assert.True(t, periodEndNotBeforeStart(&Period{End: MakeDate(2022, 2, 28)}))

	assert.True(t, periodEndNotBeforeStart(&Period{
		Start: MakeDate(2022, 1, 25),
		End:   MakeDate(2022, 1, 25),
	}))
	assert.False(t, periodEndNotBeforeStart(&Period{
		Start: MakeDate(2022, 1, 25),
		End:   MakeDate(2022, 1, 20),
	}))
}
