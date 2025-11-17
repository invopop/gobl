# 🇸🇬 GOBL Singapore Tax Regime

This document provides an overview of the tax regime in Singapore

Find example SG GOBL files in the [`examples`](../../examples/sg) (uncalculated documents) and [`examples/out`](../../examples/sg/out) (calculated envelopes) subdirectories.

## Rates

1. Standard rate of **9%**. (Since 01/01/2024)
2. Zero-rate which applies to international services and export of goods.
3. Exempt Supplies which include financial services, sale and lease of residential properties, digital payment tokens, and the import of investment precious metals.
4. Out-of-scope supplies refer to supplies which are outside the scope of the GST Act. Some examples are salaries paid to employees for their services or supplies where the place of supply is outside of Singapore.

In Singapore there are some schemes that allow reduced rates:

- **Discounted Sale Price Scheme**: When you sell a second-hand or used vehicle using this scheme, you can charge GST on 50% of the selling price.
- **Gross Margin Scheme**: GST is accounted for on the gross margin (i.e. selling price less purchase price) instead of full value of the goods supplied.

*We are not considering these for now*

## Invoicing methods

### Tax Invoice
This invoice, reference in GOBL by the use of the tax tag "standard" represents a basic Invoice. This Invoice has to meet certain requirements:

1. The words “tax invoice” in a prominent place.
2. An identifying number (e.g. invoice number).
3. Date of issue of the invoice.
4. Supplier business name, address and GST registration number.
5. Customer’s name and address.
6. A description sufficient to identify the goods or services supplied and the type of supply.
7. For each description of goods or services supplied, the quantity of goods or the extent of services, and the amount payable, excluding GST.
8. Any cash discount offered.
9. The total amount payable (excluding GST), the GST rate and the total amount of GST chargeable.
10. The total amount payable (including the total amount of GSTchargeable).
11. A breakdown of exempt, zero-rated or other supplies, stating separatelythe gross total amount payable in respect of each type of supply.

A tax invoice must not be issued if:
- You are not registered for GST
- You are selling goods using the Gross Margin Scheme (GMS)
- You are the supplier in a self-billing arrangement where your customer
issues the tax invoice

### Simplified Tax Invoice

This invoice is referenced by the tax tag "simplified". This invoice can only be used when the total amount (inclusive of GST) is less than 1000 SGD. This invoice has less requirements:

1. Supplier name, address and GST registration number;
2. An identifying number, e.g. invoice number.
3. The date of issue of the invoice.
4. Description of the goods or services supplied.
5. The total amount payable including tax.
6. The word “Price Payable includes GST”.

### Receipt

This type of invoice can be issued to a non-GST registered customer. A receipt must be serially printed and must show the following:

1. Supplier name and GST registration number;
2. The date of issue of the invoice.
3. The total amount payable including tax.
4. The word “Price Payable includes GST”.

## Correction method

### Credit Note
A credit note is issued to correct a mistake or to give a credit to your customer. A credit note must include:

1. An identifying number e.g. a serial number
2. Date of issue
3. Your name, address and GST registration number
4. Your customer's name and address
5. Reason for the credit, e.g. "returned goods"
6. Detailed description to identify the goods and services that credit is allowed for
7. Quantity and amount credited for each description
8. Total amount credited, excluding tax
9. Rate and amount of tax credited
10. Total amount credited, including tax
11. Number and date of the original tax invoice

### Debit Note
You should only issue a debit note to request for payment for transactions where no GST is charged (e.g. internal billings within the same company), or to suppliers from whom credit is due. **Not used for correcting an invoice**

## Singapore tax IDs
Suppliers in Singapore may use several official tax identification numbers on invoices, depending on their entity type and tax status:

- **Unique Entity Number (UEN)**: All registered business entities in Singapore, not specific to GST. Every GST-registered business has a UEN, but UEN itself is the company’s main ID, not a GST number. 
- **GST Registration number**: Any business entity that is registered for GST with IRAS. Overseas suppliers who register for GST also receive one. GST-registered suppliers are required to print their GST Registration Number on every tax invoice and receipt issued. In most cases GST registration number is the same as the UEN.

When the numbers differ, the GST registration number will be included in identities.

There are other tax IDs in Singapore like the NRIC or FIN, but they are used mainly for personal identification, and not mandated on invoices.


### References

[GST General Guide for Businesses](https://www.iras.gov.sg/media/docs/default-source/e-tax/etaxguide_gst_gst-general-guide-for-businesses(1).pdf?sfvrsn=8a66716d_97)

[GST Rates](https://www.iras.gov.sg/taxes/goods-services-tax-(gst)/basics-of-gst/current-gst-rates)




