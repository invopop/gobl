package jp

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// References:
// - Corporate Number: https://www.houjin-bangou.nta.go.jp/en/setsumei/
// - Qualified Invoice: https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu1.htm

// Identity type codes for Japan.
const (
	// IdentityTypeCorporateNumber represents the Japanese Corporate Number (法人番号)
	// A 13-digit unique identifier assigned to corporations by the National Tax Agency.
	IdentityTypeCorporateNumber cbc.Code = "CN"

	// IdentityTypeQualifiedInvoiceIssuer represents the Qualified Invoice Issuer Registration Number
	// (適格請求書発行事業者登録番号).
	// Format: T followed by 13 digits (e.g., T1234567890123)
	IdentityTypeQualifiedInvoiceIssuer cbc.Code = "QII"

	// IdentityTypeMyNumber represents the Individual Number (個人番号/マイナンバー)
	// A 12-digit identifier assigned to individuals in Japan.
	IdentityTypeMyNumber cbc.Code = "MN"

	// IdentityTypeResidentRegistration represents the Resident Registration Number (住民票コード) used for local
	// government identification.
	IdentityTypeResidentRegistration cbc.Code = "RR"
)

// ValidOrgIdentityTypes lists the org identity types validated by the regime.
var ValidOrgIdentityTypes = []cbc.Code{
	IdentityTypeCorporateNumber,
	IdentityTypeQualifiedInvoiceIssuer,
	IdentityTypeMyNumber,
}

// orgIdentityDefs provides bilingual definitions for each supported org identity type, used for documentation and
// schema generation.
var orgIdentityDefs = []*cbc.Definition{
	// Corporate Number (法人番号)
	{
		Code: IdentityTypeCorporateNumber,
		Name: i18n.String{
			i18n.EN: "Corporate Number",
			i18n.JA: "法人番号",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
					The Corporate Number (法人番号, Hōjin Bangō) is a 13-digit unique
					identification number assigned by the National Tax Agency to corporations,
					other organizations, and businesses registered in Japan.

					All companies registered in Japan receive a Corporate Number, which is used
					for tax filing, social security procedures, and other administrative purposes.

					The Corporate Number can be verified at the National Tax Agency's
					Corporate Number Publication Site: https://www.houjin-bangou.nta.go.jp/en/
				`),
			i18n.JA: here.Doc(`
					法人番号は、国税庁が日本に登録されている法人その他の団体や事業者に割り当てる
					13桁の固有識別番号です。

					日本に登録されているすべての企業は法人番号を受け取り、
					税務申告、社会保険手続き、その他の行政目的に使用されます。

					法人番号は国税庁の法人番号公表サイトで確認できます：
					https://www.houjin-bangou.nta.go.jp/
				`),
		},
	},

	// Qualified Invoice Issuer Registration Number (適格請求書発行事業者登録番号)
	{
		Code: IdentityTypeQualifiedInvoiceIssuer,
		Name: i18n.String{
			i18n.EN: "Qualified Invoice Issuer Registration Number",
			i18n.JA: "適格請求書発行事業者登録番号",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
					The Qualified Invoice Issuer Registration Number is issued to businesses
					registered under the Qualified Invoice System introduced on October 1, 2023.

					Format: The letter "T" followed by a 13-digit Corporate Number.
					Example: T1234567890123

					This number must be included on qualified invoices (適格請求書) for
					recipients to claim input tax credits. The registration number can
					be verified on the National Tax Agency's Qualified Invoice Issuer
					Publication System.
				`),
			i18n.JA: here.Doc(`
					適格請求書発行事業者登録番号は、2023年10月1日に導入された
					適格請求書保存方式に登録された事業者に発行されます。

					形式：英字「T」に続く13桁の法人番号
					例：T1234567890123

					この番号は、受け手が仕入税額控除を請求するために適格請求書に
					記載する必要があります。登録番号は国税庁の適格請求書発行事業者
					公表システムで確認できます。
				`),
		},
	},

	// My Number (個人番号)
	{
		Code: IdentityTypeMyNumber,
		Name: i18n.String{
			i18n.EN: "Individual Number (My Number)",
			i18n.JA: "個人番号（マイナンバー）",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
					The Individual Number (個人番号), commonly known as "My Number" (マイナンバー),
					is a 12-digit identification number assigned to all residents of Japan,
					including foreign residents.

					It is used for tax, social security, and disaster response administrative
					procedures. Due to privacy concerns, handling of My Number requires
					strict compliance with the Act on the Use of Numbers to Identify
					Specific Individuals in Administrative Procedures.

					Note: For business and invoicing purposes, the Corporate Number or
					Qualified Invoice Issuer Registration Number is typically used instead.
				`),
			i18n.JA: here.Doc(`
					個人番号（マイナンバー）は、外国人居住者を含む日本のすべての居住者に
					割り当てられる12桁の識別番号です。

					税務、社会保障、災害対策の行政手続きに使用されます。
					プライバシーの観点から、マイナンバーの取扱いには
					行政手続における特定の個人を識別するための番号の利用等に関する法律
					（マイナンバー法）への厳格な準拠が必要です。

					注：ビジネスおよび請求書発行の目的では、通常、法人番号または
					適格請求書発行事業者登録番号が使用されます。
				`),
		},
	},

	// Resident Registration Code (住民票コード)
	{
		Code: IdentityTypeResidentRegistration,
		Name: i18n.String{
			i18n.EN: "Resident Registration Code",
			i18n.JA: "住民票コード",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
					The Resident Registration Code (住民票コード) is an 11-digit number
					assigned to residents for local government administrative purposes.
					It appears on resident registration certificates (住民票).
				`),
			i18n.JA: here.Doc(`
					住民票コードは、地方自治体の行政目的で居住者に割り当てられる
					11桁の番号です。住民票に記載されます。
				`),
		},
	},

	// Standard identity types
	{
		Key: org.IdentityKeyPassport,
		Name: i18n.String{
			i18n.EN: "Passport",
			i18n.JA: "旅券",
		},
	},
	{
		Key: org.IdentityKeyForeign,
		Name: i18n.String{
			i18n.EN: "Foreign ID Card (Residence Card)",
			i18n.JA: "在留カード",
		},
	},
	{
		Key: org.IdentityKeyResident,
		Name: i18n.String{
			i18n.EN: "Certificate of Residence",
			i18n.JA: "住民票",
		},
	},
}

