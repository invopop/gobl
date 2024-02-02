# 🇪🇸 GOBL Spain Tax Regime

Example ES GOBL files can be found in the [`examples`](./examples) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Corrective Invoices

According to Spanish law on invoicing [Real Decreto 1619/2012, de 30 de noviembre](https://www.boe.es/buscar/act.php?id=BOE-A-2012-14696), only "rectified" invoices are recognized. There are in fact no mentions of internationally credit or debit notes at all.

To figure out how to map GOBL correction types to Spanish requirements, we've used the FacturaE format which defines a list of Correction Methods. The following assumptions are made on the key corrective invoice types:

- `corrective` - the previous invoice has been completely replaced and for accounting purposes can be discarded. This maps to code `01` (Rectificación modelo integro) in FacturaE.
- `credit-note` - some or potentially all of the line items in the previous invoice have been cancelled. In GOBL the quantities of line items will be positive, but for presentation in Spain should be inverted to reflect negative values. This maps to `02` in FacturaE (Rectificación modelo por diferencias.)
- `debit-note` - effectively the same as for `credit-note`, with the exception that values do not need to be inverted. Debit notes are very rarely used, its often easier to just issue a new invoice.

The FacturaE format defines two other correction methods which are currently considered out of scope for GOBL.

Previous versions of GOBL attempted to automatically invert corrective invoices instead of issuing credit notes, however, this was deemed to be confusing and not compatible with internationally recognized credit notes. The [GOBL to FacturaE convertor](https://github.com/invopop/gobl.facturae) will automatically invert quantities in a Credit Note.
