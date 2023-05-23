# IT

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

|      |                                                                                                                                                                                |
| ---- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| RF01 | Ordinary                                                                                                                                                                       |
| RF02 | "Minimum taxpayers (Art. 1, section 96-117, Italian Law 244/07)"                                                                                                               |
| RF04 | "Agriculture and connected activities and fishing (Arts. 34 and 34-bis, Italian Presidential Decree 633/72)"                                                                   |
| RF05 | "Sale of salts and tobaccos (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                          |
| RF06 | "Match sales (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                                         |
| RF07 | "Publishing (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                                                          |
| RF08 | "Management of public telephone services (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                             |
| RF09 | "Resale of public transport and parking documents (Art. 74, section 1, Italian Presidential Decree 633/72)"                                                                    |
| RF10 | "Entertainment, gaming and other activities referred to by the tariff attached to Italian Presidential Decree 640/72 (Art. 74, section 6, Italian Presidential Decree 633/72)" |
| RF11 | "Travel and tourism agencies (Art. 74-ter, Italian Presidential Decree 633/72)"                                                                                                |
| RF12 | "Farmhouse accommodation/restaurants (Art. 5, section 2, Italian law 413/91)"                                                                                                  |
| RF13 | "Door-to-door sales (Art. 25-bis, section 6, Italian Presidential Decree 600/73)"                                                                                              |
| RF14 | "Resale of used goods, artworks, antiques or collector's items (Art. 36, Italian Decree Law 41/95)"                                                                            |
| RF15 | "Artwork, antiques or collector's items auction agencies (Art. 40-bis, Italian Decree Law 41/95)"                                                                              |
| RF16 | "VAT paid in cash by P.A. (Art. 6, section 5, Italian Presidential Decree 633/72)"                                                                                             |
| RF17 | "VAT paid in cash by subjects with business turnover below Euro 200,000 (Art. 7, Italian Decree Law 185/2008)"                                                                 |
| RF18 | Other                                                                                                                                                                          |
| RF19 | "Flat rate (Art. 1, section 54-89, Italian Law 190/2014)"                                                                                                                      |

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

##### ModalitaPagamento (Payment Method)

| Code | Key               | SubKey      | Description                                       |
| ---- | ----------------- | ----------- | ------------------------------------------------- |
| MP01 | `cash`            |             | Cash                                              |
| MP02 | `cheque`          |             | cheque                                            |
| MP03 | `bank-draft`      |             | Banker's draft                                    |
| MP04 | `cash`            | `treasury`  | Cash at Treasury                                  |
| MP05 | `credit-transfer` |             | bank transfer                                     |
| MP06 | NA                |             | money order                                       |
| MP07 | NA                |             | pre-compiled bank payment slip                    |
| MP08 | `card`            |             | payment card                                      |
| MP09 | `direct-debit`    |             | direct debit                                      |
| MP10 | `direct-debit`    | `utilities` | Utilities direct debit (must be rare!)            |
| MP11 | `direct-debit`    | `fast`      | fast direct debit                                 |
| MP12 | NA                |             | collection order                                  |
| MP13 | NA                |             | payment by notice                                 |
| MP14 | NA                |             | tax office quittance                              |
| MP15 | NA                |             | transfer on special accounting accounts           |
| MP16 | NA                |             | order for direct payment from bank account        |
| MP17 | NA                |             | order for direct payment from post office account |
| MP18 | NA                |             | bulletin postal account                           |
| MP19 | `direct-debit`    | `sepa`      | SEPA Direct Debit                                 |
| MP20 | `direct-debit`    | `sepa-core` | SEPA Direct Debit CORE                            |
| MP21 | `direct-debit`    | `sepa-b2b`  | SEPA Direct Debit B2B                             |
| MP22 | `credit`          |             | Deduction on sums already collected               |
| MP23 | `online`          | `pagopa`    | PagoPA                                            |

##### TipoDocumento (Document Type)

