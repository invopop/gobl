package it

import (
	"github.com/invopop/gobl/tax"
)

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {

	switch obj := doc.(type) {
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Calculate will perform any regime specific calculations.
func Calculate(doc interface{}) error {
	switch obj := doc.(type) {
	case *tax.Identity:
		return normalizeTaxIdentity(obj)
	}
	return nil
}
