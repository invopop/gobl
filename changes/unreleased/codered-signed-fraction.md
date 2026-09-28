## Fixed

- `num`: `AmountFromString` rejects a sign in the fraction. `1.-25` was read as `0.975`. `-245.12` still parses.
