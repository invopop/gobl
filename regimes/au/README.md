# 🇦🇺 GOBL Australia Tax Regime

This document provides an overview of the tax regime in Australia.

Find example AU GOBL files in the [`examples`](../../examples/au) (uncalculated documents) subdirectory.

## Tax Regime

Australia operates a Goods and Services Tax (GST) system administered by the Australian Taxation Office (ATO). GST was introduced on 1 July 2000 and is a broad-based tax of 10% on most goods, services and other items sold or consumed in Australia.

### GST Rates

Australia has two GST rates:

1. **Standard rate: 10%** - Applies to most goods and services
2. **Zero-rated (GST-free): 0%** - Applies to:
   - Basic food items (bread, milk, eggs, vegetables, fruit, meat, etc.)
   - Most health and medical services
   - Medical aids and appliances
   - Educational courses and course materials
   - Childcare services
   - Exports of goods and services
   - Precious metals
   - Going concerns (business sales)

**Note:** Input-taxed (exempt) supplies are not currently modeled in this v1 implementation. Input-taxed supplies include financial services, residential rent, and residential property sales.

### Tax Invoices

For tax invoices with a taxable value of **A$1,000 or more** (GST-exclusive), the invoice must include either:
- The buyer's name, OR
- The buyer's ABN (Australian Business Number)

For invoices under A$1,000, the buyer's name is recommended but not legally required.

## Australian Business Number (ABN)

The ABN is an 11-digit identifier for businesses registered with the Australian Business Register (ABR). All businesses operating in Australia must have an ABN.

### Format

- **Length:** 11 digits
- **Example:** `51 824 753 556` (commonly formatted with spaces, but stored without)

### Validation

ABN validation uses a modulus 89 weighted checksum algorithm:

1. Subtract 1 from the first (left-most) digit
2. Multiply each digit by its corresponding weight: `[10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]`
3. Sum all the products
4. If the sum divided by 89 has a remainder of 0, the ABN is valid

**Example validation for ABN `51824753556`:**
```
Step 1: 5 - 1 = 4
Modified digits: [4, 1, 8, 2, 4, 7, 5, 3, 5, 5, 6]
Weights:        [10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
Products:       [40, 1, 24, 10, 28, 63, 55, 39, 75, 85, 114]
Sum:            534
534 % 89 = 0 ✓ Valid
```

### Normalization

GOBL automatically normalizes ABNs by:
- Removing spaces, hyphens, and dots
- Removing "AU" country prefix if present
- Result: 11-digit numeric string

## Scope

This v1 implementation includes:

✅ **Implemented:**
- Standard 10% GST rate
- Zero-rated (GST-free) supplies at 0%
- ABN validation with modulus 89 checksum
- Invoice validation for A$1,000 threshold rule
- Comprehensive test coverage

❌ **Not yet implemented (deferred to v2):**
- Input-taxed (exempt) supplies
- Special GST schemes (margin scheme, second-hand goods, etc.)
- Luxury Car Tax (LCT)
- Wine Equalisation Tax (WET)
- State-based taxes (payroll tax, land tax, stamp duty)
- Tax periods and BAS (Business Activity Statement) reporting

## Testing

Run tests for the Australia regime:

```bash
go test ./regimes/au/...
```

Run with verbose output:

```bash
go test -v ./regimes/au/...
```

Run with coverage:

```bash
go test -cover ./regimes/au/...
```

## Examples

See example invoices in [`examples/au/`](../../examples/au):

- `invoice-simple.json` - Basic invoice under A$1,000 threshold
- `invoice-over-threshold.json` - Invoice over A$1,000 with buyer ABN
- `invoice-gst-free.json` - Invoice with GST-free (zero-rated) exports

To validate an example:

```bash
gobl validate examples/au/invoice-simple.json
```

## References

- [Australian Taxation Office - GST](https://www.ato.gov.au/business/gst)
- [ATO - When to charge GST](https://www.ato.gov.au/business/gst/when-to-charge-gst)
- [ATO - Issuing tax invoices](https://www.ato.gov.au/business/gst/issuing-tax-invoices)
- [ABN Format and validation](https://abr.business.gov.au/Help/AbnFormat)
- [Australian Business Register](https://abr.business.gov.au/)
