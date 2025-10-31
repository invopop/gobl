# 🇦🇷 GOBL Argentina Tax Regime Development Guide

This document provides comprehensive instructions for building the Argentina tax regime for GOBL. It outlines the structure, requirements, and implementation steps based on patterns observed in existing regimes.

## Overview

Argentina uses a complex tax system with multiple tax categories including VAT (IVA), income tax (Ganancias), and various provincial taxes. The implementation should follow GOBL's regime architecture while accommodating Argentina's specific tax requirements.

Based on existing regimes, you must build the Argentina regime in `/regimes/ar/`. Take as reference other regimes for other countries, to know what to build, what to do research on and the format. The common files are:

### Regime files

1. **`ar.go`** - Main regime definition file
2. **`tax_categories.go`** - Tax category definitions and rates
3. **`tax_identity.go`** - Tax ID validation and normalization
4. **`tax_identity_test.go`** - Tests for tax ID validation
5. **`README.md`** - Argentina-specific documentation (this file)

6. **`scenarios.go`** - Tax scenarios and special regimes
7. **`scenarios_test.go`** - Tests for scenarios
8. **`invoices.go`** - Invoice-specific validations
9. **`invoices_test.go`** - Tests for invoice validations
10. **`party.go`** - Party/organization validations
11. **`party_test.go`** - Tests for party validations
12. **`identities.go`** - Identity type definitions
13. **`identities_test.go`** - Tests for identity definitions
14. **`corrections.go`** - Correction/credit note definitions
15. **`notes.go`** - Special note requirements

## Research Requirements

### Official Sources

**CRITICAL**: All tax rates, validation rules, and business logic must be based on official Argentine government sources. Include references in code comments.

#### Primary Sources to Research:

1. **AFIP (Administración Federal de Ingresos Públicos)**
   - Official tax regulations
   - Current tax rates
   - Identification number formats and validation rules
   - Electronic invoicing requirements (if applicable)

2. **Provincial Tax Authorities**
   - Ingresos Brutos rates by province
   - Provincial identification requirements

3. **Official Legal Documents**
   - Tax laws and regulations
   - Administrative resolutions
   - Technical specifications for electronic systems

### Key Research Areas

1. **Tax Categories and Rates**
   - Current IVA rates (general, reduced, exempt)
   - Ganancias withholding rates
   - Provincial tax rates
   - Monotributo categories and rates

2. **Identification Numbers**
   - CUIT/CUIL format and validation algorithm
   - CDI format for foreigners
   - Validation rules and check digits

3. **Business Rules**
   - Invoice numbering requirements
   - Electronic invoicing obligations
   - Special tax regimes
   - Exemptions and special cases

4. **Address Requirements**
   - Province and municipality codes
   - Postal code formats
   - Address validation rules

## Testing Requirements

Create comprehensive tests for:

1. **Tax Identity Validation**
   - Valid CUIT/CUIL numbers
   - Invalid formats
   - Edge cases and special numbers

2. **Tax Calculations**
   - Rate calculations
   - Rounding rules
   - Date-based rate changes

3. **Business Rules**
   - Invoice validation
   - Special regime handling
   - Exemption scenarios

## Documentation Requirements

### README.md Structure

The Argentina README should include:

1. **Overview** of Argentina's tax system
2. **Official Documentation** links and references
3. **Argentina-specific Requirements**:
   - Tax identification requirements
   - Invoice numbering rules
   - Electronic invoicing obligations
   - Provincial tax considerations
4. **Code Examples** for common scenarios
5. **Extension Usage** for special cases

### Code Documentation

- **Reference all official sources** in code comments
- **Include effective dates** for tax rates and rules
- **Document validation algorithms** with official references
- **Explain business logic** with regulatory context

## Quality Assurance

### Validation Checklist

- [ ] All tax rates match official AFIP sources
- [ ] Tax ID validation follows official algorithms
- [ ] Business rules comply with current regulations
- [ ] All official sources are referenced in code
- [ ] Comprehensive test coverage
- [ ] Documentation is complete and accurate
- [ ] Code follows GOBL patterns and conventions

### Review Process

1. **Legal Review**: Verify compliance with Argentine tax law
2. **Technical Review**: Ensure code quality and GOBL compliance
3. **Testing Review**: Validate test coverage and accuracy
4. **Documentation Review**: Confirm completeness and accuracy

## Examples and References

Study existing regimes for implementation patterns:

- **Spain (ES)**: Complex VAT system with multiple rates
- **Brazil (BR)**: Multiple tax categories and provincial considerations
- **Mexico (MX)**: VAT and special taxes
- **Italy (IT)**: Complex tax identity validation
- **India (IN)**: Dual tax system (similar to federal/provincial in Argentina)

## Getting Started

1. **Research Phase**: Study official AFIP documentation
2. **Design Phase**: Plan the regime structure and components
3. **Implementation Phase**: Start with core files and basic validation
4. **Testing Phase**: Develop comprehensive tests
5. **Documentation Phase**: Create detailed README and code documentation
6. **Review Phase**: Legal and technical review
7. **Integration Phase**: Add to GOBL registry and test integration

## Important Notes

- **Always use official sources** for tax rates and validation rules
- **Include effective dates** for all tax rates and rules
- **Reference official documents** in all code comments
- **Test with real-world examples** where possible
- **Consider future changes** in tax law and regulations
- **Follow GOBL conventions** for consistency across regimes

This guide provides the foundation for implementing a robust, compliant Argentina tax regime for GOBL that follows established patterns while meeting Argentina's specific requirements.
