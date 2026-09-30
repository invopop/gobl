## Changed

- `bill`: due date percentages are calculated from the amount due after advances, falling back to the payable amount.
- **breaking**: `pay`: `Record.Ref` and `DirectDebit.Ref` are `cbc.Code` instead of `string`, matching `Instructions.Ref`. JSON documents are unchanged, but references are limited to 128 characters and Go code assigning a `string` needs a conversion.
