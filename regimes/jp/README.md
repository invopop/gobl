# GOBL Japan Tax Regime — `JP`

Japan's tax regime for GOBL covering the Consumption Tax, Withholding Incom  and the Qualified Invoice System.

## Tax Categories

### Consumption Tax (JCT - Japanese Consumption Tax, 消費税, Shōhizei)

- [What is the Japanese Consumption Tax?](https://stripe.com/en-es/resources/more/what-is-japanese-consumption-tax)
- [Consumption Tax](https://www.nta.go.jp/english/taxes/consumption_tax/01.htm#c01)

Japan's consumption tax is equivalent to VAT, applied to goods and services sales.

| Rate | Percentage | Since |
|------|-----------|-------|
| Standard | 10% | 2019-10-01 |
| Standard | 8% | 2014-04-01 |
| Standard | 5% | 1997-04-01 |
| Standard | 3% | 1989-04-01 |
| Reduced | 8% | 2019-10-01 |

The reduced rate of 8% applies to food and beverages (excluding dining out and alcohol) and newspaper subscriptions.

### Withholding Income Tax (源泉徴収, Gensen Chōshū)

- [Japan Withholding Taxes](https://taxsummaries.pwc.com/japan/corporate/withholding-taxes)
- [Withholding Tax Rates - National Tax Agency](https://www.nta.go.jp/publication/pamph/gensen/shikata_r08/pdf/15.pdf)

Japan does have a withholding tax system where businesses must withhold income tax from payments to certain service 
providers (freelancers, lawyers, accountants, etc.). This is similar to Spain's IRPF in GOBL — a retained tax that does appear on invoices.

| Rate | Percentage | Since | Notes |
|------|-----------|-------|-------|
| Professional | 10.21% | 2013-01-01 | Includes 2.1% reconstruction surtax |
| Professional | 10% | 1989-04-01 | Base rate |
| Professional (over ¥1M) | 20.42% | 2013-01-01 | Includes 2.1% reconstruction surtax |
| Professional (over ¥1M) | 20% | 1989-04-01 | Base rate |


## Tax Identity

- [Corporate Number System](https://www.houjin-bangou.nta.go.jp/en/shitsumon/shosai.html?selQaId=00001)
- [Japan Corporate Number](https://en.wikipedia.org/wiki/Corporate_Number)

The Corporate Numbers (Japanese: 法人番号, Hepburn: hōjin bangō) are 13-digit identifiers assigned by the National Tax Agency 
to companies and other organizations registered in Japan. When filing tax returns or other forms related to taxation, 
employment or social insurance, assignees are required to print their own Corporate Number on the document.


- 13 digits, numeric only
- First digit is a check digit (1-9)
- Checksum: `9 - ((Σ Pn × Qn) mod 9)`

Corporate numbers are publicly searchable on the [NTA Corporate Number Publication Site](https://www.houjin-bangou.nta.go.jp/en/).

## Qualified Invoice System (QIS)

- [EU-Japan Centre - Qualified Invoice System](https://www.eu-japan.eu/qualified-invoice-system)
- [Qualified Invoice System Instructions](https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu/pdf/0024006-039_01.pdf)
- 
Since October 1, 2023, Japan's Qualified Invoice System (適格請求書等保存方式) requires registered issuers to include their registration number on invoices for customers to claim input tax credits.
Businesses must deliver and retain a qualified invoice (commonly just called an invoice) that meets the system’s requirements to correctly apply tax credits for purchases and accurately pay consumption tax.

The **Invoice Registration Number** (適格請求書発行事業者登録番号) has the format `T` followed by 13 digits. For corporations, this is `T` + their corporate number.

This is modeled as an organization identity with the key `jp-invoice-registration-number`.


