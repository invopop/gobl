package org

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/alecthomas/jsonschema"
)

// Date represents a simple date without time used most frequently
// with business documents.
type Date struct {
	civil.Date
}

// NewDate provides a pointer to a new date instance.
func NewDate(year int, month time.Month, day int) *Date {
	d := MakeDate(year, month, day)
	return &d
}

// MakeDate provides a new date instance.
func MakeDate(year int, month time.Month, day int) Date {
	return Date{
		civil.Date{
			Year:  year,
			Month: month,
			Day:   day,
		},
	}
}

func (Date) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "string",
		Format:      "date",
		Title:       "Date",
		Description: "Civil date in simplified ISO format, like 2021-05-26",
	}
}
