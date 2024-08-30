# 🇭🇺 GOBL Hungary Tax Regime

## Public documentation

- [API documentation and e-reporting format](https://onlineszamla-test.nav.gov.hu/files/container/download/Online%20Invoice%20System%203.0%20Interface%20Specification.pdf)

## Hungary Specifics
E-invoicing is not mandatory but e-reporting is. However, since 2021 the XML file submitted to comply with the real-time reporting obligations can be delivered by the tax authority to the customer and used as an e-invoice. To this end, the issuers must indicate that it is an eInvoice, generate a hash value from the invoice data and insert it into the XML file. In addition to the data mandatory for Real-Time Information Reporting (RTIR) system, all data required for invoices must be included into the XML file. 

In Hungary there are not different types of invoices with specific names. To create a credit note, it is the same as a regular invoice with the only difference that the values of the line items are negative and that you should include the reference number of the invoice to be modified. 

In Hungary there are VAT groups. If a supplier/customer is a VAT group member, the VAT group ID must be placed as the tax ID of the supplier/customer and there must be another field with the taxID of the group member. To support this in GOBL we include the second VAT ID in the Identities field. These VAT IDs are easy to differentiate as VAT Group ID has a 5 as the VAT Code (9th digit of the VAT ID) and the group member VAT ID has a 4.

The customer must be classified in one of these groups: DOMESTIC, PRIVATE_PERSON or OTHER. 

