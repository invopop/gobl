# 🇸🇪 GOBL Sweden Tax Regime

Sweden uses the PEPPOL BIS Billing 3.0 format (CIUS based on EN 16931) for their e-invoicing system.

Find example SEGOBL files in the [`examples`](../../examples/se) (uncalculated documents) and [`examples/out`](../../examples/se/out) (calculated envelopes) subdirectories.

## Public Documentation

- [PEPPOL BIS Billing 3.0 Specification](https://docs.peppol.eu/poacc/billing/3.0/)
- [Agency for Digital Government (DIGG) - E-invoicing](https://www.digg.se/e-handel-och-e-faktura/obligatorisk-e-fakturering-i-offentlig-sektor)
- [Swedish Tax Agency (Skatterverket)](https://www.skatteverket.se/foretag.4.76a43be412206334b89800052908.html)
  - [English version](https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers.4.12815e4f14a62bc048f5159.html)
- [Swedish Tax Agency (Skatteverket) - VAT Rules](https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momssatspavarorochtjanster.4.58d555751259e4d66168000409.html)
- [2025 Official Skatterverket VAT Guide](https://www.skatteverket.se/download/18.7be5268414bea0646946f3e/1428566850726/552B14.pdf)

## E-invoicing Requirements

Since April 1, 2019, all suppliers to Swedish public authorities must send and receive invoices electronically via the PEPPOL network using the PEPPOL BIS Billing 3.0 profile. While not mandatory for B2B transactions, e-invoicing is highly recommended and effectively a de-facto standard when trading with larger enterprises.

A compliant invoice must include all Core Invoice elements as specified by EN 16931 Business Rules and the PEPPOL profile, including:

1. Date of issue of the invoice.
2. A unique serial number for each invoice based on one or more series (unique and sequential per fiscal year).
3. The seller's VAT registration number.
4. The buyer's VAT registration number if the buyer is liable for payment for the purchase, so-called reverse charge.
5. Name and address of the seller and buyer.
6. The quantity and nature of the goods or the scope and nature of the services.
7. The date on which the sale of the goods or services was made or completed or the date on which the advance or on-account payment was made, if such a date can be determined and it is different from the invoice date.
8. The tax base for each VAT rate or exemption, the unit price excluding VAT, and any price reduction or discount not included in the unit price.
9. Applied VAT rate.
9. The amount of VAT to be paid. If the seller uses profit margin taxation on the transaction, the VAT amount should not be stated on the invoice.

[Source](https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momslagensregleromfakturering.4.58d555751259e4d66168000403.html#fakturansinnehall)

## Sweden-specific Requirements

[Source](https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/394342.html)

### Identification Numbers

**All identification numbers are 10 digits long**. Individuals and businesses have different formats.

The **supplier must always be identified by their tax identification number**.

If the customer is liable to pay VAT (they are registered for VAT or the operation is a reverse charge), their VAT number must also be present. Otherwise name and address are sufficient.

> **VAT ID format:** `SE` + 10-digit organization or personal number + `01` as check digits

#### Individuals:

- **Social security number (Personnummer)** - For individuals registered in Sweden.
  - Birthdate in format YYMMDD, "-" if less than 100 years old, "+" otherwise, 3-digit birthnumber (last digit odd for biological males, even for females) and a checksum.
  - Example: `990101-1230` Male born in 1999-01-01
- **Coordination number (Samordningsnummer)** - For individuals not registered in Sweden.
  - Birthdate in format YYMMDD (day + 60), "-" if less than 100 years old, "+" otherwise, 3-digit birthnumber (last digit odd for biological males, even for females) and a checksum.
  - Example: `990161+1229` Female born in 1899-01-01

#### Businesses:

- **Organization number (Organisationsnummer, `Org.nr`)**
  <details>
    <summary>First digit (the "group number") may identify the type of entity (but not necessarily):</summary>
    <ul>
      <li>1 - Death certificate</li>
      <li>2 or 8 - Religious denominations</li>
      <li>20 - Government agencies (assigned by the Statistics Sweden)</li>
      <li>3 - Foreign companies engaged in business activities or own real estate in Sweden</li>
      <li>5 - limited liability companies (Aktiebolag), branches, banks, insurance companies and European companies</li>
      <li>6 - Single company</li>
      <li>7 or 8 - Tenant-owner associations, economic associations, non-profit associations, housing associations, cooperative tenancy associations, European cooperatives and European groupings for territorial cooperation</li>
      <li>9 - Partnerships and limited partnerships</li>
    </ul>
  </details>
  - The last digit is a checksum.
  - Sole proprietorships use the owner's personal number.
  - Example: `556036-0793` Private limited company

#### Checksum

<details>
<summary>The checksum is calculated using the <a href="https://stripe.com/en-es/resources/more/how-to-use-the-luhn-algorithm-a-guide-in-applications-for-businesses">Luhn algorithm</a>:</summary>

- Start with the payload digits. Moving from right to left, double every second digit, starting from the last digit. If doubling a digit results in a value > 9, subtract 9 from it (or sum its digits).
- Sum all the resulting digits (including the ones that were not doubled).
- The check digit is then calculated using the formula $(10 - (s \pmod{10})) \pmod{10}$, where $s$ is the sum from the previous step. This yields the smallest non-negative number which, when added to $s$, results in a multiple of 10.

</details>

### VAT Rates

In Sweden, VAT is called "Moms" (Mervärdesskatt). The following rates (Skattesatser) are used in Sweden:

| Rate            | Swedish Term          | Percentage | Description                                                                                                              |
| --------------- | --------------------- | ---------- | ------------------------------------------------------------------------------------------------------------------------ |
| Standard        | Normalskattesats      | 25%        | Most goods and services                                                                                                  |
| Reduced         | Skattesats 12 procent | 12%        | Food products, hotel accommodations, restaurant and catering services, shoe repair, leather goods, clothing, bicycles... |
| Heavily reduced | Skattesats 6 procent  | 6%         | Passenger transport, intellectual property, cultural services (except cinema), books, newspapers...                      |
| Exempt          | Momsfri               | 0%         | Exports, intra-community supplies, pharmaceuticals, certain financial and healthcare services                            |

### F-Tax

When a business is registered for F-tax, their customers do not have to deduct taxes on payments made to them for work performed in Sweden.

This is usually done by domesticsole proprietorships and foreign companies.

### Reverse Charge

Cases where reverse charge applies:

- **Domestic Transactions**: Applicable to specific sectors like construction services, trading of certain metals, waste and scrap materials, emission rights, and services related to real estate.

- **Cross-Border Transactions**:
  - Intra-Community Supplies: When goods or services are supplied between EU member states, and the customer is VAT-registered in another member state.
  - Services from Abroad: When services are provided by a supplier not established in Sweden to a VAT-registered customer in Sweden.

Implications:

- VAT rate = 0
- VAT category code "Reverse charge"
- Include both supplier and customer VAT IDs
- Add note or exemption reason "Reverse charge"
- Customer must account for VAT

### PEPPOL Technical Requirements

Invoices must reference the PEPPOL BIS identifiers:
- Specification Identifier: `urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0`
- Business Process Specified Document Context Parameter: `urn:fdc:peppol.eu:2017:poacc:billing:01:1.0`

No empty XML elements are allowed, and all mandatory fields must be present to pass validation.
