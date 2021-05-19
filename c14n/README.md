# GoBL Canonicalization

## Introduction

One of the hardest issues to solve around digital signatures is generating a digest or hash of the source data consistently and despite any structural changes that may have a happened while in transit. The objective of canonicalization is to ensure that the data is logically equivalent at the source and destination so that the digest can be calculated reliably on both sides, and thus be used in digital signatures. 

In the world of XML, the [Canonical XML Version 1.1](https://www.w3.org/TR/xml-c14n11/) W3C recommendation aims to set out rules to be used to create consistent documents. Anyone who's worked with XML signatures however ([XML-DSIG](https://www.w3.org/TR/xmldsig-core/)) knows that despite the best intentions and libraries, it can still be difficult to get the expected results, especially when using different languages at the source and destination.

JSON conversely lacks clearly defined or active industrial standard around canonicalization, despite having a much simpler syntax. Indeed, the [JSON Web Token specification](https://datatracker.ietf.org/doc/html/rfc7519) gets around pesky canonical issues by including the actual signed payload data as a Base 64 string inside the signatures.

One of the objectives of GoBL is to create a document that could potentially be stored in any key-value format alternative to JSON, like [YAML](https://yaml.org/), [Protobuf](https://developers.google.com/protocol-buffers), or maybe even XML. Perhaps GoBL documents need to be persisted to a document database like [CouchDB](https://couchdb.apache.org/) or a [JSONB field in PostgreSQL](https://www.postgresql.org/docs/13/functions-json.html). It should not matter what the underlying format or persistence engine is, as long as the logical contents are exactly the same. Thus when signing documents it's essential we have a reliable canonical version of JSON, even if the data is stored somewhere else.

This `c14n` package, inspired by the works of others, thus aims to define a simple standardized approach to canonical JSON that could potentially be implemented easily in other languages. More than just a definition, the code here is a reference implementation from which other implementations can be made in languages other than Go.

## GoBL JSON C14n

GoBL considers the following JSON values as explicit types:

* a string
* a number, which extends the JSON spec and is split into:
    * an integer
    * a float
* an object
* an array
* a boolean
* null

JSON in canonical form:

1. MUST be encoded in VALID [UTF-8](https://tools.ietf.org/html/rfc3629). A document with invalid character encoding will be rejected.
2. MUST NOT include insignificant whitespace.
3. MUST order the attributes of objects lexicographically by the UCS (Unicode Character Set) code points of their names.
4. MUST remove attributes from objects whose value is `null`, but maintain them in arrays.
5. MUST represent integer numbers, those with a zero-valued fractional part, WITHOUT:
    1. a leading minus sign when the value is zero,
    2. a decimal point,
    3. an exponent, thus limiting numbers to 64 bits, and
    4. insignificant leading zeroes, as already required by JSON.
6. MUST represent floating point numbers in exponential notation, INCLUDING:
    1. a nonzero single-digit part to the left of the decimal point,
    2. a nonempty fractional part to the right of the decimal point,
    3. no trailing zeroes to right of the decimal point except to comply with the previous point,
    4. a capital `E` for the exponent indicator,
    5. no plus sign in the [mantissa](https://en.wikipedia.org/wiki/Significand) nor exponent, and
    6. no insignificant leading zeros in the exponent.
7. MUST represent all strings, including object attribute keys, in their minimal length UTF-8 encoding:
    1. using two-character escape sequences where possible for characters that require escaping, specifically:
        * `\"` U+0022 Quotation Mark
        * `\\` U+005C Reverse Solidus (backslash)
        * `\b` U+0008 Backspace
        * `\t` U+0009 Character Tabulation (tab)
        * `\n` U+000A Line Feed (newline)
        * `\f` U+000C Form Feed
        * `\r` U+000D Carriage Return
    2. using six-character `\u00XX` uppercase hexadecimal escape sequences for control characters that require escaping but lack a two-character sequence described previously, and
    3. reject any string containing invalid encoding.

The GoBL JSON c14n package has been designed to operate using any raw JSON source and uses the Go [`encoding/json`](https://golang.org/pkg/encoding/json/) library's streaming methods to parse and recreate a document in memory. A simplified object model is used to map JSON structures ready to be again converted back into canonical JSON.

## Usage Example

```go
d := `{ "foo":"bar", "c": 123.4, "a": 56, "b": 0.0, "y":null}`
r := strings.NewReader(data)
res, err := c14n.CanonicalJSON(r)
if err != nil {
  panic(err.Error())
}
fmt.Printf("Result: %v\n", string(res))
// Result: {"a":56,"b":0.0E0,"c":1.234E2,"foo":"bar"}
```

## Prior Art

This specification and implementation is based on the [gibson042 canonicaljson specification](https://gibson042.github.io/canonicaljson-spec/) with simplifications concerning invalid UTF-8 characters, null values in objects, and a reference implementation that is more explicit making it potentially easier to be recreate in other programming languages.

The gibson042 specification is in turn based on the now expired [JSON Canonical Form internet draft](https://datatracker.ietf.org/doc/html/draft-staykov-hu-json-canonical-form-00) which lacks clarity on the handling of integer numbers, is missing details on escape sequences, and doesn't consider invalid UTF-8 characters.

Canonical representation of floats is consistent with [XML Schema 2, section 3.2.4.2](https://www.w3.org/TR/xmlschema-2/#float-canonical-representation), and expects integer numbers without an exponential component as defined in [RFC 7638 - JSON Web Key Thumbprint](https://datatracker.ietf.org/doc/html/rfc7638#section-3.3).





