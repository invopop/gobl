# Israel (IL) Tax Regime

This document provides an overview of the tax regime in Israel.

Find example IL GOBL files in the [`examples`](../../examples/il) (uncalculated documents) and [`examples/out`](../../examples/il/out) (calculated envelopes) subdirectories.


## Value-Added Tax (VAT)

Known in Hebrew as *Ma'am* (מע"מ), Israel's VAT system categorizes goods and services into three main rates:

- **General Rate (18%)**: Applies to most goods and services. Raised from 17% on 1 January 2025.
- **Zero Rate (0%)**: Applies to exports of goods, services provided to non-residents, supplies to the Eilat Free Trade Zone, fresh fruits and vegetables, and certain tourism-related services.
- **Exempt**: Certain financial services, residential real estate leases (up to 25 years), educational and healthcare services, and services provided by non-profit organisations are exempt from VAT.

VAT is administered by the **Israel Tax Authority (ITA)** under the VAT Law 5736-1975. For more information, visit the [Israel Tax Authority website](https://www.gov.il/en/departments/israel_tax_authority).

**Note:** Financial institutions pay an equivalent tax at the standard rate based on total payroll and profits, while non-profit organisations with *Malkar* status pay a wage tax of 7.5% on total payroll.

## VAT Registration Requirements

Israeli businesses are classified into the following categories based on their annual turnover:

- **Authorized Dealer (Osek Murshe — עוסק מורשה)**: Businesses with annual turnover above **NIS 120,000** must register, charge VAT at the general rate, and file bimonthly returns.
- **Exempt Dealer (Osek Patur — עוסק פטור)**: Businesses below the threshold are not required to register, cannot charge VAT, and may not recover input VAT.
- **Small Dealer (Osek Zair — עוסק זעיר)**: Introduced in 2025, an alternative to Osek Patur for businesses below the threshold with no employees, offering a simplified 30% automatic expense deduction.

**Note:** Businesses below the mandatory threshold are not permitted to register as an Authorized Dealer and will therefore issue invoices without a tax ID.

### Tax ID Validation

Israeli VAT-registered businesses are identified by a **Mispar Osek Murshe** (מספר עוסק מורשה), a 9-digit numeric identifier assigned by the ITA upon registration.

For sole proprietors (Osek Murshe), the Mispar Osek is typically the same as the personal **Mispar Zehut** (מספר זהות / Teudat Zehut), a 9-digit national ID number. For companies and other legal entities, the Mispar Osek corresponds to the entity's registration number issued by the **Corporations Authority** (רשות התאגידים), where the first two digits indicate the entity type:
- `50` — Public institutions
- `51` — Companies
- `56` — Foreign non-profit corporations
- `58` — Associations (Amutot)

### Checksum validation

The personal Mispar Zehut is known to use the **Luhn algorithm (mod 10)** for its check digit. However, no official ITA source has been found confirming that the same algorithm applies to all Mispar Osek numbers, particularly those assigned to companies and other entity types. Some third-party sources suggest that company numbers may use a different check digit scheme (two control digits instead of one), but this has not been verified against official documentation.

For this reason, the current implementation validates only the 9-digit numeric format. Full verification of whether a number is active and registered must be performed through the [official government entity register](https://www.gov.il/en/service/search-the-entity-register).

## VAT Invoicing Requirements

There are two types of VAT invoices in Israel: the standard tax invoice and the simplified tax invoice.

**Simplified Tax Invoice**: Allowed in the following cases, per Section 46 of the VAT Law 5736-1975:

- When the recipient of goods or services is **not a registered dealer** (Osek Murshe).

## Electronic Invoicing — SHAAM

Since May 2024, Israel requires B2B invoices above certain thresholds to be pre-cleared with the ITA via the **SHAAM platform** (מערכת שידור חשבוניות מס) before being issued. Upon approval, the ITA assigns an **Allocation Number** (מספר הקצאה) that must appear on the invoice for the buyer to deduct input VAT.

The threshold has decreased in phases:

| Date | Threshold (pre-VAT) |
|------|---------------------|
| May 2024 | NIS 25,000 |
| January 2025 | NIS 20,000 |
| January 2026 | NIS 10,000 |
| June 2026 | NIS 5,000 |

**Note:** SHAAM integration is not yet implemented as a GOBL addon. This is planned for a future release.

### References

- [Israel Tax Authority](https://www.gov.il/en/departments/israel_tax_authority)
- [VAT Law 5736-1975 (English)](https://www.icnl.org/wp-content/uploads/Israel_vat1975.pdf)
- [OECD — Israel TIN](https://www.oecd.org/tax/automatic-exchange/crs-implementation-and-assistance/tax-identification-numbers/Israel-TIN.pdf)
- [ITA Invoice API Specification v1.0 (July 2023)](https://www.gov.il/BlobFolder/generalpage/israel-invoice-160723/he/IncomeTax_software-houses-en-040723.pdf)
