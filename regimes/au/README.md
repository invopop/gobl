# Australia (`AU`)

Australia uses a Goods and Services Tax (GST) system administered by the Australian Taxation Office (ATO). GOBL models the Australian regime with a 10% standard GST rate, support for zero-rated supplies, ABN validation, and invoice validation rules for supplier and customer tax identifiers.

## Public Documentation

- [ATO - GST](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst)
- [ATO - Tax invoices](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/tax-invoices)
- [ABR - ABN format](https://abr.business.gov.au/Help/AbnFormat)
- [Peppol PINT A-NZ BIS](https://docs.peppol.eu/poac/aunz/pint-aunz/bis/)

## Tax Identity (ABN)

Australian businesses are commonly identified by an Australian Business Number (ABN). The ABN is 11 digits long, usually written with spaces for display, but normalized in GOBL without separators.

Validation follows the ABR checksum algorithm:

| Position | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Weight | 10 | 1 | 3 | 5 | 7 | 9 | 11 | 13 | 15 | 17 | 19 |

1. Subtract 1 from the first digit.
2. Multiply each digit by its weight.
3. Sum the products.
4. The total must be divisible by 89.

## GST

| Rate Name | GOBL Rate Key | Percent | Since |
| --- | --- | --- | --- |
| General rate | `standard` / `general` | 10% | 2000-07-01 |
| Zero rate | `zero` / `zero` | 0% | 2000-07-01 |
