package pa

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

/*
 * Sources of data:
 *
 *  - https://dgi.mef.gob.pa/facturaelectronica
 *  - https://dgi.mef.gob.pa/FacturaElectronica/Documentacion.html
 *  - https://www.oecd.org/content/dam/oecd/en/topics/policy-issue-focus/aeoi/panama-tin.pdf
 *
 * DV algorithm reference:
 *  - https://www.anip.gob.pa/documentos/DV_RUC.pdf
 *
 * The DV (Dígito Verificador) is a 2-digit check digit pre-assigned by the DGI.
 * In the SFEP XML schema, RUC and DV are separate fields (dRuc and dDV). In GOBL,
 * the DV is embedded as the last hyphen-separated segment of tax.Identity.Code
 * (e.g., "8-888-888-08") to be consistent with how all other regimes store check
 * digits. The SFEP addon splits the last segment when building XML.
 */

// Tax identity codes for special cases.
const (
	TaxIdentityCodeFinalConsumer cbc.Code = "CIP-000-000-0000"
)

// Tax Identity Types determined from the RUC format.
const (
	TaxIdentityTypeNatural     cbc.Key = "natural"
	TaxIdentityTypeForeigner   cbc.Key = "foreigner"
	TaxIdentityTypeNaturalized cbc.Key = "naturalized"
	TaxIdentityTypeLegal       cbc.Key = "legal"
)

// Tax Identity Patterns (all include the trailing -DV suffix except final consumer)
//
// Natural person (cédula): [Province]-[Book]-[Entry]-[DV] (e.g., 8-442-445-08)
// AV - Antes de la Vigencia (before the current ID system): [Province]AV-[Book]-[Entry]-[DV]
// PI - Población Indígena (indigenous population): [Province]PI-[Book]-[Entry]-[DV]
// Foreigner - Extranjero (E): E-[Book]-[Entry]-[DV] (e.g., E-12-342-10)
// Naturalized - Naturalizado (N): N-[Book]-[Entry]-[DV] (e.g., N-45-832-58)
// PE - Panameño en el Exterior (Panamanian abroad): PE-[Book]-[Entry]-[DV]
// Legal entity: [Seq]-[Section]-[Number]-[DV] (e.g., 2486589-1-816994-62)
// NT - Número Tributario (tax number for entities without RUC): [Province]NT-[Seq]-[Book]-[Entry]-[DV]
// Final consumer: CIP-000-000-0000 (no DV)
const (
	TaxIdentityPatternNatural       = `^\d{1,2}-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternForeigner     = `^E-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternNaturalized   = `^N-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternPE            = `^PE-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternAV            = `^\d{1,2}AV-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternPI            = `^\d{1,2}PI-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternLegal         = `^\d{3,}-\d{1,4}-\d{1,7}-\d{2}$`
	TaxIdentityPatternNT            = `^\d{1,2}NT-\d{1,4}-\d{1,4}-\d{1,6}-\d{2}$`
	TaxIdentityPatternFinalConsumer = `^CIP-000-000-0000$`
)

var (
	taxCodeNaturalRegexp       = regexp.MustCompile(TaxIdentityPatternNatural)
	taxCodeForeignerRegexp     = regexp.MustCompile(TaxIdentityPatternForeigner)
	taxCodeNaturalizedRegexp   = regexp.MustCompile(TaxIdentityPatternNaturalized)
	taxCodePERegexp            = regexp.MustCompile(TaxIdentityPatternPE)
	taxCodeAVRegexp            = regexp.MustCompile(TaxIdentityPatternAV)
	taxCodePIRegexp            = regexp.MustCompile(TaxIdentityPatternPI)
	taxCodeLegalRegexp         = regexp.MustCompile(TaxIdentityPatternLegal)
	taxCodeNTRegexp            = regexp.MustCompile(TaxIdentityPatternNT)
	taxCodeFinalConsumerRegexp = regexp.MustCompile(TaxIdentityPatternFinalConsumer)

	// Keeps hyphens since they are structural separators in Panamanian RUC codes,
	// unlike other regimes where separators are cosmetic.
	taxCodeBadCharsRegexp = regexp.MustCompile(`[^A-Z0-9-]+`)
)

// RUC type identifiers as they appear in the RUC string itself.
const (
	rucPrefixForeigner   = "E"  // Extranjero (foreigner)
	rucPrefixNaturalized = "N"  // Naturalizado (naturalized citizen)
	rucPrefixPE          = "PE" // Panameño en el Exterior (Panamanian abroad)
	rucSuffixNT          = "NT" // Número Tributario (tax number for entities without RUC)
	rucSuffixAV          = "AV" // Antes de la Vigencia (before the current ID system)
	rucSuffixPI          = "PI" // Población Indígena (indigenous population)
)

