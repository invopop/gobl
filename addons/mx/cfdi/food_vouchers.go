package cfdi

import (
	"regexp"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/validation"
)

// Constants for the precision of the complement's amounts
const (
	FoodVouchersFinalPrecision = 2
)

// Complement's Codes Patterns
const (
	CURPPattern           = "^[A-Z][A,E,I,O,U,X][A-Z]{2}[0-9]{2}[0-1][0-9][0-3][0-9][M,H][A-Z]{2}[B,C,D,F,G,H,J,K,L,M,N,Ñ,P,Q,R,S,T,V,W,X,Y,Z]{3}[0-9,A-Z][0-9]$"
	SocialSecurityPattern = "^[0-9]{11}$"
	PostCodePattern       = "^[0-9]{5}$"
)

// Complement's Codes Regexps
var (
	CURPRegexp           = regexp.MustCompile(CURPPattern)
	SocialSecurityRegexp = regexp.MustCompile(SocialSecurityPattern)
	PostCodeRegexp       = regexp.MustCompile(PostCodePattern)
)

// FoodVouchers carries the data to produce a CFDI's "Complemento de
// Vales de Despensa" (version 1.0) providing detailed information about food
// vouchers issued by an e-wallet supplier to its customer's employees.
//
// This struct maps to the `ValesDeDespensa` root node in the CFDI's complement.
type FoodVouchers struct {
	// Customer's employer registration number (maps to `registroPatronal`).
	EmployerRegistration string `json:"employer_registration,omitempty" jsonschema:"title=Employer Registration"`
	// Customer's account number (maps to `numeroDeCuenta`).
	AccountNumber string `json:"account_number" jsonschema:"title=Account Number"`
	// Sum of all line amounts (calculated, maps to `total`).
	Total num.Amount `json:"total" jsonschema:"title=Total" jsonschema_extras:"calculated=true"`
	// List of food vouchers issued to the customer's employees (maps to `Conceptos`).
	Lines []*FoodVouchersLine `json:"lines" jsonschema:"title=Lines"`
}

// FoodVouchersLine represents a single food voucher issued to the e-wallet of
// one of the customer's employees. It maps to one `Concepto` node in the CFDI's
// complement.
type FoodVouchersLine struct {
	// Line number starting from 1 (calculated).
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Identifier of the e-wallet that received the food voucher (maps to `Identificador`).
	EWalletID cbc.Code `json:"e_wallet_id" jsonschema:"title=E-wallet Identifier"`
	// Date and time of the food voucher's issue (maps to `Fecha`).
	IssueDateTime cal.DateTime `json:"issue_date_time" jsonschema:"title=Issue Date and Time"`
	// Employee that received the food voucher.
	Employee *FoodVouchersEmployee `json:"employee,omitempty" jsonschema:"title=Employee"`
	// Amount of the food voucher (maps to `importe`).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// FoodVouchersEmployee represents an employee that received a food voucher. It
// groups employee related field that appears under the `Concepto` node in the
// CFDI's complement.
type FoodVouchersEmployee struct {
	// Employee's tax identity code (maps to `rfc`).
	TaxCode cbc.Code `json:"tax_code" jsonschema:"title=Employee's Tax Identity Code"`
	// Employee's CURP ("Clave Única de Registro de Población", maps to `curp`).
	CURP cbc.Code `json:"curp" jsonschema:"title=Employee's CURP"`
	// Employee's name (maps to `nombre`).
	Name string `json:"name" jsonschema:"title=Employee's Name"`
	// Employee's Social Security Number (maps to `numSeguridadSocial`).
	SocialSecurity cbc.Code `json:"social_security,omitempty" jsonschema:"title=Employee's Social Security Number"`
}

// Validate checks the FoodVouchers data according to the SAT's
// rules for the "Complemento de Vales de Despensa".
func (fvc *FoodVouchers) Validate() error {
	return validation.ValidateStruct(fvc,
		validation.Field(&fvc.EmployerRegistration, validation.Length(0, 20)),
		validation.Field(&fvc.AccountNumber,
			validation.Required,
			validation.Length(0, 20),
		),
		validation.Field(&fvc.Total, validation.Required),
		validation.Field(&fvc.Lines, validation.Required),
	)
}

// Validate checks the FoodVouchersLine data is valid.
func (fvl *FoodVouchersLine) Validate() error {
	return validation.ValidateStruct(fvl,
		validation.Field(&fvl.EWalletID,
			validation.Required,
			validation.Length(0, 20),
		),
		validation.Field(&fvl.IssueDateTime, cal.DateTimeNotZero()),
		validation.Field(&fvl.Employee, validation.Required),
		validation.Field(&fvl.Amount, validation.Required),
	)
}

// Validate checks the FoodVouchersEmployee data is valid.
func (fve *FoodVouchersEmployee) Validate() error {
	return validation.ValidateStruct(fve,
		validation.Field(&fve.TaxCode,
			validation.Required,
			validation.By(mx.ValidateTaxCode),
		),
		validation.Field(&fve.CURP,
			validation.Required,
			validation.Match(CURPRegexp),
		),
		validation.Field(&fve.Name,
			validation.Required,
			validation.Length(0, 100),
		),
		validation.Field(&fve.SocialSecurity,
			validation.Match(SocialSecurityRegexp),
		),
	)
}

// Calculate performs the complement's calculations and normalisations.
func (fvc *FoodVouchers) Calculate() error {
	fvc.Total = num.MakeAmount(0, FoodVouchersFinalPrecision)

	for i, l := range fvc.Lines {
		l.Index = i + 1
		l.Amount = l.Amount.Rescale(FoodVouchersFinalPrecision)

		fvc.Total = fvc.Total.Add(l.Amount)
	}

	return nil
}
