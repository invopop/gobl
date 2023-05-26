# PT

Portugal doesn't have an e-invoicing format per se. Tax information is reported electronically to the AT (Autoridade Tributária e Aduaneira) either periodically in batches via a SAF-T (PT) report or individually in real time via a web service.

## Public Documentation

* [Portaria n.o 302/2016 – SAF-T Data Structure & Taxonomies](https://info.portaldasfinancas.gov.pt/pt/informacao_fiscal/legislacao/diplomas_legislativos/Documents/Portaria_302_2016.pdf)
* [Portaria n.o 195/2020 – Comunicação de Séries Documentais, Aspetos Específicos](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Comunicacao_Series_ATCUD/Documents/Comunicacao_de_Series_Documentais_Manual_de_Integracao_de_SW_Aspetos_Genericos.pdf)
* [Portaria n.o 195/2020 – Especificações Técnicas Código QR](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Novas_regras_faturacao/Documents/Especificacoes_Tecnicas_Codigo_QR.pdf)
* [Comunicação dos elementos dos documentos de faturação à AT, por webservice](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Fatcorews/Documents/Comunicacao_dos_elementos_dos_documentos_de_faturacao.pdf)

## Local Codes

### `InvoiceType` (Tipo de documento)

| Code | Name | GOBL Type | GOBL Tax Tag |
| --- | --- | --- | --- |
| FT | Fatura, emitida nos termos do artigo 36.o do Código do IVA | `standard` | |
| FS | Fatura simplificada, emitida nos termos do artigo 40.o do Código do IVA | `standard` | `simplified` |
| FR | Fatura-recibo | `standard` | `invoice-receipt` |
| ND | Nota de débito | `credit-note` | |
| NC | Nota de crédito | `debit-note` | |

### `TaxExemptionCode` (Código do motivo de isenção de imposto)

| Code | Description | GOBL Line Tax (VAT) Tag |
| --- | --- | --- |
| M01 | Artigo 16.°, n.° 6 do CIVA | `exempt+outlay` |
| M02 | Artigo 6.° do Decreto-Lei n.° 198/90, de 19 de junho | `exempt+intrastate-export` |
| M04 | Isento artigo 13.° do CIVA | `exempt+imports` |
| M05 | Isento artigo 14.° do CIVA | `exempt+exports` |
| M06 | Isento artigo 15.° do CIVA | `exempt+suspension-scheme` |
| M07 | Isento artigo 9.° do CIVA | `exempt+internal-operations` |
| M09 | IVA - não confere direito a dedução / Artigo 62.° alínea b) do CIVA | `exempt+small-retail-scheme` |
| M10 | IVA - regime de isenção / Artigo 57.° do CIVA | `exempt+exempt-scheme` |
| M11 | Regime particular do tabaco / Decreto-Lei n.° 346/85, de 23 de agosto | `exempt+tobacco-scheme` |
| M12 | Regime da margem de lucro - Agências de viagens / Decreto-Lei n.° 221/85, de 3 de julho | `exempt+margin-scheme+travel` |
| M13 | Regime da margem de lucro - Bens em segunda mão / Decreto-Lei n.° 199/96, de 18 de outubro | `exempt+margin-scheme+second-hand` |
| M14 | Regime da margem de lucro - Objetos de arte / Decreto-Lei n.° 199/96, de 18 de outubro | `exempt+margin-scheme+art` |
| M15 | Regime da margem de lucro - Objetos de coleção e antiguidades / Decreto-Lei n.° 199/96, de 18 de outubro | `exempt+margin-scheme+antiques` |
| M16 | Isento artigo 14.° do RITI | `exempt+goods-transmission` |
| M19 | Outras isenções - Isenções temporárias determinadas em diploma próprio | `exempt+other` |
| M20 | IVA - regime forfetário / Artigo 59.°-D n.°2 do CIVA | `exempt+flat-rate-scheme` |
| M21 | IVA - não confere direito à dedução (ou expressão similar) - Artigo 72.° n.° 4 do CIVA | `exempt+non-deductible` |
| M25 | Mercadorias à consignação - Artigo 38.° n.° 1 alínea a) do CIVA | `exempt+consignment-goods` |
| M30 | IVA - autoliquidação / Artigo 2.° n.° 1 alínea i) do CIVA | `exempt+reverse-charge+waste` |
| M31 | IVA - autoliquidação / Artigo 2.° n.° 1 alínea j) do CIVA | `exempt+reverse-charge+civil-eng` |
| M32 | IVA - autoliquidação / Artigo 2.° n.° 1 alínea I) do CIVA | `exempt+reverse-charge+greenhouse` |
| M33 | IVA - autoliquidação / Artigo 2.° n.° 1 alínea m) do CIVA | `exempt+reverse-charge+woods` |
| M40 | IVA - autoliquidação / Artigo 6.° n.° 6 alínea a) do CIVA, a contrário | `exempt+reverse-charge+b2b` |
| M41 | IVA - autoliquidação / Artigo 8.° n.° 3 do RITI | `exempt+reverse-charge+intraeu` |
| M42 | IVA - autoliquidação / Decreto-Lei n.° 21/2007, de 29 de janeiro | `exempt+reverse-charge+real-estate` |
| M43 | IVA - autoliquidação / Decreto-Lei n.° 362/99, de 16 de setembro | `exempt+reverse-charge+gold` |
| M99 | Não sujeito ou não tributado | `exempt+non-taxable` |