// DGI-defined constants used in the 20-digit DV algorithm representation.
const (
	dvPersonFixedSegment = "0000005"

	dvProvinceNone    = "00"
	dvTypeCedula      = "00"
	dvTypeForeigner   = "50"
	dvTypeNaturalized = "40"
	dvTypePE          = "75"
	dvTypeAV          = "15"
	dvTypePI          = "79"
	dvTypeNT          = "43"
)

var (
	errDVMismatch = errors.New("dv checksum failed")
)

// validateTaxIdentity checks that the Panamanian RUC code has a valid format
// and that the DV checksum is correct.
func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil {
		return nil
	}

	return validation.ValidateStruct(tID,
		validation.Field(&tID.Code, validation.By(validateTaxCodeFormat)),
		validation.Field(&tID.Code, validation.By(validateTaxCodeDV)),
	)
}

// normalizeTaxIdentity cleans the RUC code while preserving structural hyphens.
func normalizeTaxIdentity(tID *tax.Identity) {
	if tID == nil {
		return
	}

	c := strings.ToUpper(tID.Code.String())
	c = taxCodeBadCharsRegexp.ReplaceAllString(c, "")
	c = strings.TrimPrefix(c, string(l10n.PA))
	tID.Code = cbc.Code(c)
}

func validateTaxCodeFormat(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code.IsEmpty() {
		return nil
	}

	if typ := determineTaxCodeType(code); typ.IsEmpty() {
		return tax.ErrIdentityCodeInvalid
	}

	return nil
}

func validateTaxCodeDV(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code.IsEmpty() {
		return nil
	}

	codeStr := code.String()
	if taxCodeFinalConsumerRegexp.MatchString(codeStr) {
		return nil
	}

	// Every valid RUC format contains hyphens; reaching here without one means
	// the format validator already flagged the code, so we skip silently.
	lastHyphen := strings.LastIndex(codeStr, "-")
	if lastHyphen < 0 {
		return nil
	}

	ruc := codeStr[:lastHyphen]
	providedDV := codeStr[lastHyphen+1:]

	computed, err := calculateDV(ruc)
	if err != nil {
		return err
	}

	if computed != providedDV {
		return errDVMismatch
	}

	return nil
}

// determineTaxCodeType determines the type of RUC or returns an empty key
// if the code does not match any recognized format.
func determineTaxCodeType(code cbc.Code) cbc.Key {
	codeStr := code.String()
	switch {
	case taxCodeFinalConsumerRegexp.MatchString(codeStr):
		return TaxIdentityTypeNatural

	case taxCodeForeignerRegexp.MatchString(codeStr):
		return TaxIdentityTypeForeigner

	case taxCodeNaturalizedRegexp.MatchString(codeStr):
		return TaxIdentityTypeNaturalized

	case taxCodePERegexp.MatchString(codeStr):
		return TaxIdentityTypeForeigner

	case taxCodeAVRegexp.MatchString(codeStr):
		return TaxIdentityTypeNatural

	case taxCodePIRegexp.MatchString(codeStr):
		return TaxIdentityTypeNatural

	case taxCodeNTRegexp.MatchString(codeStr):
		return TaxIdentityTypeLegal

	case taxCodeNaturalRegexp.MatchString(codeStr):
		return TaxIdentityTypeNatural

	case taxCodeLegalRegexp.MatchString(codeStr):
		return TaxIdentityTypeLegal

	default:
		return cbc.KeyEmpty
	}
}

// calculateDV computes the 2-digit DV (Dígito Verificador) from a RUC string
// (without the DV suffix) using the official DGI algorithm.
//
// The algorithm:
//  1. Transform the RUC into a 20-digit numeric string based on its type.
//  2. For old-format legal entity RUCs, apply a cross-reference substitution.
//  3. Compute two digits via weighted mod-11 (two passes).
func calculateDV(ruc string) (string, error) {
	segments := strings.Split(ruc, "-")

	if len(segments) < 3 || len(segments) > 4 {
		return "", fmt.Errorf("invalid RUC segment count: %d", len(segments))
	}

	if len(segments) == 4 && !strings.HasSuffix(segments[0], rucSuffixNT) {
		return "", fmt.Errorf("4-segment RUC must have NT suffix in first segment")
	}

	ructb, err := buildNumericRUC(ruc, segments)
	if err != nil {
		return "", err
	}

	oldFormat := isOldFormatRUC(ructb)

	if oldFormat {
		ructb = applyOldFormatCrossRef(ructb)
	}

	dv1 := digitDV(oldFormat, ructb)
	dv2 := digitDV(oldFormat, fmt.Sprintf("%s%d", ructb, dv1))

	return fmt.Sprintf("%d%d", dv1, dv2), nil
}

