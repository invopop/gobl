## Fixed

- `org`: `Inbox` codes use the default `cbc.NormalizeCode` instead of strict
  normalization, so participant IDs such as
  `0225:533304358_53330435800013_DEF0__CODIAL` are kept as issued.
