## Changed

- `es-facturae-v3`: **breaking** the FacturaE addon moved to its own module,
  [gobl.es.facturae](https://github.com/invopop/gobl.es.facturae), which must
  be imported (`github.com/invopop/gobl.es.facturae/addon`) to calculate or
  validate documents that declare it. The key stays valid in the `$addons`
  JSON Schema as an approved external addon. Fault codes are unchanged.

## Removed

- `addons/es/facturae`: moved to gobl.es.facturae. The Spanish examples no
  longer declare the addon.
