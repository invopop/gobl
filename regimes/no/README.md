# Norway (NO) Tax Regime

This document provides an overview of the Norwegian tax regime, with an emphasis on the Value-Added Tax (VAT) system, registration requirements, invoicing rules, and related administration.

---

## Overview

Norway's tax system is administered by the Norwegian Tax Administration (Skatteetaten). Although Norway is not part of the European Union, it maintains its own VAT system that applies broadly to most goods and services. In addition to VAT, Norwegian businesses are subject to corporate tax, income tax, and other duties. This document focuses on VAT-related matters which are critical for most enterprises operating in Norway.

---

## Value-Added Tax (VAT)

Norway levies VAT on goods and services with several distinct rates. The rates, which are subject to periodic review, currently include:

- **Standard Rate (25%)**  
  Applies to most goods and services.

- **Reduced Rate (15%)**  
  Applies primarily to food and non-alcoholic beverages.

- **Reduced Rate (12%)**  
  Applies to:
  - Passenger transport (including local public transport)
  - Hotel accommodation
  - Cultural events (e.g., cinema, theater, exhibitions)
  - Certain services such as vehicle ferry transport and broadcasting
  
- **Zero Rate (0%)**  
  Applies to:
  - Exports of goods and certain international transport services
  - Specific activities defined in Norwegian law (e.g., cross-border supplies)

- **Exemptions**  
  Some goods and services are exempt from VAT. These typically include:
  - Financial and banking services
  - Insurance
  - Certain healthcare and educational services

*Note: The precise application of rates, as well as exemptions, may be subject to special rules or evolving policy. Always consult the latest guidance from Skatteetaten.*

---

## VAT Registration Requirements

### Who Must Register?
- **Mandatory Registration**:  
  Any business whose taxable turnover exceeds **NOK 50,000** within a 12-month period must register for VAT. Taxable turnover includes both sales and purchases of goods and services that are subject to VAT.

- **Voluntary Registration**:  
  Businesses with a turnover below **NOK 50,000** may opt to register voluntarily. Voluntary registration allows the business to:
  - Reclaim VAT on business-related purchases (input VAT).
  - Issue VAT invoices to customers, even if not strictly required by law.

### Registration Process
- Businesses register through the online portal provided by Skatteetaten.
- Upon approval, the business is issued an Organisation Number appended with "MVA" (e.g., `123456789 MVA`).

---

## Tax Registration Number (TRN)

The Tax Registration Number, commonly reflected as an Organisation Number followed by "MVA", is a unique identifier for businesses registered for VAT.

- **Format**: A 9-digit number (e.g., `123456789`) that, when appended with "MVA", becomes the VAT number.
- **Validation**:
  - The TRN uses a checksum process (similar in concept to the Luhn algorithm) to ensure its validity.
  - The algorithm typically involves multiplying digits by weight factors, summing the results, and ensuring that the total modulo 10 is zero.
  
*Example*:  
For the TRN `123456789`, an illustrative checksum is calculated as:  
`(1×2 + 2×1 + 3×2 + 4×1 + 5×2 + 6×1 + 7×2 + 8×1 + 9×2) % 10 = 0`  
If the checksum equals zero, the TRN is considered valid.

*Note: The specifics of the weight factors may vary. Businesses should refer to Skatteetaten’s guidelines on TRN validation for complete details.*

---

## VAT Filing and Payment

Businesses registered for VAT must adhere to strict filing and payment schedules. The filing frequency is determined by the business’s annual taxable turnover:

- **Monthly Filing**:  
  For businesses with an annual taxable turnover exceeding **NOK 6 million**.  
  **Due Date**: Typically, the return and payment are due on the 10th day of the month following the reporting period.

- **Quarterly Filing**:  
  For businesses with an annual taxable turnover between **NOK 1 million** and **NOK 6 million**.  
  **Due Date**: Returns and payments are generally due on the 10th day of the month following each quarter.