| Code | Type          | Scheme / Condition         | Description                                                                                                                                                                       |
| ---- | ------------- | -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| TD01 |               |                            | Regular Invoice                                                                                                                                                                   |
| TD02 | `partial`     |                            | Advance or down payment on invoice                                                                                                                                                |
| TD03 | `partial`     | `freelancer`               | Advance or down payment on freelance invoice                                                                                                                                      |
| TD04 | `credit-note` |                            | Credit note                                                                                                                                                                       |
| TD05 | `debit-note`  |                            | Debit note                                                                                                                                                                        |
| TD06 |               | `freelancer`               | Freelancer Invoice - includes retained taxes                                                                                                                                      |
| TD07 | `simplified`  |                            | Simplified (\*)                                                                                                                                                                   |
| TD08 | `credit-note` | no customer                | Simplified Credit Note (\*)                                                                                                                                                       |
| TD09 | `debit-note`  | no customer                | Simplified Debit Note (\*)                                                                                                                                                        |
| TD16 |               | `reverse-charge`           | reverse charge internal invoice integration                                                                                                                                       |
| TD17 |               | TBD                        | integration/self invoicing for purchase services from abroad                                                                                                                      |
| TD18 |               | `eu-goods`                 | integration for purchase of intra UE goods                                                                                                                                        |
| TD19 |               | TBD                        | integration/self invoicing for purchase of goods ex art.17 c.2 DPR 633/72                                                                                                         |
| TD20 | `self-billed` | TBD                        | self invoicing for regularisation and integration of invoices (ex art.6 c.8 and 9-bis d.lgs 471/97 or art.46 c.5 D.L. 331/93)                                                     |
| TD21 | `self-billed` | `ceiling-exceeded`         | Self invoicing when goods are bought for export without VAT until a certain limit. If limit is exceeded, they must issue an invoice of type TD21. (Autofaturra per splafonamento) |
| TD22 |               | `goods`                    | extractions of goods from VAT Warehouse                                                                                                                                           |
| TD23 |               | `goods-with-tax`           | extractions of goods from VAT Warehouse with payment of VAT                                                                                                                       |
| TD24 |               | `deferred`                 | "deferred invoice ex art.21, c.4, lett. a) DPR 633/72"                                                                                                                            |
| TD25 |               | `deferred`, `third-period` | "deferred invoice ex art.21, c.4, third period lett. b) DPR 633/72"                                                                                                               |
| TD26 |               | `depreciable-assets`       | sale of depreciable assets and for internal transfers (ex art.36 DPR 633/72)                                                                                                      |
| TD27 | `self-billed` | TBD                        | self invoicing for self consumption or for free transfer without recourse                                                                                                         |
| TD28 |               | `san-marino-paper`         | Purchases from San Marino with VAT (paper invoice)                                                                                                                                |

Note: fields marked with (\*) are for simplified invoice documents.

##### Natura (Nature)

|      |                                                                                                                                                                                                                                                          |
| ---- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --- |
| N1   | excluded pursuant to Art. 15, DPR 633/72                                                                                                                                                                                                                 |
| N2   | not subject (this code is no longer permitted to use on invoices emitted from 1 January 2021 )                                                                                                                                                           |
| N2.1 | not subject to VAT under the articles from 7 to 7-septies of DPR 633/72                                                                                                                                                                                  |
| N2.2 | not subject – other cases                                                                                                                                                                                                                                |
| N3   | not taxable (this code is no longer permitted to use on invoices emitted from 1 January 2021 )                                                                                                                                                           |
| N3.1 | not taxable - exportations                                                                                                                                                                                                                               |
| N3.2 | not taxable - intra Community transfers                                                                                                                                                                                                                  |
| N3.3 | not taxable - transfers to San Marino                                                                                                                                                                                                                    |
| N3.4 | not taxable - transactions treated as export supplies                                                                                                                                                                                                    |
| N3.5 | not taxable - for declaration of intent                                                                                                                                                                                                                  |
| N3.6 | not taxable – other transactions that don’t contribute to the determination of ceiling                                                                                                                                                                   |
| N4   | exempt                                                                                                                                                                                                                                                   |
| N5   | margin regime / VAT not exposed on invoice                                                                                                                                                                                                               |
| N6   | "reverse charge (for transactions in reverse charge or for self invoicing for purchase of extra UE services or for import of goods only in the cases provided for) — (this code is no longer permitted to use on invoices emitted from 1 January 2021 )" |     |
| N6.1 | reverse charge - transfer of scrap and of other recyclable materials                                                                                                                                                                                     |
| N6.2 | reverse charge - trasnfer of gold and pure silver pursuant to law 7/2000 as well as used jewelery to OPO                                                                                                                                                 |
| N6.3 | reverse charge - subcontracting in the construction sector                                                                                                                                                                                               |
| N6.4 | reverse charge - transfer of buildings                                                                                                                                                                                                                   |
| N6.5 | reverse charge - transfer of mobile phones                                                                                                                                                                                                               |
| N6.6 | reverse charge - transfer of electronic products                                                                                                                                                                                                         |
| N6.7 | reverse charge - provisions in the construction and related sectors                                                                                                                                                                                      |
| N6.8 | reverse charge - transactions in the energy sector                                                                                                                                                                                                       |
| N6.9 | reverse charge - other cases                                                                                                                                                                                                                             |
| N7   | "VAT paid in other EU countries (telecommunications, tele-broadcasting and electronic services provision pursuant to Art. 7 -octies letter a, b, art. 74-sexies Italian Presidential Decree 633/72)"                                                     |

##### TipoRitenuta (Withholding Type)

| Code | Description                        |
| ---- | ---------------------------------- |
| RT01 | witholding tax natural persons     |
| RT02 | witholding corporate entities      |
| RT03 | INPS contribution                  |
| RT04 | ENASARCO contribution              |
| RT05 | ENPAM contribution                 |
| RT06 | Other social security contribution |

## TODO
- Document Codice Destinatario (uses inbox codes)
- Document how local codes are mapped
