package gobl

import (
	"time"

	"cloud.google.com/go/civil"
)

// Date represents a simple date without time used most frequently
// with business documents.
type Date struct {
	civil.Date
}

// NewDate provides a new date instance.
func NewDate(year int, month time.Month, day int) Date {
	return Date{
		civil.Date{
			Year:  year,
			Month: month,
			Day:   day,
		},
	}
}
