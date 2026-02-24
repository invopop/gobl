# 🇦🇩 Andorra Tax Regime (AD)

Implementation of the tax regime for Andorra. The main indirect tax in Andorra is the **Impost General Indirecte (IGI)**, and it is enforced since 1st of January 2013.

## Tax Rates

These are the current IGI tax rates:

| Rate | Key | Percent | Description |
| ---- | --- | ------- | ----------- |
| General | `general` | 4.5% | Standard rate for most goods and services. |
| Reduced | `reduced` | 1.0% | Food, water, books, newspapers. |
| Super-Reduced | `super-reduced` | 0.0% | Health, education, social services. |
| Special | `special` | 2.5% | Transport, libraries, museums. |
| Increased | `increased` | 9.5% | Banking and financial services. |


## Date requirements

- Quarterly: Companies with a turnover of more than €250,000 (April, July, October, January).
- Semestral: Companies with a turnover of less than €250,000 (July, January).

Start of activity: Generally declared semestrally (July and January), unless the special regime applies.

## Tax Identity (NRT)

The **Número de Registre Tributari (NRT)** is the tax identification number in Andorra.

Format: `X-999999-X`
- A leading letter (identifying the type of person/entity):
  - F: Individual Residents
  - E: Non-resident Individuals
  - L: Limited Liability Companies (S.L.)
  - A: Joint-stock Corporations (S.A.)
- Six digits.
- A trailing control letter.

## References

- [Departament de Tributs i de Fronteres - Andorra](https://www.impostos.ad)
- [Andorra NRT number guide](https://lookuptax.com/docs/tax-identification-number/andorra-tax-id-guide)
