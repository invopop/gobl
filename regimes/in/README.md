# India (IN) Tax Regime

This document provides an overview of the tax regime in India.

---

## Overview of GST

India follows a **dual GST model**, where both the Central and State Governments levy taxes on a shared tax base:

- **Central GST (CGST)**: Levied by the Central Government.
- **State GST (SGST) / Union Territory GST (UTGST)**: Levied by State or Union Territory Governments.
- **Integrated GST (IGST)**: Levied by the Central Government on **interstate supplies** and imports.

### Application of GST

- **Intrastate Supplies**: Subject to CGST and SGST/UTGST in equal proportions.
- **Interstate Supplies and Imports**: Subject to IGST, equivalent to CGST + SGST.
- **Compensation Cess**: Additional tax on luxury and sin goods, such as tobacco and motor vehicles.

---

## Rates and Categories

### Taxable Rates

1. **0.25%–3%**: Precious metals like gold and diamonds.
2. **5%**: Basic goods and services (e.g., economy air travel, basic restaurants).
3. **12%–18%**: Standard services and goods (e.g., hotels, banking, construction).
4. **28%**: Luxury items (e.g., air conditioners, motor vehicles).

### Zero-Rated Supplies

- Exports.
- Supplies to Special Economic Zones (SEZs).

### Exempt Supplies

- Fresh fruits and vegetables.
- Educational services.
- Public road tolls.

### Note on GOBL Tax Categories

Due to the **dual GST model**, which divides taxes between the Central and State Governments, GOBL does not include predefined rate values for tax categories (e.g., CGST, SGST/UTGST, IGST). This choice prioritizes simplicity, avoiding the added complexity of managing split tax rate allocations.

---

### GSTIN (Goods and Services Tax Identification Number)

The GSTIN is a unique 15-digit identifier assigned to every registered taxpayer under the GST system.

#### Validation

GOBL includes built-in validation for the GSTIN field to ensure compliance with the GST system. This validation verifies the format and checksum of the GSTIN, ensuring that only correctly structured identifiers are accepted.

---

### HSN (Harmonized System of Nomenclature) Code

The HSN code is an internationally recognized system for classifying goods, adopted in India under the GST regime. It is used to identify and categorize items systematically, ensuring the correct application of tax rates.

Under the Indian GST system, the HSN code is mandatory for each item on a tax invoice, helping maintain uniformity and compliance across goods and services transactions.

---

Find example IN GOBL files in the [`examples`](../../examples/in) (uncalculated documents) and [`examples/out`](../../examples/in/out) (calculated envelopes) subdirectories.

For additional details, visit the official [GST Portal](https://www.gst.gov.in/).