// References:
// - Corporate Number: https://www.houjin-bangou.nta.go.jp/en/setsumei/
// - Qualified Invoice: https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu1.htm

var (
	// regexpCorporateNumber matches a 13-digit Corporate Number
	// The Corporate Number is assigned by the National Tax Agency to corporations
	regexpCorporateNumber = regexp.MustCompile(`^\d{13}$`)

	// regexpQualifiedInvoiceIssuer matches a Qualified Invoice Issuer Registration Number
	// Format: T (or t) followed by 13 digits
	regexpQualifiedInvoiceIssuer = regexp.MustCompile(`^[Tt]\d{13}$`)

	// regexpMyNumber matches a 12-digit Individual Number (My Number)
	regexpMyNumber = regexp.MustCompile(`^\d{12}$`)
)

// normalizeOrgIdentity normalizes organization identity codes
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}

	code := strings.ToUpper(strings.TrimSpace(id.Code.String()))
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.JP))

	id.Code = cbc.Code(code)
}

// validateOrgIdentity validates organization identity codes
func validateOrgIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}

	switch id.Type {
	case IdentityTypeCorporateNumber:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.By(validateCorporateNumber),
				validation.Skip,
			),
		)
	case IdentityTypeQualifiedInvoiceIssuer:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.By(validateQualifiedInvoiceIssuer),
				validation.Skip,
			),
		)
	case IdentityTypeMyNumber:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.By(validateMyNumber),
				validation.Skip,
			),
		)
	}

	return nil
}

// validateCorporateNumber validates a 13-digit Corporate Number
func validateCorporateNumber(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()
	if !regexpCorporateNumber.MatchString(val) {
		return errors.New("must be exactly 13 digits")
	}

	// Validate check digit
	if err := validateCorpNumberCheckDigit(val); err != nil {
		return errors.New("invalid check digit")
	}

	return nil
}

// validateCorpNumberCheckDigit validates the check digit of a Corporate Number.
// The first digit is the check digit calculated from the remaining 12 digits.
func validateCorpNumberCheckDigit(corpNum string) error {
	if len(corpNum) != 13 {
		return errors.New("invalid length")
	}

	// NTA checksum algorithm: sum base digits (indices 1–12) with weights
	// alternating 2,1 from left (odd index → weight 2, even index → weight 1).
	// Check digit = 9 - (sum mod 9), or 0 if sum mod 9 == 0.
	var sum int
	for j := 1; j <= 12; j++ {
		d := int(corpNum[j] - '0')
		w := 1
		if j%2 == 1 {
			w = 2
		}
		sum += d * w
	}

	remainder := sum % 9
	var expectedCheck int
	if remainder == 0 {
		expectedCheck = 0
	} else {
		expectedCheck = 9 - remainder
	}

	actualCheck := int(corpNum[0] - '0')
	if actualCheck != expectedCheck {
		return errors.New("check digit mismatch")
	}

	return nil
}

// validateQualifiedInvoiceIssuer validates a Qualified Invoice Issuer Registration Number
func validateQualifiedInvoiceIssuer(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()
	if !regexpQualifiedInvoiceIssuer.MatchString(val) {
		return errors.New("invalid Qualified Invoice Issuer format, expected T + 13 digits")
	}

	// Normalize to uppercase T for check digit validation
	tNum := strings.ToUpper(val)
	if err := validateTNumberCheckDigit(tNum); err != nil {
		return errors.New("invalid check digit")
	}

	return nil
}

// validateMyNumber validates a 12-digit Individual Number (My Number)
func validateMyNumber(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()
	if !regexpMyNumber.MatchString(val) {
		return errors.New("invalid My Number format, expected 12 digits")
	}

	return nil
}
