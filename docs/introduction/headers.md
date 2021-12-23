# Headers

The GOBL Header serves to describe the document payload. Headers are considered mutable, they can be added to, but cannot be reduced without risking the signatures being invalidated.

The basic fields included in a header are:

* UUID (`uuid`) - Required. A unique UUIDv1 identifier for the envelope. The document may also include its own IDs. Using a v1 as opposed to a v4 UUID implies a timestamp is always included.
* Type (`typ`) - Required. The type of document payload. See the [Documents](documents.md) section for details on what is currently supported.
* Region (`rgn`) - Required. Fiscal or tax region for the contents. If any local validation rules need to be applied, the region determines which they are.
* Digest (`dig`) - Required. A digest object that defines the encoding mechanism and value of the document's content. More details [below](headers.md#digest).
* Stamps (`stamps`) - Optional. Stamps are meant to be used when an additional localised signature or seal is generated for a document such as when sent to a governmental agency. Provider (`prv`) and value (`val`) fields are included with each.
* Tags (`tags`) - Optional. A simple set of text tags that may be useful for organising envelopes.
* Meta (`meta`) - Optional. A hash of additional meta data that doesn't fit into the existing structures.
* Notes (`notes`) - Optional. Additional text notes.
* Draft (`draft`) - True or false boolean. When true, the document should be considered a draft and not be used for official purposes.

### Digest

WIP.
