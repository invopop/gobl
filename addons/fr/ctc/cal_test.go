package ctc

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
)

func TestDateValidation(t *testing.T) {
	t.Run("valid dates within range", func(t *testing.T) {
		testCases := []struct {
			name string
			date *cal.Date
		}{
			{
				name: "min boundary (2000-01-01)",
				date: cal.NewDate(2000, 1, 1),
			},
			{
				name: "mid range (2024-06-15)",
				date: cal.NewDate(2024, 6, 15),
			},
			{
				name: "max boundary (2099-12-31)",
				date: cal.NewDate(2099, 12, 31),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := validateDate(tc.date)
				assert.NoError(t, err, "date %v should be valid", tc.date)
			})
		}
	})

	t.Run("dates outside valid range", func(t *testing.T) {
		testCases := []struct {
			name        string
			date        *cal.Date
			errContains string
		}{
			{
				name:        "before 2000",
				date:        cal.NewDate(1999, 12, 31),
				errContains: "too early",
			},
			{
				name:        "after 2099",
				date:        cal.NewDate(2100, 1, 1),
				errContains: "too late",
			},
			{
				name:        "year 1950",
				date:        cal.NewDate(1950, 6, 15),
				errContains: "too early",
			},
			{
				name:        "year 2150",
				date:        cal.NewDate(2150, 6, 15),
				errContains: "too late",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := validateDate(tc.date)
				assert.Error(t, err, "date %v should be invalid", tc.date)
				if err != nil {
					assert.Contains(t, err.Error(), tc.errContains, "error should indicate date is out of range")
				}
			})
		}
	})
}
