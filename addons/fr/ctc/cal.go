package ctc

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/validation"
)

func validateDate(d *cal.Date) error {
	return validation.Validate(d,
		cal.DateAfter(cal.MakeDate(2000, 1, 1)),    // >= 2000-01-01
		cal.DateBefore(cal.MakeDate(2099, 12, 31)), // <= 2099-12-31
	)
}
