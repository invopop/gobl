package bill

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// Preceding allows for information to be provided about a previous invoice that this one
// will replace or subtract from. If this is used, the invoice type code will most likely need
// to be set to `corrected` or `credit-note`.
type Preceding struct {
	// Preceding document's UUID if available can be useful for tracing.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Identity code fo the previous invoice.
	Code string `json:"code" jsonschema:"title=Code"`
	// Additional identification details
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// When the preceding invoice was issued.
	IssueDate org.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Tax period in which the previous invoice has an effect.
	Period *org.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Specific codes for the corrections made.
	Corrections []CorrectionCode `json:"corrections,omitempty" jsonschema:"title=Corrections"`
	// How has the previous invoice been corrected?
	CorrectionMethod CorrectionMethodCode `json:"correction_method,omitempty" jsonschema:"title=Correction Method"`
	// Additional details regarding preceding invoice
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Additional semi-structured data that may be useful in specific regions
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the preceding details look okay
func (p *Preceding) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.IssueDate, org.DateNotZero()),
		validation.Field(&p.Period),
		validation.Field(&p.Corrections),
		validation.Field(&p.CorrectionMethod),
	)
}

// CorrectionCode helps identify from a set of reasons why this correction
// is happening
type CorrectionCode string

// CorrectionMethodCode identifies that type of correction being applied.
type CorrectionMethodCode string

// List of currently supported correction codes. These are determined by local needs
// and could be increased as new regions are added with local requirements.
const (
	CodeCorrectionCode            CorrectionCode = "code"       // Invoice Code
	SeriesCorrectionCode          CorrectionCode = "series"     // Invoice series number
	IssueDateCorrectionCode       CorrectionCode = "issue-date" // Issue Date
	SupplierCorrectionCode        CorrectionCode = "supplier"   // General supplier details
	CustomerCorrectionCode        CorrectionCode = "customer"   // General customer details
	SupplierNameCorrectionCode    CorrectionCode = "supplier-name"
	CustomerNameCorrectionCode    CorrectionCode = "customer-name"
	SupplierTaxIDCorrectionCode   CorrectionCode = "supplier-tax-id"
	CustomerTaxIDCorrectionCode   CorrectionCode = "customer-tax-id"
	SupplierAddressCorrectionCode CorrectionCode = "supplier-addr"
	CustomerAddressCorrectionCode CorrectionCode = "customer-addr"
	LineCorrectionCode            CorrectionCode = "line"
	PeriodCorrectionCode          CorrectionCode = "period"
	TypeCorrectionCode            CorrectionCode = "type"
	LegalDetailsCorrectionCode    CorrectionCode = "legal-details"
	TaxRateCorrectionCode         CorrectionCode = "tax-rate"
	TaxAmountCorrectionCode       CorrectionCode = "tax-amount"
	TaxBaseCorrectionCode         CorrectionCode = "tax-base"
	TaxCorrectionCode             CorrectionCode = "tax"          // General issue with tax calculations
	TaxRetainedCorrectionCode     CorrectionCode = "tax-retained" // Error in retained tax calculations
	RefundCorrectionCode          CorrectionCode = "refund"       // Goods or materials have been returned to supplier
	DiscountCorrectionCode        CorrectionCode = "discount"     // New discounts or rebates added
	JudicialCorrectionCode        CorrectionCode = "judicial"     // Court ruling or administrative decision
	InsolvencyCorrectionCode      CorrectionCode = "insolvency"   // the customer is insolvent and cannot pay
)

// CorrectionCodeList provides a fixed list of all the correction
// codes that are currently supported by GOBL.
var CorrectionCodeList = []CorrectionCode{
	CodeCorrectionCode,
	SeriesCorrectionCode,
	IssueDateCorrectionCode,
	SupplierCorrectionCode,
	CustomerCorrectionCode,
	SupplierNameCorrectionCode,
	CustomerNameCorrectionCode,
	SupplierTaxIDCorrectionCode,
	CustomerTaxIDCorrectionCode,
	SupplierAddressCorrectionCode,
	CustomerAddressCorrectionCode,
	LineCorrectionCode,
	PeriodCorrectionCode,
	TypeCorrectionCode,
	LegalDetailsCorrectionCode,
	TaxRateCorrectionCode,
	TaxAmountCorrectionCode,
	TaxBaseCorrectionCode,
	TaxCorrectionCode,
	TaxRetainedCorrectionCode,
	RefundCorrectionCode,
	DiscountCorrectionCode,
	JudicialCorrectionCode,
	InsolvencyCorrectionCode,
}

// Validate ensures the correction code is part of the accepted list
func (cc CorrectionCode) Validate() error {
	for _, code := range CorrectionCodeList {
		if code == cc {
			return nil
		}
	}
	return errors.New("invalid")
}

// Defined list of correction methods
const (
	CompleteCorrectionMethodCode   CorrectionMethodCode = "complete"   // everything has changed
	PartialCorrectionMethodCode    CorrectionMethodCode = "partial"    // only differences corrected
	DiscountCorrectionMethodCode   CorrectionMethodCode = "discount"   // deducted from future invoices
	AuthorizedCorrectionMethodCode CorrectionMethodCode = "authorized" // Permitted by tax agency
)

// CorrectionMethodCodeList provides a fixed list of codes for validation
// purposes.
var CorrectionMethodCodeList = []CorrectionMethodCode{
	CompleteCorrectionMethodCode,
	PartialCorrectionMethodCode,
	DiscountCorrectionMethodCode,
	AuthorizedCorrectionMethodCode,
}

// Validate ensures the correction code is part of the accepted list
func (cc CorrectionMethodCode) Validate() error {
	for _, code := range CorrectionMethodCodeList {
		if code == cc {
			return nil
		}
	}
	return errors.New("invalid")
}
