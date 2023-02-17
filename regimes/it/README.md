# IT

Italy uses the FatturaPA format for their e-invoicing system.

## Public Documentations

### <b>FatturaPA</b>

[Historical Documentations](https://www.fatturapa.gov.it/en/norme-e-regole/documentazione-fattura-elettronica/formato-fatturapa/)

[Schema V1.2.1 Spec Table View (EN)](https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.1/Table-view-B2B-Ordinary-invoice.pdf) - by far the most comprehensible spec doc. Since the difference between 1.2.2 and 1.2.1 is minimal, this is perfectly usable.

[Schema V1.2.2 PDF (IT)](https://www.fatturapa.gov.it/export/documenti/Specifiche_tecniche_del_formato_FatturaPA_v1.3.1.pdf) - most up-to-date but difficult

#### Changes from 1.2.1 to 1.2.2
- Documentation changes: TD25, N1, N6.2, N7
- Addition of TD28: Acquisti da San Marino con IVA (fattura cartacea)

### <b>Tax Definitions</b>

[Fiscal Code (Codice Fiscale)](https://en.wikipedia.org/wiki/Italian_fiscal_code)

[VAT Number (Partita IVA)](https://en.wikipedia.org/wiki/VAT_identification_number)

[Agenzia Entrate (Tax Office) IVA Doc](https://www.agenziaentrate.gov.it/portale/web/english/nse/business/vat-in-italy)

### Challenges

#### Special Codes (WIP)

FatturaPA demands very specific categorization for the type of economic activity,
document type, fund type, etc. It will be a challenge to map these onto GOBL
constructs.

##### RegimeFiscale (Tax System)
...
##### TipoCassa (Fund Type)
...
##### ModalitaPagamento (Payment Method)
...
##### TipoDocumento (Document Type)
...
##### Natura (Nature)
...
##### TipoRitenuta (Withholding Type)
...
