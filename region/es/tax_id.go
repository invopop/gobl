package es

// TaxCodeType represents the types of tax code supported
// in Spain.
type TaxCodeType string

const (
	National     TaxCodeType = "0"
	Foreign      TaxCodeType = "6"
	Organisation TaxCodeType = "10"
)
