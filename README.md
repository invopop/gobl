# GOBL

<img src="https://github.com/invopop/gobl/blob/main/gobl_logo_black_rgb.svg?raw=true" width="181" height="219" alt="GOBL Logo">

Go Business Language.

Released under the Apache 2.0 [LICENSE](https://github.com/invopop/gobl/blob/main/LICENSE), Copyright 2021,2022 [Invopop Ltd.](https://invopop.com).

## Introduction

GOBL (pronounced "gobble") is a hybrid solution between a document standard and a software library. On one side it defines a language for business documents in JSON, while at the same time providing a library and services written in Go to help build, validate, sign, and localise.

Traditional software business document standards consist of a set of definitions layered into different namespaces, followed by an independent set of implementations. A well written standard can be implemented anywhere and be compatible with any existing library. In practice however, and especially for XML-base standards, dealing with multiple namespaces adds a lot of complexity.

For GoBL a different approach is being taken. The code and library implementation is in itself the standard. Rather than trying to be flexible at the cost of complexity, GoBL includes everything needed for digital signatures, validation, and local implementations including tax definitions, all maintained under an open source license.

In our opinion, Go is an ideal language for this type of project due to its simple and concise syntax, performance, testing capabilities, and portability. We understand however that Go won't be everyone's cup of tea, so the project is designed to offer a server component (you could call it a microservice) to be launched in a docker container inside your own infrastructure, alongside a JSON Schema definition of all the key documents. Building your own documents and using the GoBL services to validate and sign them should be relatively straight forward.

## GoBL Standard

### Packages

### Documents

### Envelope

## Serialization

### Amounts & Percentages

Marshalling numbers can be problematic with JSON as the standard dictates that numbers should be represented as integers or floats, without any tailing 0s. This is fine for maths problems, but not useful when trying to convey monetary values or rates with a specific level of accuracy. GoBL will always serialise Amounts as strings to avoid any potential issues with number conversion.

### ID vs UUID

Traditionally when dealing with databases a sequential ID is assigned to each tuple (document, row, item, etc.) starting from 1 going up to whatever the integer limit is. Creating a scalable and potentially global system however with regular numbers is not easy as you require a single point in the system responsible for ensuring that numbers are always provided in the correct order. Single points of failure are not good for scalability.

To get around the issues with sequential numbers, the [UUID standard](https://tools.ietf.org/html/rfc4122) (Universally Unique IDentifier) was defined as a mechanism to create a very large number that can be potentially used to identify anything.

The GoBL library forces you to use UUIDs instead of sequential IDs for anything that could be referenced in the future. To enforce this fact, instead of naming fields `id`, we make an effort to call them `uuid`.

Sometimes sequential IDs are however required, usually for human consumption purposes such as ensuring order when generating invoices so that authorities can quickly check that dates and numbers are in order. Our recommendation for such codes is to use a dedicated service to generate sequential IDs based on the UUIDs, such as [Invopop's Sequence Service](https://invopop.com).


## CLI tool

This repo contains a `gobl` CLI tool which can be used to manipulate GOBL documents from the command line or shell scripts.

### `gobl build`

Build a complete GOBL document from one or more input sources.  Example uses:

```sh
# Finalize a complete invoice
gobl build invoice.yaml

# Set the supplier from an external file
gobl build invoice.yaml \
    --set-file doc.supplier=supplier.yaml

# Set arbitrary values from the command line. Inputs are parsed as YAML.
gobl build invoice.yaml \
    --set doc.foo.bar="a long string" \
    --set doc.foo.baz=1234

# Set an explicit string value (to avoid interpetation as a boolean or number)
gobl build invoice.yaml \
    --set-string doc.foo.baz=1234 \
    --set-string doc.foo.quz=true
```
