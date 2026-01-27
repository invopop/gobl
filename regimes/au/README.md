# 🇦🇺 GOBL Australia Tax Regime

Find example AU GOBL files in the [`examples`](../../examples/au) (uncalculated documents) subdirectory.

## Tax System

Australia operates a Goods and Services Tax (GST) system administered by the Australian Taxation Office (ATO). GST was introduced on 1 July 2000 as a broad-based consumption tax.

### GST Rates

**Standard rate: 10%**

Applies to most goods and services sold or consumed in Australia.

**GST-free (zero-rated): 0%**

Certain supplies are GST-free, meaning GST is charged at 0% but input tax credits can still be claimed. Common GST-free supplies include:

- Exports of goods and services
- Basic food items (bread, milk, eggs, fruit, vegetables, meat)
- Most health and medical services
- Medical aids and appliances
- Educational courses and course materials
- Childcare services
- Precious metals
- Going concerns (business sales)

**Input-taxed (exempt)**

Some supplies are input-taxed, meaning no GST is charged and no input tax credits can be claimed on related purchases:

- Most financial supplies (lending, credit, equity transactions)
- Residential rent
- Sales of residential premises
- Certain insurance products

## Tax Invoices

The requirements for tax invoices depend on the taxable value.

### Tax Invoices Under A$1,000

For sales with a taxable value less than A$1,000 (GST-exclusive), a tax invoice must show:

1. The seller's identity or business name
2. The seller's ABN
3. The date of issue
4. A brief description of items sold
5. The GST amount or a statement that the total price includes GST

### Tax Invoices A$1,000 or More

For sales with a taxable value of A$1,000 or more (GST-exclusive), a tax invoice must include all of the above, plus:

- **The buyer's identity (name) OR the buyer's ABN**

This is a key requirement - either the buyer's name or their ABN must appear on the invoice.

Additional requirements for invoices A$1,000 and above:

- Quantity of goods or extent of services
- Price per item (if applicable)
- The GST amount (can be shown separately or as "Total price includes GST")

## Australian Business Number (ABN)

The ABN is an 11-digit identifier for businesses registered with the Australian Business Register (ABR).

### Format

- 11 digits
- Example: `51 824 753 556`

### Validation

ABN validation uses a modulus 89 weighted checksum algorithm:

1. Subtract 1 from the first (left-most) digit
2. Multiply each digit (including the modified first digit) by its weight: `[10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]`
3. Sum all the products
4. Divide the sum by 89
5. If the remainder is zero, the ABN is valid

**Example for ABN 51824753556:**

```
Step 1: 5 - 1 = 4
Modified digits: [4, 1, 8, 2, 4, 7, 5, 3, 5, 5, 6]
Weights:        [10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
Products:       [40, 1, 24, 10, 28, 63, 55, 39, 75, 85, 114]
Sum:            534
534 % 89 = 0 ✓ Valid
```

## GST Credits (Input Tax Credits)

Registered businesses can claim GST credits for GST paid on purchases used in making taxable supplies. This ensures GST is only paid on the value added at each stage.

**Key distinction:**

- **GST-free supplies**: GST charged at 0%, but input tax credits can be claimed
- **Input-taxed supplies**: No GST charged, and no input tax credits can be claimed

This makes GST-free supplies (like exports) more favorable for businesses, as they can recover GST paid on costs while charging 0% to customers.

## References

**Australian Taxation Office:**
- [GST Overview](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst)
- [When to charge GST](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/when-to-charge-gst-and-when-not-to)
- [Tax Invoices](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/tax-invoices)

**Australian Business Register:**
- [ABN Format](https://abr.business.gov.au/Help/AbnFormat)
- [ABN Lookup](https://abr.business.gov.au/)
