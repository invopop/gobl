package gobl

import (
	"cloud.google.com/go/civil"
)

// Date represents a simple date without time used most frequently
// with business documents.
type Date struct {
	civil.Date
}
