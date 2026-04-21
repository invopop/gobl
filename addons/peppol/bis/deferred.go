package bis

// The rules below are defined by the Peppol BIS 3.0 schematron but are not
// enforced by this addon. They are satisfied more naturally at serialization
// time by gobl.ubl, which either owns the UBL structure being checked or can
// synthesize the required text from structured GOBL data.
//
// When the corresponding work lands in gobl.ubl, remove the matching comment
// from the country file (e.g. "// SE-R-005 … deferred to gobl.ubl — see
// deferred.go") and, if appropriate, file a follow-up issue noting the move.
//
// Deferred rules:
//
//   - DE-R-022 (attachment filenames must be unique, case-insensitive)
//     GOBL has no single attachment list analogous to
//     cac:AdditionalDocumentReference. Uniqueness is a property of the UBL
//     emitted by gobl.ubl and belongs there.
//
//   - DE-R-018 (early payment discount note must follow the #SKONTO# format)
//     Early-payment terms are already structured in GOBL as
//     bill.Payment.Terms.DueDates. gobl.ubl should render the #SKONTO# note
//     from that field rather than having callers hand-write it.
//
//   - IS-R-008 / IS-R-009 / IS-R-010 (EINDAGI date note, format + due-date
//     presence + comparison with BT-9)
//     GOBL already models the due date structurally via
//     bill.Payment.Terms.DueDates. gobl.ubl should emit the EINDAGI
//     AdditionalDocumentReference from that data; when driven from the same
//     source, the format, presence, and comparison rules become
//     structurally impossible to violate.
//
//   - GR-R-001-2 through GR-R-001-7 (Greek invoice ID segment contents:
//     supplier TIN, YYYYMMDD, sequence, document type, and two further
//     segments)
//     The 6-segment shape (GR-R-001-1) IS enforced here: we join Series +
//     "_" + Code and assert that the split yields six tokens. The contents
//     of each segment are not checked, because the supplier TIN, issue date,
//     and sequence are already on the GOBL document as structured fields and
//     asking callers to additionally encode them into `Code` is duplicative.
//     gobl.ubl should build the Peppol-visible ID from those fields at
//     serialization time; GR-specific callers then only need to supply the
//     last two segments (likely via a future GR-specific field or pair of
//     extensions — TBD during the gobl.ubl work).