- **Annual Filing**:  
  For businesses with an annual taxable turnover below **NOK 1 million**.  
  **Due Date**: The VAT return and payment are due on the 10th day of the month after the end of the financial year.

### Payment Methods and Penalties

- **Payment Methods**:  
  VAT payments can be made electronically via Skatteetaten’s online system or by bank transfer.

- **Penalties and Interest**:  
  Late filing or delayed payment may result in fines and interest charges. The interest rates for late payments are set by Skatteetaten and are subject to periodic adjustment.

*Businesses are advised to ensure prompt filing and payment to avoid penalties and to monitor any changes in the filing frequencies or deadlines published by the tax authorities.*

---

## Invoicing Requirements

To ensure transparency and enable accurate tax reporting, Norwegian invoicing rules are comprehensive:

### Invoice Format and Content

- **Format**:
  - Invoices may be issued in either electronic or paper format.
  - Electronic invoices must comply with security and archiving requirements set by Skatteetaten.

- **Mandatory Information**:
  - **Invoice Date** and a **Unique Invoice Number**
  - **Seller’s Details**: Name, address, and VAT number (TRN).
  - **Buyer’s Details**: Name, address, and VAT number (if applicable).
  - **Description of Goods/Services**: A clear description of the items or services supplied.
  - **Quantity and Unit Price**: Breakdowns per item or service.
  - **Total Amount**: Total sum due, including VAT.
  - **VAT Rate(s) Applied**: Indicating the applicable rate for each item.
  - **Payment Terms and Due Date**

- **Currency and Exchange Rates**:
  - The invoice must specify the currency.
  - For invoices in a foreign currency, include the exchange rate applied at the time of the transaction.

- **Language Requirements**:
  - Invoices can be issued in Norwegian or English.
  - If using a foreign language, a Norwegian translation must be provided upon request.

### Special Invoicing Scenarios

- **Retention Period**:  
  Invoices must be retained for at least **5 years** from the end of the financial year during which the invoice was issued.

- **Electronic Invoicing**:
  - Must meet the technical standards and security requirements specified by Skatteetaten.
  - Must be stored securely and remain accessible for audits.

- **Credit and Debit Notes**:
  - **Credit Notes**: Must reference the original invoice and restate the mandatory details, clearly indicating the corrections.
  - **Debit Notes**: Used when additional amounts are due; they must include a reference to the original invoice along with all mandatory information.

- **Simplified Invoices**:
  - For transactions below **NOK 1,000**, simplified invoices can be issued.
  - These must include: Invoice date, unique invoice number, seller’s name and address, description of goods or services, and the total payable amount (including VAT).

---

## Additional Considerations

### Digital Reporting and Compliance

- **E-Filing and Digital Archiving**:  
  Norwegian businesses are increasingly required to use digital systems for filing VAT returns and maintaining electronic records.  
- **Integration with Accounting Software**:  
  Many ERP and accounting systems in Norway are designed to integrate directly with Skatteetaten’s reporting systems to facilitate seamless compliance.

### Regular Updates

- **Legislative Changes**:  
  Tax rates, registration thresholds, and filing rules may be subject to change. It is crucial to follow updates through Skatteetaten’s official communications and the Norwegian legal gazette.
- **Consulting Professionals**:  
  Given the complexities of tax regulations, businesses should consult tax professionals or legal advisors to ensure compliance with the latest laws.

---

## Useful Resources

- [Skatteetaten (Norwegian Tax Administration)](https://www.skatteetaten.no)
- [Norwegian VAT Guidelines](https://www.skatteetaten.no/en/business-and-organisation/vat-and-duties)
- [Official Norwegian Legal Information](https://lovdata.no/)

---

*This document is intended as a general guide. For detailed cases and personalized advice, please consult a tax professional or the Norwegian Tax Administration.*

---

This adapted document should now reflect a comprehensive and structured view of the Norwegian VAT system and related tax matters.