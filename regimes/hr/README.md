# 🇭🇷 GOBL Croatia Tax Regime

Croatia joined the European Union in 2013 and adopted the **euro (EUR)** as its official currency on 1 January 2023, replacing the Croatian kuna (HRK). Tax administration is managed by the **Porezna uprava** (Croatian Tax Administration).

Find example HR GOBL files in the [`examples`](../../examples/hr) and [`examples/out`](../../examples/hr/out) subdirectories.

## Public Documentation

* [Croatian Tax Administration – VAT](https://porezna-uprava.gov.hr/en/vat/7362)
* [OIB – Personal Identification Number (Official Gazette 60/2008)](https://narodne-novine.nn.hr/clanci/sluzbeni/2008_05_60_2033.html)
* [OIB Check Digit Algorithm (Official Gazette 1/2009)](https://narodne-novine.nn.hr/clanci/sluzbeni/2009_01_1_6.html)

## Tax Identity (OIB)

Croatia assigns every legal and natural person a unique 11-digit **Osobni identifikacijski broj** (OIB, Personal Identification Number). The OIB is the primary tax identifier used in all business and tax documents.

### Format

* **11 numeric digits** — 10 random digits followed by 1 check digit.
* The check digit is computed using the **ISO 7064 MOD 11.10** algorithm.

## Value Added Tax (VAT / PDV)

Croatian VAT is called **PDV** (_Porez na dodanu vrijednost_) and is regulated at the national level by the Croatian Tax Administration. Croatia applies four VAT rates.

| Rate | GOBL Rate | Percent | Since |
| ---- | --------- | ------- | ----- |
| General (_Opća stopa_) | `standard` | 25% | 2013-07-17 |
| Reduced (_Snižena stopa_) | `reduced` | 13% | 2013-07-17 |
| Super Reduced (_Super snižena stopa_) | `super-reduced` | 5% | 2013-07-17 |
| Zero (_Nulta stopa_) | `zero` | 0% | 2013-07-17 |

The **25% general rate** applies to most goods and services. The **13% reduced rate** covers items such as accommodation services, newspapers, and certain food items. The **5% super-reduced rate** applies to essential goods such as basic foodstuffs, medicines, books, and scientific journals.
