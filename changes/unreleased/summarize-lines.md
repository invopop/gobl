## Added

- `bill`: `SummarizeLines` replaces an invoice's lines with one summary line
  for each group of lines with the same taxes, extensions, seller, and item
  key, keeping the totals, for documents with too many lines to be handled in
  full. Summarized invoices are tagged with `summarized` and given a note with
  the number of lines replaced.
- `tax`: `summarized` tag, available on all invoices.
