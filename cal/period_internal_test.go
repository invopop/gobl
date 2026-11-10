package cal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeriodHasBound(t *testing.T) {
	assert.True(t, periodHasBound(nil))
	assert.True(t, periodHasBound((*Period)(nil)))
	assert.True(t, periodHasBound("not a period"))

	assert.False(t, periodHasBound(&Period{}))
	assert.False(t, periodHasBound(&Period{Label: "Q3"}))
	assert.True(t, periodHasBound(&Period{Start: NewDate(2022, 1, 25)}))
	assert.True(t, periodHasBound(&Period{End: NewDate(2022, 2, 28)}))
	assert.True(t, periodHasBound(&Period{
		Start: NewDate(2022, 1, 25),
		End:   NewDate(2022, 2, 28),
	}))
}

func TestPeriodEndNotBeforeStart(t *testing.T) {
	assert.True(t, periodEndNotBeforeStart(nil))
	assert.True(t, periodEndNotBeforeStart((*Period)(nil)))
	assert.True(t, periodEndNotBeforeStart("not a period"))

	assert.True(t, periodEndNotBeforeStart(&Period{Start: NewDate(2022, 1, 25)}))
	assert.True(t, periodEndNotBeforeStart(&Period{End: NewDate(2022, 2, 28)}))
	assert.True(t, periodEndNotBeforeStart(&Period{
		Start: new(Date),
		End:   NewDate(2022, 2, 28),
	}))
	assert.True(t, periodEndNotBeforeStart(&Period{
		Start: NewDate(2022, 1, 25),
		End:   new(Date),
	}))

	assert.True(t, periodEndNotBeforeStart(&Period{
		Start: NewDate(2022, 1, 25),
		End:   NewDate(2022, 1, 25),
	}))
	assert.False(t, periodEndNotBeforeStart(&Period{
		Start: NewDate(2022, 1, 25),
		End:   NewDate(2022, 1, 20),
	}))
}
