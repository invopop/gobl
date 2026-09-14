## Added

- `untdid`: new `untdid-unit` extension for preserving any UN/ECE
  Recommendations 20 or 21 unit code. The EN16931 addon provides
  `UnitToUNTDID` and `UnitFromUNTDID` converters.

## Changed

- **breaking**: removed `org.Unit`; unit fields and constants now use `cbc.Key`,
  consistent with other GOBL key-based vocabularies. `UnitDefinitions`,
  `HasValidUnitKey`, and `ExtendUnitKeySchema` provide shared definitions,
  contextual validation, and schema choices. During normalization, `org.Item`
  moves legacy UN/ECE unit values to the `untdid-unit` extension without
  interpreting them. The `eu-en16931-v2017` addon owns bidirectional mapping
  between that extension and GOBL unit keys.
