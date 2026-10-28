## Added

- `bill`: `no-lines` invoice tag for documents that carry their totals but not
  their lines, such as Factur-X BASIC WL. An invoice with the tag no longer
  needs lines, and must also have the `bypass` tag so the totals it states are
  kept.

## Changed

- `eu-en16931-v2017`: invoices with the `no-lines` tag are not required to
  have a line (BR-16).
