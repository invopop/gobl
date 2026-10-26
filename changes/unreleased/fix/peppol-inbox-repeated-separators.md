## Changed

- `cbc`: strict codes now accept runs of `_` between blocks, and
  `NormalizeStrictCode` keeps them instead of collapsing them to a single `_`.
  Other repeated separators are still collapsed.

## Fixed

- `org`: inbox codes such as the Peppol participant ID
  `0225:533304358_53330435800013_DEF0__CODIAL` keep their `__` in both the
  inbox and the derived Peppol endpoint.
