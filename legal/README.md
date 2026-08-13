# GOBL Legal Package

> [!IMPORTANT]
> This package is experimental. It adapts the body of a legal document, and
> specifically a contract, to GOBL's canonical structured format. It is not a
> complete legal-agreement ontology and using it does not make a contract valid
> or enforceable.

The current model provides a deterministic hierarchy of chapters and nested
sections. Chapters and sections may have unique anchors and URI references;
local fragment references are checked against anchors in the same contract.
Normalization cleans text, removes nil list entries, and derives one-based
chapter and section indexes. GOBL envelopes can then canonicalize and sign the
resulting document body.

The model deliberately does not yet define parties and their roles, signer
authority or capacity, offer and acceptance, effective and execution dates,
governing law, conditions, schedules and exhibits, amendment/supersession
relationships, or jurisdiction-specific execution requirements. Those concepts
need explicit semantics and lifecycle rules before this package should be
presented as a general contract format.

An external `$ref` currently commits only to the reference URI. It does not
commit to the bytes found at that URI, so mutable external terms are not safely
incorporated by reference. A future reference model should carry a digest and,
where practical, an immutable snapshot or attachment.
