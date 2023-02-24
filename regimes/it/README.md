# IT

Italy uses the FatturaPA format for their e-invoicing system.

## Public Documentations

### FatturaPA

[Historical Documentations](https://www.fatturapa.gov.it/en/norme-e-regole/documentazione-fattura-elettronica/formato-fatturapa/)

[Schema V1.2.1 Spec Table View (EN)](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.1/Table-view-B2B-Ordinary-invoice.pdf) - by far the most comprehensible spec doc. Since the difference between 1.2.2 and 1.2.1 is minimal, this is perfectly usable.

[Schema V1.2.2 PDF (IT)](https://www.fatturapa.gov.it/export/documenti/Specifiche_tecniche_del_formato_FatturaPA_v1.3.1.pdf) - most up-to-date but difficult

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

### Challenges

#### Special Codes (WIP)

FatturaPA demands the classificationss of the invoice data using predefined
alphanmueric code refered to as FPACodes in this package. These codes can be
unbelieveably specific (e.g. TD28: "Purchase from San Marino with VAT (paper
invoice)", RF06: "Match sales") and includes distinctions we normally would not
(MP09, MP10, MP11, MP19, MP20, and MP21 all refer to different types of direct
debit payments).

Additionally, there is no straightforward mapping between the FPACodes and
fields supported in `bill.Invoice`. "Nature" codes, for example, include things
like reverse charges (what we would find in a `inv.Tax.Schemes`) as well as
classifications for "non-taxable" items (not really a scheme—more like a `Note`
attached to a line item?).

##### Tax System (Regime Fiscale)

A "tax system" in Italy is a property of the seller and not the product or the
service provided.

<b>bill.Invoice Mapping:</b> none (or could we use `inv.Tax.Schemes`?)

| Code | support | Description                                                    |
|------|---------|----------------------------------------------------------------|
| RF01 | ✓       | Ordinary                                                       |
| RF02 |         | Minimum taxpayers                                              |
| RF04 |         | Agriculture and connected activities and fishing               |
| RF05 |         | Sale of salts and tobaccos                                     |
| RF06 |         | Match sales                                                    |
| RF07 |         | Publishing                                                     |
| RF08 |         | Management of public telephone services                        |
| RF09 |         | Resale of public transport and parking documents               |
| RF10 |         | Entertainment, gaming, etc referred to by tariff in DPR 640/72 |
| RF11 |         | Travel and tourism agencies                                    |
| RF12 |         | Agrotourism (Farmhouse accomodations and restaurants)          |
| RF13 |         | Door-to-door sales                                             |
| RF14 |         | Resale of used goods, artworks, antiques or collector's items  |
| RF15 |         | Artwork, antiques or collector's items auction agencies        |
| RF16 |         | VAT paid in cash by P.A.                                       |
| RF17 |         | VAT paid in cash by subjects with business turnover <€200,000  |
| RF18 |         | Other                                                          |
| RF19 |         | Flat rate                                                      |

##### Fund Type (TipoCassa)

<b>bill.Invoice Mapping:</b> none

| Code | Supported | Description                                                                  |
|------|-----------|------------------------------------------------------------------------------|
| TC01 |           | National Pension and Welfare Fund for Lawyers and Solicitors                 |
| TC02 |           | Pension fund for accountants                                                 |
| TC03 |           | Pension and welfare fund for surveyors                                       |
| TC04 |           | National pension and welfare fund for self-employed engineers and architects |
| TC05 |           | National fund for solicitors                                                 |
| TC06 |           | National pension and welfare fund for bookkeepers and commercial experts     |
| TC07 |           | ENASARCO (National welfare board for sales agents and representatives)       |
| TC08 |           | ENPACL (National pension and welfare board for employment consultants)       |
| TC09 |           | ENPAM (National pension and welfare board for doctors)                       |
| TC10 |           | ENPAF (National pension and welfare board for pharmacists)                   |
| TC11 |           | ENPAV (National pension and welfare board for veterinary physicians)         |
| TC12 |           | ENPAIA (National pension and welfare board for agricultural employees)       |
| TC13 |           | Pension fund for employees of shipping companies and maritime agencies)      |
| TC14 |           | INPGI (National pension institute for Italian journalists)                   |
| TC15 |           | ONAOSI (National welfare board for orphans of Italian doctors)               |
| TC16 |           | CASAGIT (Autonomous supplementary welfare fund for Italian journalists)      |
| TC17 |           | EPPI (Pension board for industrial experts and graduate industrial experts)  |
| TC18 |           | EPAP (National multi-category pension and welfare board)                     |
| TC19 |           | ENPAB (National pension and welfare board for biologists)                    |
| TC20 |           | ENPAPI (National pension and welfare board for the nursing profession)       |
| TC21 |           | ENPAP (National pension and welfare board for psychologists)                 |
| TC22 |           | INPS (National Social Security Institute)                                    |

##### Payment Method (Modalita di Pagamento)

<b>bill.Invoice Mapping:</b> inv.Payment.Instructions

| Code | Support | Description                                       |
|------|---------|---------------------------------------------------|
| MP01 | ✓       | cash                                              |
| MP02 |         | cheque                                            |
| MP03 |         | banker's draft                                    |
| MP04 |         | cash at Treasury                                  |
| MP05 | ✓       | bank transfer                                     |
| MP06 |         | money order                                       |
| MP07 |         | pre-compiled bank payment slip                    |
| MP08 | ✓       | payment card                                      |
| MP09 | ✓       | direct debit                                      |
| MP10 | ✓       | utilities direct debit                            |
| MP11 | ✓       | fast direct debit                                 |
| MP12 |         | collection order                                  |
| MP13 |         | payment by notice                                 |
| MP14 |         | tax office quittance                              |
| MP15 |         | transfer on special accounting accounts           |
| MP16 |         | order for direct payment from bank account        |
| MP17 |         | order for direct payment from post office account |
| MP18 |         | bulletin postal account                           |
| MP19 | ✓       | SEPA Direct Debit                                 |
| MP20 | ✓       | SEPA Direct Debit CORE                            |
| MP21 | ✓       | SEPA Direct Debit B2B                             |
| MP22 |         | Deduction on sums already collected               |
| MP23 |         | PagoPA                                            |

##### Document Type (Tipo Documento)

<b>bill.Invoice Mapping:</b> inv.Type

| Code | Support | Description                                                          |
|------|---------|----------------------------------------------------------------------|
| TD01 | ✓       | invoice                                                              |
| TD02 |         | advance/down payment on invoice                                      |
| TD03 |         | advance/down payment on fee                                          |
| TD04 | ✓       | credit note                                                          |
| TD05 |         | debit note                                                           |
| TD06 |         | fee                                                                  |
| TD16 |         | reverse charge internal invoice integration                          |
| TD17 |         | integration/self invoicing for purchase services from abroad         |
| TD18 |         | integration for purchase of intra UE goods                           |
| TD19 |         | integration/self invoicing for purchase of goods                     |
| TD20 |         | self invoicing for regularisation and integration of invoices        |
| TD21 |         | self invoicing for splaphoning                                       |
| TD22 |         | extractions of goods from VAT Warehouse                              |
| TD23 |         | extractions of goods from VAT Warehouse with payment of VAT          |
| TD24 |         | "deferred invoice ex art.21, c.4, lett. a) DPR 633/72"               |
| TD25 |         | "deferred invoice ex art.21, c.4, third period lett. b) DPR 633/72"  |
| TD26 |         | sale of depreciable assets and for internal transfers                |
| TD27 |         | self invoicing for self consumption / free transfer without recourse |
| TD28 |         | Purchases from San Marino with VAT (paper invoice)                   |

##### Nature (Natura)

<b>bill.Invoice Mapping:</b> inv.Tax.Schemes, (potentially) inv.Lines[].Notes

| Code | Support | Description                                                             |
|------|---------|-------------------------------------------------------------------------|
| N1   |         | excluded pursuant to Art. 15, DPR 633/72                                |
| N2.1 |         | not subject to VAT under the articles from 7 to 7-septies of DPR 633/72 |
| N2.2 |         | not subject – other cases                                               |
| N3.1 |         | not taxable - exportations                                              |
| N3.2 |         | not taxable - intra Community transfers                                 |
| N3.3 |         | not taxable - transfers to San Marino                                   |
| N3.4 |         | not taxable - transactions treated as export supplies                   |
| N3.5 |         | not taxable - for declaration of intent                                 |
| N3.6 |         | not taxable – other transactions that do not count towards the plafond. |
| N4   |         | exempt                                                                  |
| N5   |         | margin regime / VAT not exposed on invoice                              |
| N6.1 | ✓       | reverse charge - transfer of scrap and of other recyclable materials    |
| N6.2 | ✓       | reverse charge - transfer of gold, silver, and jewelry                  |
| N6.3 | ✓       | reverse charge - subcontracting in the construction sector              |
| N6.4 | ✓       | reverse charge - transfer of buildings                                  |
| N6.5 | ✓       | reverse charge - transfer of mobile phones                              |
| N6.6 | ✓       | reverse charge - transfer of electronic products                        |
| N6.7 | ✓       | reverse charge - provisions in the construction and related sectors     |
| N6.8 | ✓       | reverse charge - transactions in the energy sector                      |
| N6.9 | ✓       | reverse charge - other cases                                            |
| N7   |         | VAT paid in other EU countries (telecommunications)                     |

##### Withholding Type (Tipo Ritenuta)

<b>bill.Invoice Mapping:</b> inv.Totals.Taxes

| Code | Support | Description                        |
|------|---------|------------------------------------|
| RT01 | ✓       | witholding tax natural persons     |
| RT02 | ✓       | witholding corporate entities      |
| RT03 | ✓       | INPS contribution                  |
| RT04 | ✓       | ENASARCO contribution              |
| RT05 | ✓       | ENPAM contribution                 |
| RT06 | ✓       | Other social security contribution |