// isOldFormatRUC checks whether the 20-digit numeric RUC uses the old legal entity
// format. In the 20-digit representation, positions [3:5] == "00" and [5] < '5'
// identify these codes.
func isOldFormatRUC(ructb string) bool {
	return len(ructb) >= 6 &&
		digitAt(ructb, 3) == 0 && digitAt(ructb, 4) == 0 && digitAt(ructb, 5) < 5
}

// applyOldFormatCrossRef substitutes the type code at positions [5:7] of the
// 20-digit numeric RUC using the DGI cross-reference table.
func applyOldFormatCrossRef(ructb string) string {
	key := ructb[5:7]
	replacement := dvCrossRefLookup(key)

	return ructb[:5] + replacement + ructb[7:]
}

// dvCrossRefLookup returns the cross-reference substitution for old-format RUC codes.
// Source: DGI's DV algorithm, refer to source URLs above.
func dvCrossRefLookup(key string) string {
	switch key {
	case "00":
		return "00"
	case "10", "19", "34", "43":
		return "01"
	case "11", "20", "26", "35", "44":
		return "02"
	case "12", "21", "27", "36", "45":
		return "03"
	case "13", "22", "28", "37", "46":
		return "04"
	case "14", "29", "38", "47":
		return "05"
	case "15", "30", "39", "48":
		return "06"
	case "16", "23", "31", "40", "49":
		return "07"
	case "17", "24", "32", "41":
		return "08"
	case "18", "25", "33", "42":
		return "09"
	default:
		return key
	}
}

// buildNumericRUC transforms a RUC into the 20-digit numeric representation
// required by the DV algorithm. The transformation depends on the RUC type.
func buildNumericRUC(ruc string, segments []string) (string, error) {
	switch {
	case ruc[0] == rucPrefixForeigner[0]:
		return buildPersonNumeric(dvProvinceNone, dvTypeForeigner, segments[1], segments[2]), nil

	case strings.HasSuffix(segments[0], rucSuffixNT):
		province := segments[0][:len(segments[0])-len(rucSuffixNT)]
		return buildPersonNumeric(province, dvTypeNT, segments[2], segments[3]), nil

	case strings.HasSuffix(segments[0], rucSuffixAV):
		province := segments[0][:len(segments[0])-len(rucSuffixAV)]
		return buildPersonNumeric(province, dvTypeAV, segments[1], segments[2]), nil

	case strings.HasSuffix(segments[0], rucSuffixPI):
		province := segments[0][:len(segments[0])-len(rucSuffixPI)]
		return buildPersonNumeric(province, dvTypePI, segments[1], segments[2]), nil

	case segments[0] == rucPrefixPE:
		return buildPersonNumeric(dvProvinceNone, dvTypePE, segments[1], segments[2]), nil

	case ruc[0] == rucPrefixNaturalized[0]:
		return buildPersonNumeric(dvProvinceNone, dvTypeNaturalized, segments[1], segments[2]), nil

	case len(segments[0]) <= 2:
		return buildPersonNumeric(segments[0], dvTypeCedula, segments[1], segments[2]), nil

	default:
		return buildLegalNumeric(segments[0], segments[1], segments[2]), nil
	}
}

func buildPersonNumeric(province, typeCode, book, entry string) string {
	return padLeft(book, 4) +
		dvPersonFixedSegment +
		padLeft(province, 2) +
		typeCode +
		padLeft(book, 3) +
		padLeft(entry, 5)
}

func buildLegalNumeric(seq, section, number string) string {
	return padLeft(seq, 10) +
		padLeft(section, 4) +
		padLeft(number, 6)
}

func digitAt(s string, i int) int {
	return int(s[i] - '0')
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}

	return strings.Repeat("0", width-len(s)) + s
}

// digitDV computes a single DV digit using weighted mod-11.
// Weights start at 2 for the rightmost digit and increment as it moves to the left,
// a common convention in mod-11 check digit algorithms (see also BR and NL regimes).
// Example for "456": 4*4 + 5*3 + 6*2 = 43, 43%11 = 10, 11-10 = 1.
// For old-format RUCs, weight 12 is skipped.
func digitDV(oldFormat bool, ructb string) int {
	weight := 2
	sum := 0

	for i := len(ructb) - 1; i >= 0; i-- {
		// Old-format RUCs skip weight 12 in the first pass per the DGI spec.
		if oldFormat && weight == 12 {
			weight--
		}

		sum += weight * digitAt(ructb, i)
		weight++
	}

	remainder := sum % 11
	if remainder > 1 {
		return 11 - remainder
	}

	return 0
}
