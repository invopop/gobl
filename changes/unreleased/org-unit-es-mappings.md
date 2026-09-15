## Added

- `eu-en16931-v2017`: the `portion` unit maps to the UNTDID `13` (ration) code,
  so it can satisfy BR-23, which requires a unit code on every invoice line.
  Every unit GOBL defines now has an exact UNTDID equivalent.

## Removed

- **breaking**: `org`: the non-standard `6pack` and `tetrabrik` units, which had
  no UN/ECE equivalent and so could never satisfy BR-23. Normalization replaces
  them with `pkg` and `carton` respectively, so existing documents keep
  validating.
