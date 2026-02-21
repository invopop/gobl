# New Zealand (NZ)

This regime models New Zealand’s Goods and Services Tax (GST) using a standard rate and a zero rate, plus basic IRD number format validation. It follows the same structure as other GST regimes in this repo while keeping rules minimal and well-sourced.

Find example NZ GOBL files in the `examples/nz` (uncalculated documents) and `examples/nz/out` (calculated envelopes) subdirectories.

## Rates

- **Standard rate (15%)** for most taxable supplies.
- **Zero-rated (0%)** applies to specific supplies such as exports and other defined cases.
- **Exempt** supplies are not charged GST, but do not allow input tax credits.

## IRD number format

IRD numbers are 8 or 9 digits. 8-digit numbers may be expressed with a leading zero for display. This regime validates format only (no checksum), consistent with other regimes that lack a public checksum.

Sources:
- https://www.ird.govt.nz/gst/charging-gst
- https://www.ird.govt.nz/en/gst/charging-gst/zero-rated-supplies
- https://www.ird.govt.nz/en/gst/charging-gst/exempt-supplies
- https://www.ird.govt.nz/gst/registering-for-gst
- https://www.ird.govt.nz/myir-help/logging-in/ird-numbers
