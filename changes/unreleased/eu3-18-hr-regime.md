## Added

- `hr`: added the Croatia tax regime, covering the 25% general rate and the 13% and
  5% reduced rates of Article 38 of the VAT Act (NN 73/2013), plus normalization and
  ISO 7064 MOD 11,10 validation of OIB tax codes and `OIB` identities. Rates are
  dated from the Act's commencement on 1 July 2013, including the reduced rate's rise
  from 10% to 13% on 1 January 2014.
- `hr`: added the `cash-basis` invoice tag for suppliers that apply the cash accounting
  scheme of Article 125.i of the VAT Act. It adds the required "Obračun prema
  naplaćenim naknadama" VAT note.
- `hr`: added the `PP` tax category for the consumption tax (porez na potrošnju), the
  local tax on drinks served in catering establishments under the Local Taxes Act. Each
  municipality sets its rate, up to 3%, so the category has no fixed rates and each tax
  line gives its percent. The code is the one the Croatian eRačun specification uses.
- `hr`: added the `hr-pass-through` extension for pass-through items, amounts charged
  together with taxable items that are left out of the VAT base: the packaging return
  fee (`POVNAK`), the special tax on motor vehicles (`PPMV`) and any other non-taxable
  fee (`N`), with the HR tax category codes of the Croatian eRačun specification
  (HR-TB-2) without the `HR:` prefix. The amount is a line or charge of its own, and
  the extension goes on its VAT combo, which must use the exempt key. The consumption
  tax, also a pass-through item, keeps its own `PP` tax category.
