package bill

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
)

// Preceding allows for information to be provided about a previous invoice that this one
// will replace or subtract from. If this is used, the invoice type code will most likely need
// to be set to `corrected` or `credit-note`.
type Preceding struct {
	// Preceding document's UUID if available can be useful for tracing.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Identity code of the previous invoice.
	Code string `json:"code" jsonschema:"title=Code"`
	// Additional identification details
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// When the preceding invoice was issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Tax period in which the previous invoice has an effect.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Specific codes for the corrections made.
	Corrections []CorrectionKey `json:"corrections,omitempty" jsonschema:"title=Corrections"`
	// How has the previous invoice been corrected?
	CorrectionMethod CorrectionMethodKey `json:"correction_method,omitempty" jsonschema:"title=Correction Method"`
	// Additional details regarding preceding invoice
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Additional semi-structured data that may be useful in specific regions
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the preceding details look okay
func (p *Preceding) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.UUID),
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.IssueDate, cal.DateNotZero()),
		validation.Field(&p.Period),
		validation.Field(&p.Corrections),
		validation.Field(&p.CorrectionMethod),
		validation.Field(&p.Meta),
	)
}

// CorrectionKey helps identify from a set of reasons why this correction
// is happening
type CorrectionKey string

// CorrectionMethodKey identifies that type of correction being applied.
type CorrectionMethodKey string

// List of currently supported correction codes. These are determined by local needs
// and could be increased as new regions are added with local requirements.
const (
	CodeCorrectionKey            CorrectionKey = "code"       // Invoice Code
	SeriesCorrectionKey          CorrectionKey = "series"     // Invoice series number
	IssueDateCorrectionKey       CorrectionKey = "issue-date" // Issue Date
	SupplierCorrectionKey        CorrectionKey = "supplier"   // General supplier details
	CustomerCorrectionKey        CorrectionKey = "customer"   // General customer details
	SupplierNameCorrectionKey    CorrectionKey = "supplier-name"
	CustomerNameCorrectionKey    CorrectionKey = "customer-name"
	SupplierTaxIDCorrectionKey   CorrectionKey = "supplier-tax-id"
	CustomerTaxIDCorrectionKey   CorrectionKey = "customer-tax-id"
	SupplierAddressCorrectionKey CorrectionKey = "supplier-addr"
	CustomerAddressCorrectionKey CorrectionKey = "customer-addr"
	LineCorrectionKey            CorrectionKey = "line"
	PeriodCorrectionKey          CorrectionKey = "period"
	TypeCorrectionKey            CorrectionKey = "type"
	LegalDetailsCorrectionKey    CorrectionKey = "legal-details"
	TaxRateCorrectionKey         CorrectionKey = "tax-rate"
	TaxAmountCorrectionKey       CorrectionKey = "tax-amount"
	TaxBaseCorrectionKey         CorrectionKey = "tax-base"
	TaxCorrectionKey             CorrectionKey = "tax"          // General issue with tax calculations
	TaxRetainedCorrectionKey     CorrectionKey = "tax-retained" // Error in retained tax calculations
	RefundCorrectionKey          CorrectionKey = "refund"       // Goods or materials have been returned to supplier
	DiscountCorrectionKey        CorrectionKey = "discount"     // New discounts or rebates added
	JudicialCorrectionKey        CorrectionKey = "judicial"     // Court ruling or administrative decision
	InsolvencyCorrectionKey      CorrectionKey = "insolvency"   // the customer is insolvent and cannot pay
)

// CorrectionKeyList provides a fixed list of all the correction
// codes that are currently supported by GOBL.
var CorrectionKeyList = []CorrectionKey{
	CodeCorrectionKey,
	SeriesCorrectionKey,
	IssueDateCorrectionKey,
	SupplierCorrectionKey,
	CustomerCorrectionKey,
	SupplierNameCorrectionKey,
	CustomerNameCorrectionKey,
	SupplierTaxIDCorrectionKey,
	CustomerTaxIDCorrectionKey,
	SupplierAddressCorrectionKey,
	CustomerAddressCorrectionKey,
	LineCorrectionKey,
	PeriodCorrectionKey,
	TypeCorrectionKey,
	LegalDetailsCorrectionKey,
	TaxRateCorrectionKey,
	TaxAmountCorrectionKey,
	TaxBaseCorrectionKey,
	TaxCorrectionKey,
	TaxRetainedCorrectionKey,
	RefundCorrectionKey,
	DiscountCorrectionKey,
	JudicialCorrectionKey,
	InsolvencyCorrectionKey,
}

// Validate ensures the correction code is part of the accepted list
func (cc CorrectionKey) Validate() error {
	for _, code := range CorrectionKeyList {
		if code == cc {
			return nil
		}
	}
	return errors.New("invalid")
}

// Defined list of correction methods
const (
	CompleteCorrectionMethodKey   CorrectionMethodKey = "complete"   // everything has changed
	PartialCorrectionMethodKey    CorrectionMethodKey = "partial"    // only differences corrected
	DiscountCorrectionMethodKey   CorrectionMethodKey = "discount"   // deducted from future invoices
	AuthorizedCorrectionMethodKey CorrectionMethodKey = "authorized" // Permitted by tax agency
)

// CorrectionMethodKeyList provides a fixed list of codes for validation
// purposes.
var CorrectionMethodKeyList = []CorrectionMethodKey{
	CompleteCorrectionMethodKey,
	PartialCorrectionMethodKey,
	DiscountCorrectionMethodKey,
	AuthorizedCorrectionMethodKey,
}

// Validate ensures the correction code is part of the accepted list
func (cc CorrectionMethodKey) Validate() error {
	for _, code := range CorrectionMethodKeyList {
		if code == cc {
			return nil
		}
	}
	return errors.New("invalid")
}
