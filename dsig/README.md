# GoBL Digital Signatures

## Introduction

Digital signatures are one of the fundamental features of GoBL as they bring the ability to be able to mathematically confirm using a public key that the person who owns the private key really did create the document.

This `dsig` package aims to bring together the functionality required to handle GoBL document digests and signatures in one place so they are easy and convenient to use.

Signatures in GoBL use the [Javascript Object Signing and Encryption (JOSE)](https://datatracker.ietf.org/wg/jose/about/) standards specifically around [JSON Web Signatures (JWS) (RFC7515)](https://datatracker.ietf.org/doc/html/rfc7515) and [JSON Web Keys (JWK) (RFC7517)](https://datatracker.ietf.org/doc/html/rfc7517).

Behind the scenes, GoBL uses the [go-jose](https://github.com/go-jose/go-jose) library to do all the heavy lifting and provides wrappers that make it easy to use sensible defaults. There should not be anything that cannot be implemented in another language, but helpers do make life easier and limit what is available to the use-cases of GoBL documents.

There are four key components to the dsig implementation:

 * **Private Key** - Private JSON Web Keys (JWK), that can be used to create signatures. Currently, GoBL only supports ECDSA keys using a 256-bit curve. The private key is used to create a public counterpart and in addition to the JWK standards, every key *must* be identified with a UUID.
 * **Public Key** -  Public JSON Web Keys used to verify signatures. These can be shared freely and persisted or cached wherever they are to be used. Like the private key, they *must* include the same UUID assigned to the private counterpart.
 * **Signature** - A JSON Web Signature which (JWS) is always serialized to JSON in compact form. The signature headers will always include the key's UUID to make it easier to find the public key used for validation.
 * **Digest** - Defines the algorithm used to create a digest or hash of the GoBL document body and the resulting value in hexadecimal format. The digest is expected to be included in a document header and consequently in the signature payload. SHA256 digests are only supported at this time.

This package aims to make it easier to use digital signatures with GoBL documents, but it should be just as easy to use this library with any software, document, or message that could benefit from a simplified approach to dealing with JSON Web Signatures.
