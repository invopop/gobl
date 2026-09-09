# Legal examples

These files demonstrate the three-artifact legal architecture using the
[Invopop Terms & Conditions](https://www.invopop.com/legal-documents#terms-and-conditions)
as source material.

- `invopop-terms.yaml` is an abridged, non-authoritative structural
  representation of the general terms. It intentionally paraphrases the live
  page and excludes the separate API terms and privacy materials.
- `invopop-terms-assent.yaml` shows a hypothetical customer filling the open
  `user` party role and accepting the exact contract digest.
- `invopop-terms-analysis.yaml` shows replaceable semantic annotations that
  point to stable anchors without modifying the contract.

The generated envelopes in `out/` are executable fixtures. If the contract
changes, its calculated envelope digest must also be updated in the assent and
analysis examples. These examples are for format development and are not legal
documents or legal advice.

The `saas-subscription-*` files form a second, wholly fictional example of a
negotiated B2B SaaS subscription. The contract includes pricing, invoicing,
usage allowance, an availability SLA, service credits, support targets, data
protection, liability, renewal, and termination. Two assent documents represent
counterpart execution. The analysis extracts selected commercial and SLA facts
as machine-usable values while leaving the signed prose authoritative.
