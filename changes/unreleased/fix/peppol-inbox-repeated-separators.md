## Fixed

- `org`: inbox codes with a scheme or the `peppol` key are no longer stripped
  of repeated separators during normalization, so participant IDs such as
  `0225:533304358_53330435800013_DEF0__CODIAL` keep their `__` in both the
  inbox and the derived Peppol endpoint.
