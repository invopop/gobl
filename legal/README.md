# GOBL Legal Package

> [!IMPORTANT]
> This package is experimental. It provides a canonical container for legal
> documents; it is not a complete legal ontology, legal advice, or a guarantee
> that an agreement is valid or enforceable in any jurisdiction.

The package uses a hybrid model. Authoritative prose remains readable and
largely unconstrained, while a small deterministic spine identifies the
parties, versions, dates, governing context, execution requirements, related
documents, and stable locations in the text. GOBL envelopes can canonicalize
and sign each artifact.

## Architecture

The model separates three artifacts that have different trust and lifecycle
properties:

1. `Contract` contains the authoritative text and essential legal metadata.
   Its UUID identifies one immutable version; `agreement` identifies the
   enduring relationship across amendments and restatements.
2. `Assent` records a party or signatory's intent toward the digest of an exact
   contract version. Keeping assent separate avoids trying to embed a signature
   in the bytes that the same signature must authenticate.
3. `Analysis` contains replaceable human-, rules-, or model-produced
   annotations. It is digest-bound to a contract version but is not part of the
   authoritative bargain.

```text
agreement UUID
    |
    +-- Contract v1 UUID + canonical digest <--- Assent(s)
    |                 ^
    |                 +------------------------- Analysis(es)
    |
    +-- amendment UUID + canonical digest <---- Assent(s)
```

Execution requirements state who is expected to sign and with what intent.
Collected `Assent` documents are evidence that those requirements were met;
the contract does not contain a mutable `signed` or `effective` flag. Systems
should derive lifecycle state from the immutable documents and relevant facts.

## Deliberate boundaries

Chapters, sections, recitals, and definitions may have stable anchors. Local
references and analysis annotations resolve to those anchors. The prose itself
is not forced into a universal schema of obligations, permissions, warranties,
or remedies. Those interpretations belong in `Analysis`, where they can carry
the producing method, be reviewed, replaced, and compared without silently
changing the contract.

External material that is incorporated by reference requires a digest. A URL
alone identifies a location whose contents may change; it is suitable for
supporting material, but not enough to commit to incorporated terms.

The current experiment still leaves substantial work to applications and
jurisdiction-specific extensions, including:

- identity assurance, signing technology, authority, witnesses, and deeds;
- offer, delivery, acceptance, withdrawal, and effective-date event evidence;
- amendment application, consolidated views, and conflict rules;
- localized execution formalities, consumer protections, and required notices;
- confidentiality and access control for contract text and personal data; and
- reliable rendering and preservation of the human presentation that was
  actually reviewed.

Standard terms may describe an open party role, such as any future user,
instead of pretending that the adhering legal person is already known. The
subsequent `Assent` can bind that role to a concrete entity and representative.

The intended role of this package is therefore a stable interoperability and
evidence layer, not a replacement for prose, legal judgment, or contract
management.
