# Envelopes

"GOBL Envelope" is the name we've given to the structure that wraps around a document, adding headers and seals (signatures), just like an envelope in real-life. There are three key parts to an envelope in GOBL:

* Header (`head`) - Meta data that describes the included document.
* Document (`doc`) - The actual payload of the envelope, like an Invoice.
* Signatures (`sigs`) -  A set of [JSON Web Signatures](https://en.wikipedia.org/wiki/JSON\_Web\_Signature) that can be used to verify the headers included in the signature, and thus de document, have not been modified.

We'll go through each of this in a bit more detail in the following chapters.
