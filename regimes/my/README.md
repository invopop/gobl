# 🇲🇾 GOBL Malaysia Tax Regime

Find example MY GOBL files in the [`examples`](../../examples/my) (uncalculated documents) and [`examples/out`](../../examples/my/out) (calculated envelopes) subdirectories.

## Malaysia-specific Requirements

### Tax IDs

In Malaysia, companies and taxable entities are primarily identified using:

- **Business Registration Number**: A 12-digit number issued upon company registration. Example: `201901234567`
- **Service Tax / Sales Tax Registration Numbers**: For companies registered under SST (Sales and Service Tax), tax IDs may appear in alphanumeric formats such as:
  - `SST1234567890`
  - `W10-12345678-123`
  - These identifiers may include letters and hyphens.

During the normalization process of Tax Identities, GOBL will automatically **uppercase** Malaysian tax IDs for consistency.

A Malaysian company's `org.Party` definition inside an invoice may look like:

```json
{
  "tax_id": {
    "country": "MY",
    "code": "201901234567"
  },
  "name": "Tech Solutions Sdn Bhd",
  "addresses": [
    {
      "street": "123 Jalan Bukit Bintang",
      "locality": "Kuala Lumpur",
      "region": "Wilayah Persekutuan",
      "code": "55100",
      "country": "MY"
    }
  ],
  "emails": [
    {
      "addr": "billing@techsol.my"
    }
  ],
  "identities": [
    {
      "type": "SST",
      "code": "SST1234567890"
    }
  ]
}
```

**Notes:**
- The primary `tax_id` field must always be filled.
- Additional tax identifiers like SST numbers can be included in the `identities` array for clarity or regulatory needs.
- Both numeric-only and alphanumeric tax IDs are accepted, following Malaysian Customs regulations.
- The logic to validate or normalize values in the identities array (e.g., SST numbers) has not been implemented yet.

## Disclaimer

> This implementation of the Malaysian tax regime (`MY`) for GOBL has been developed as part of a technical exercise and is intended for **demonstration and testing purposes only**. While every effort has been made to align with publicly available SST regulations and business practices in Malaysia, this module is **not certified for production use** and may not cover all legal, fiscal, or e-invoicing requirements.
>
> Users are advised to consult with a certified tax advisor or local regulatory authority before using this module in any commercial or legal context.


Sources:
- [Royal Malaysian Customs Department - Sales and Service Tax](https://mysst.customs.gov.my/)


