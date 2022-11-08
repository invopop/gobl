package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// Preceding allows for information to be provided about a previous invoice that this one
// will replace or subtract from. If this is used, the invoice type code will most likely need
// to be set to `corrected` or `credit-note`.
type Preceding struct {
	// Preceding document's UUID if available can be useful for tracing.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Series identification code
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Code of the previous document.
	Code string `json:"code" jsonschema:"title=Code"`
	// The issue date if the previous document.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Tax period in which the previous invoice had an effect.
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
		validation.Field(&p.Series),
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.IssueDate, cal.DateNotZero()),
		validation.Field(&p.Period),
		validation.Field(&p.Corrections, validation.Each(isValidCorrectionKey)),
		validation.Field(&p.CorrectionMethod, isValidCorrectionMethodKey),
		validation.Field(&p.Meta),
	)
}

// CorrectionKey helps identify from a set of reasons why this correction
// is happening
type CorrectionKey org.Key

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

// CorrectionKeyDef holds a definition of a correction key with it's
// description.
type CorrectionKeyDef struct {
	// Key being defined
	Key CorrectionKey `json:"key" jsonschema:"title=Key"`
	// Description of the key and how it should be used.
	Description string `json:"description" jsonschema:"title=Description"`
}

// CorrectionKeyList provides a fixed list of all the correction
// codes that are currently supported by GOBL.
var CorrectionKeyDefinitions = []CorrectionKeyDef{
	{CodeCorrectionKey, "Code has changed."},
	{SeriesCorrectionKey, "Series has changed."},
	{IssueDateCorrectionKey, "Issue date was modified."},
	{SupplierCorrectionKey, "Supplier details were changed."},
	{CustomerCorrectionKey, "Customer details were changed."},
	{SupplierNameCorrectionKey, "Supplier name was changed."},
	{CustomerNameCorrectionKey, "Customer name was changed."},
	{SupplierTaxIDCorrectionKey, "Supplier Tax ID was changed."},
	{CustomerTaxIDCorrectionKey, "Customer Tax ID was changed."},
	{SupplierAddressCorrectionKey, "Supplier address was modified."},
	{CustomerAddressCorrectionKey, "Customer address was modified."},
	{LineCorrectionKey, "Line details were corrected."},
	{PeriodCorrectionKey, "Period was changed."},
	{TypeCorrectionKey, "Type of document was corrected."},
	{LegalDetailsCorrectionKey, "Legal details were corrected."},
	{TaxRateCorrectionKey, "Tax rates were modified."},
	{TaxAmountCorrectionKey, "Tax amount was corrected."},
	{TaxBaseCorrectionKey, "Taxable base was corrected."},
	{TaxCorrectionKey, "General issue with tax calculations."},
	{TaxRetainedCorrectionKey, "Error in retained tax calculations/"},
	{RefundCorrectionKey, "Goods or materials have been returned to supplier."},
	{DiscountCorrectionKey, "New discounts or rebates added."},
	{JudicialCorrectionKey, "Court ruling or administrative decision."},
	{InsolvencyCorrectionKey, "The customer is insolvent and cannot pay."},
}

var isValidCorrectionKey = validation.In(validCorrectionKeys()...)

func validCorrectionKeys() []interface{} {
	list := make([]interface{}, len(CorrectionKeyDefinitions))
	for i, v := range CorrectionKeyDefinitions {
		list[i] = v.Key
	}
	return list
}

// Defined list of correction methods
const (
	CompleteCorrectionMethodKey   CorrectionMethodKey = "complete"   // everything has changed
	PartialCorrectionMethodKey    CorrectionMethodKey = "partial"    // only differences corrected
	DiscountCorrectionMethodKey   CorrectionMethodKey = "discount"   // deducted from future invoices
	AuthorizedCorrectionMethodKey CorrectionMethodKey = "authorized" // Permitted by tax agency
)

// CorrectionMethodKeyDef defines the fields used to describe each correction method.
type CorrectionMethodKeyDef struct {
	// Key being defined
	Key CorrectionMethodKey `json:"key" jsonschema:"title=Key"`
	// Description of the key and how it should be used.
	Description string `json:"description" jsonschema:"title=Description"`
}

// CorrectionMethodKeyList provides a fixed list of codes for validation
// purposes.
var CorrectionMethodKeyDefinitions = []CorrectionMethodKeyDef{
	{CompleteCorrectionMethodKey, "Everything has changed, this document replaces the previous one."},
	{PartialCorrectionMethodKey, "Only differences corrected."},
	{DiscountCorrectionMethodKey, "Deducted from future invoices."},
	{AuthorizedCorrectionMethodKey, "Permitted by tax agency."},
}

var isValidCorrectionMethodKey = validation.In(validCorrectionMethodKeys()...)

func validCorrectionMethodKeys() []interface{} {
	list := make([]interface{}, len(CorrectionMethodKeyDefinitions))
	for i, v := range CorrectionMethodKeyDefinitions {
		list[i] = v.Key
	}
	return list
}

// JSONSchemaExtend provides additional details to the schema.
func (CorrectionKey) JSONSchemaExtend(s *jsonschema.Schema) {
	s.AnyOf = make([]*jsonschema.Schema, len(CorrectionKeyDefinitions))
	for i, v := range CorrectionKeyDefinitions {
		s.AnyOf[i] = &jsonschema.Schema{
			Const:       v.Key,
			Description: v.Description,
		}
	}
}

// JSONSchemaExtend provides additional details to the schema.
func (CorrectionMethodKey) JSONSchemaExtend(s *jsonschema.Schema) {
	s.AnyOf = make([]*jsonschema.Schema, len(CorrectionMethodKeyDefinitions))
	for i, v := range CorrectionMethodKeyDefinitions {
		s.AnyOf[i] = &jsonschema.Schema{
			Const:       v.Key,
			Description: v.Description,
		}
	}
}
