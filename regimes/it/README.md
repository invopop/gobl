# 🇮🇹 GOBL Italy Tax Regime

Italy uses the FatturaPA format for their e-invoicing system.

## Public Documentations

### FatturaPA

[Historical Documentations](https://www.fatturapa.gov.it/en/norme-e-regole/documentazione-fattura-elettronica/formato-fatturapa/)

[Schema V1.2.1 Spec Table View (EN)](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.1/Table-view-B2B-Ordinary-invoice.pdf) - by far the most comprehensible spec doc. Since the difference between 1.2.2 and 1.2.1 is minimal, this is perfectly usable.

[Schema V1.2.2 PDF (IT)](https://www.fatturapa.gov.it/export/documenti/Specifiche_tecniche_del_formato_FatturaPA_v1.3.1.pdf) - most up-to-date but difficult

[XSD V1.2.2](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.2/Schema_del_file_xml_FatturaPA_v1.2.2.xsd)

### Tax Rates

[IRPEF](https://www.agenziaentrate.gov.it/portale/imposta-sul-reddito-delle-persone-fisiche-irpef-/aliquote-e-calcolo-dell-irpef)
[IVA (VAT)](https://www.agenziaentrate.gov.it/portale/web/guest/iva-regole-generali-aliquote-esenzioni-pagamento/norme-generali-e-aliquote#:~:text=In%20Italia%20l'aliquota%20ordinaria,per%20esempio%20per%20alcuni%20alimenti)

#### Changes from 1.2.1 to 1.2.2

- Documentation changes: TD25, N1, N6.2, N7
- Addition of TD28: Acquisti da San Marino con IVA (fattura cartacea)

### Tax Definitions

[Fiscal Code (Codice Fiscale)](https://en.wikipedia.org/wiki/Italian_fiscal_code)

[VAT Number (Partita IVA)](https://en.wikipedia.org/wiki/VAT_identification_number)

[Agenzia Entrate (Tax Office) IVA Doc](https://www.agenziaentrate.gov.it/portale/web/english/nse/business/vat-in-italy)

### Italy-specific Details

#### Stamp Duty

Add an invoice-level `bill.Charge` and use `it.ChargeKeyStampDuty` as the `bill.Charge.Key`.

#### Numero REA

`Party.Registration` is used to store the Numero REA (Registro delle Imprese) of the company.
The `Office` field is used to store the Provincia (Province) of the company, the `Entry` field is used to store the Numero REA. Additionally, the share capital is stored in the `Capital` field used in conjunction with `Currency`.

#### Local Codes

FatturaPA demands very specific categorization for the type of economic activity,
document type, fund type, etc.

##### RegimeFiscale (Tax System)

The "Regime Fiscale" or tax system needs to be applied to all Italian FatturaPA invoices. The default code `RF01` will always be used unless overridden by specifying one of the following tax tags in the invoice document:

| Code | Document Tags        | Description                                                                                                                                                                    |
| ---- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| RF01 | default              | Ordinary                                                                                                                                                                       |
| RF02 | `minimum-tax-payers` | "Minimum taxpayers (Art. 1, section 96-117, Italian Law 244/07)"                                                                                                               |
| RF04 |                      | "Agriculture and connected activities and fishing (Arts. 34 and 34-bis, Italian Presidential Decree 633/72)"                                                                   |
| RF05 |                      | "Sale of salts and tobaccos (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                          |
| RF06 |                      | "Match sales (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                                         |
| RF07 |                      | "Publishing (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                                          |
| RF08 |                      | "Management of public telephone services (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                             |
| RF09 |                      | "Resale of public transport and parking documents (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                    |
| RF10 |                      | "Entertainment, gaming and other activities referred to by the tariff attached to Italian Presidential Decree 640/72 (Art. 74, section 6, Italian Presidential Decree 633/72)" |
| RF11 |                      | "Travel and tourism agencies (Art. 74-ter, Italian Presidential Decree 633/72)"                                                                                                |
| RF12 |                      | "Farmhouse accommodation/restaurants (Art. 5, section 2, Italian law 413/91)"                                                                                                  |
| RF13 |                      | "Door-to-door sales (Art. 25-bis, section 6, Italian Presidential Decree 600/73)"                                                                                              |
| RF14 |                      | "Resale of used goods, artworks, antiques or collector's items (Art. 36, Italian Decree Law 41/95)"                                                                            |
| RF15 |                      | "Artwork, antiques or collector's items auction agencies (Art. 40-bis, Italian Decree Law 41/95)"                                                                              |
| RF16 |                      | "VAT paid in cash by P.A. (Art. 6, section 5, Italian Presidential Decree 633/72)"                                                                                             |
| RF17 |                      | "VAT paid in cash by subjects with business turnover below Euro 200,000 (Art. 7, Italian Decree Law 185/2008)"                                                                 |
| RF18 |                      | Other                                                                                                                                                                          |
| RF19 |                      | "Flat rate (Art. 1, section 54-89, Italian Law 190/2014)"                                                                                                                      |

##### ModalitaPagamento (Payment Means)

The following table describes how to map the Italian payment means codes to those of GOBL. The list is somewhat based on the official mapping of the FatturaPA codes to EU Semantic invoice definition, more details available [here](https://www.agenziaentrate.gov.it/portale/documents/20143/288396/Technical+Rules+for+European+Invoicing+v2.1.pdf).

| Code | Key(s)                     | Description                                                     |
| ---- | -------------------------- | --------------------------------------------------------------- |
| MP01 | `cash`, `other`            | Cash                                                            |
| MP02 | `cheque`                   | Cheque                                                          |
| MP03 | `bank-draft`               | Banker's draft                                                  |
| MP04 | `cash+treasury`            | Cash at Treasury                                                |
| MP05 | `credit-transfer`          | bank transfer                                                   |
| MP06 | `promissory-note`          | Promissory Note                                                 |
| MP07 | `other+payment-slip`       | Pre-compiled bank payment slip                                  |
| MP08 | `card`, `online`           | Any type of payment card                                        |
| MP09 | `direct-debit+rid`         | Direct debit (RID)                                              |
| MP10 | `direct-debit+rid-utility` | Utilities direct debit (RID utenze)                             |
| MP11 | `direct-debit+rid-fast`    | Fast direct debit (RID veloce)                                  |
| MP12 | `direct-debit+riba`        | Collection order (RIBA)                                         |
| MP13 | `debit-transfer`           | Payment by notice (MAV)                                         |
| MP14 | `other+tax-receipt`        | Tax office quittance                                            |
| MP15 | `other+special-account`    | Transfer on special accounting accounts                         |
| MP16 | `direct-debit`             | Direct Debit                                                    |
| MP17 | `direct-debit+post-office` | Order for direct payment from post office account               |
| MP18 | `cheque+post-office`       | Bulletin postal account                                         |
| MP19 | `direct-debiy+sepa`        | SEPA Direct Debit (default type of direct debit)                |
| MP20 | `direct-debit+sepa-core`   | SEPA Direct Debit CORE                                          |
| MP21 | `direct-debit+sepa-b2b`    | SEPA Direct Debit B2B                                           |
| MP22 | `netting`                  | Deduction on sums already collected from previous transactions. |
| MP23 | `online+pagopa`            | PagoPA                                                          |

##### TipoDocumento (Document Type)

All Italian invoices must be identified with a specific type code defined by the FatturaPA format. The following table helps identify how GOBL will map the expected Italian code with a combination of the Invoice Type and tax tags.

| Code | Type          | Tax Tags                            | Description                                                                                                                                                                       |
| ---- | ------------- | ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| TD01 | `standard`    |                                     | Regular Invoice                                                                                                                                                                   |
| TD02 | `standard`    | `partial`                           | Advance or down payment on invoice                                                                                                                                                |
| TD03 | `standard`    | `partial`, `freelance`              | Advance or down payment on freelance invoice                                                                                                                                      |
| TD04 | `credit-note` |                                     | Credit note                                                                                                                                                                       |
| TD05 | `debit-note`  |                                     | Debit note                                                                                                                                                                        |
| TD06 | `standard`    | `freelance`                         | Freelancer Invoice - includes retained taxes                                                                                                                                      |
| TD07 | `standard`    | `simplified`                        | Simplified (no customer)                                                                                                                                                          |
| TD08 | `credit-note` | `simplified`                        | Simplified Credit Note (no customer)                                                                                                                                              |
| TD09 | `debit-note`  | `simplified`                        | Simplified Debit Note (no customer)                                                                                                                                               |
| TD16 | `standard`    | `self-billed`, `reverse-charge`     | reverse charge internal invoice integration                                                                                                                                       |
| TD17 | `standard`    | `self-billled`, `import`            | integration/self invoicing for purchase services from abroad                                                                                                                      |
| TD18 | `standard`    | `self-billed`, `import`, `goods-eu` | integration for purchase of intra UE goods                                                                                                                                        |
| TD19 | `standard`    | `self-billed`, `import`, `goods`    | integration/self invoicing for purchase of goods ex art.17 c.2 DPR 633/72                                                                                                         |
| TD20 | `standard`    | `self-billed`, `regularize`         | self invoicing for regularisation and integration of invoices (ex art.6 c.8 and 9-bis d.lgs 471/97 or art.46 c.5 D.L. 331/93)                                                     |
| TD21 | `standard`    | `self-billed`, `ceiling-exceeded`   | Self invoicing when goods are bought for export without VAT until a certain limit. If limit is exceeded, they must issue an invoice of type TD21. (Autofaturra per splafonamento) |
| TD22 | `standard`    | `self-billed`, `goods-extracted`    | extractions of goods from VAT Warehouse                                                                                                                                           |
| TD23 | `standard`    | `self-billed`, `goods-with-tax`     | extractions of goods from VAT Warehouse with payment of VAT                                                                                                                       |
| TD24 | `standard`    | `deferred`                          | "deferred invoice ex art.21, c.4, lett. a) DPR 633/72"                                                                                                                            |
| TD25 | `standard`    | `deferred`, `third-period`          | "deferred invoice ex art.21, c.4, third period lett. b) DPR 633/72"                                                                                                               |
| TD26 | `standard`    | `depreciable-assets`                | sale of depreciable assets and for internal transfers (ex art.36 DPR 633/72)                                                                                                      |
| TD27 | `standard`    | `self-billed`, TBD                  | self invoicing for self consumption or for free transfer without recourse                                                                                                         |
| TD28 | `standard`    | `self-billed`, `san-marino-paper`   | Purchases from San Marino with VAT (paper invoice)                                                                                                                                |

##### Natura (Nature)

The "Natura" code is required when identifying why a single row inside an invoice _does not_ include VAT.

| Code | Tax Rate Key                                 | Description                                                                                                                                                                                                                                              |
| ---- | -------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| N1   | `excluded`                                   | excluded pursuant to Art. 15, DPR 633/72                                                                                                                                                                                                                 |
| N2   | `not-subject`                                | not subject (this code is no longer permitted to use on invoices emitted from 1 January 2021 )                                                                                                                                                           |
| N2.1 | `not-subject+article-7`                      | not subject to VAT under the articles from 7 to 7-septies of DPR 633/72                                                                                                                                                                                  |
| N2.2 | `not-subject+other`                          | not subject – other cases                                                                                                                                                                                                                                |
| N3   | `not-taxable`                                | not taxable (this code is no longer permitted to use on invoices emitted from 1 January 2021 )                                                                                                                                                           |
| N3.1 | `not-taxable+export`                         | not taxable - exportations                                                                                                                                                                                                                               |
| N3.2 | `not-taxable+intra-community`                | not taxable - intra Community transfers                                                                                                                                                                                                                  |
| N3.3 | `not-taxable+san-marino`                     | not taxable - transfers to San Marino                                                                                                                                                                                                                    |
| N3.4 | `not-taxable+export-supplies`                | not taxable - transactions treated as export supplies                                                                                                                                                                                                    |
| N3.5 | `not-taxable+declaration-of-intent`          | not taxable - for declaration of intent                                                                                                                                                                                                                  |
| N3.6 | `not-taxable+other`                          | not taxable – other transactions that don’t contribute to the determination of ceiling                                                                                                                                                                   |
| N4   | `exempt`                                     | exempt                                                                                                                                                                                                                                                   |
| N5   | `margin-regime`                              | margin regime / VAT not exposed on invoice                                                                                                                                                                                                               |
| N6   | `reverse-charge`                             | "reverse charge (for transactions in reverse charge or for self invoicing for purchase of extra UE services or for import of goods only in the cases provided for) — (this code is no longer permitted to use on invoices emitted from 1 January 2021 )" |
| N6.1 | `reverse-charge+scrap`                       | reverse charge - transfer of scrap and of other recyclable materials                                                                                                                                                                                     |
| N6.2 | `reverse-charge+precious-metals`             | reverse charge - transfer of gold and pure silver pursuant to law 7/2000 as well as used jewelery to OPO                                                                                                                                                 |
| N6.3 | `reverse-charge+construction-subcontracting` | reverse charge - subcontracting in the construction sector                                                                                                                                                                                               |
| N6.4 | `reverse-charge+buildings`                   | reverse charge - transfer of buildings                                                                                                                                                                                                                   |
| N6.5 | `reverse-charge+mobile`                      | reverse charge - transfer of mobile phones                                                                                                                                                                                                               |
| N6.6 | `reverse-charge+electronics`                 | reverse charge - transfer of electronic products                                                                                                                                                                                                         |
| N6.7 | `reverse-charge+construction`                | reverse charge - provisions in the construction and related sectors                                                                                                                                                                                      |
| N6.8 | `reverse-charge+energy`                      | reverse charge - transactions in the energy sector                                                                                                                                                                                                       |
| N6.9 | `reverse-charge+other`                       | reverse charge - other cases                                                                                                                                                                                                                             |
| N7   | `vat-eu`                                     | "VAT paid in other EU countries (telecommunications, tele-broadcasting and electronic services provision pursuant to Art. 7 -octies letter a, b, art. 74-sexies Italian Presidential Decree 633/72)"                                                     |

##### TipoRitenuta (Withholding Type)

Withholding types are different categories of withheld taxes that can be applied alongside VAT. The following list identifies those currently supported by GOBL:

| Code | Tax Category Code | Description                        |
| ---- | ----------------- | ---------------------------------- |
| RT01 | `IRPEF`           | witholding tax natural persons     |
| RT02 | `IRES`            | witholding corporate entities      |
| RT03 | `INPS`            | INPS contribution                  |
| RT04 | `ENASARCO`        | ENASARCO contribution              |
| RT05 | `ENPAM`           | ENPAM contribution                 |
| RT06 | not supported     | Other social security contribution |

##### TipoCassa (Fund Type)

| Code | Description                                                                                                                                    |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| TC01 | National Pension and Welfare Fund for Lawyers and Solicitors                                                                                   |
| TC02 | Pension fund for accountants                                                                                                                   |
| TC03 | Pension and welfare fund for surveyors                                                                                                         |
| TC04 | National pension and welfare fund for self-employed engineers and architects                                                                   |
| TC05 | National fund for solicitors                                                                                                                   |
| TC06 | National pension and welfare fund for bookkeepers and commercial experts                                                                       |
| TC07 | National welfare board for sales agents and representatives (ENASARCO - Ente Nazionale Assistenza Agenti e Rappresentanti di Commercio)        |
| TC08 | National pension and welfare board for employment consultants (ENPACL - Ente Nazionale Previdenza e Assistenza Consulenti del Lavoro)          |
| TC09 | National pension and welfare board for doctors (ENPAM - Ente Nazionale Previdenza e Assistenza Medici)                                         |
| TC10 | National pension and welfare board for pharmacists (ENPAF - Ente Nazionale Previdenza e Assistenza Farmacisti )                                |
| TC11 | National pension and welfare board for veterinary physicians (ENPAV - Ente Nazionale Previdenza e Assistenza Veterinari)                       |
| TC12 | National pension and welfare board for agricultural employees (ENPAIA - Ente Nazionale Previdenza e Assistenza Impiegati dell'Agricoltura)     |
| TC13 | Pension fund for employees of shipping companies and maritime agencies)                                                                        |
| TC14 | National pension institute for Italian journalists (INPGI - Istituto Nazionale Previdenza Giornalisti Italiani)                                |
| TC15 | National welfare board for orphans of Italian doctors (ONAOSI - Opera Nazionale Assistenza Orfani Sanitari Italiani)                           |
| TC16 | Autonomous supplementary welfare fund for Italian journalists (CASAGIT - Cassa Autonoma Assistenza Integrativa Giornalisti Italiani)           |
| TC17 | Pension board for industrial experts and graduate industrial experts (EPPI - Ente Previdenza Periti Industriali e Periti Industriali Laureati) |
| TC18 | National multi-category pension and welfare board (EPAP - Ente Previdenza e Assistenza Pluricategoriale)                                       |
| TC19 | National pension and welfare board for biologists (ENPAB - Ente Nazionale Previdenza e Assistenza Biologi)                                     |
| TC20 | National pension and welfare board for the nursing profession (ENPAPI - Ente Nazionale Previdenza e Assistenza Professione Infermieristica)    |
| TC21 | National pension and welfare board for psychologists (ENPAP - Ente Nazionale Previdenza e Assistenza Psicologi)                                |
| TC22 | National Social Security Institute (INPS - Istituto Nazionale della Previdenza Sociale)                                                        |

## TODO

- Document Codice Destinatario (uses inbox codes)
- Document how local codes are mapped
