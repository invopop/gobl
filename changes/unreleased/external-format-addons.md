## Changed

- `es-facturae-v3`: **breaking** the FacturaE addon moved to its own module,
  [gobl.es.facturae](https://github.com/invopop/gobl.es.facturae), which must
  be imported (`github.com/invopop/gobl.es.facturae/addon`) to calculate or
  validate documents that declare it. The key stays valid in the `$addons`
  JSON Schema as an approved external addon. Fault codes are unchanged.
- `de-xrechnung-v3`, `de-zugferd-v2`: **breaking** the German addons moved to
  [gobl.de.xinvoice](https://github.com/invopop/gobl.de.xinvoice), imported
  from `github.com/invopop/gobl.de.xinvoice/addon/xrechnung` and
  `github.com/invopop/gobl.de.xinvoice/addon/zugferd`. The keys stay valid as
  approved external addons. Fault codes are unchanged.
- `fr-facturx-v1`: **breaking** the Factur-X addon moved to
  [gobl.fr.ctc](https://github.com/invopop/gobl.fr.ctc), imported from
  `github.com/invopop/gobl.fr.ctc/addon/facturx`. The key stays valid as an
  approved external addon.

## Removed

- `addons/es/facturae`: moved to gobl.es.facturae. The Spanish examples no
  longer declare the addon.
- `addons/de/xrechnung`, `addons/de/zugferd`: moved to gobl.de.xinvoice.
- `addons/fr/facturx`: moved to gobl.fr.ctc.
- Examples `invoice-de-de-xr`, `invoice-de-es-xr` and `invoice-fr-fr-facturx`
  are now `invoice-de-de-en16931`, `invoice-de-es-en16931` and
  `invoice-fr-fr-en16931`, declaring only `eu-en16931-v2017`.
