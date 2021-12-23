# What is GOBL?

GOBL stands for "Go Business Language", and is three things in one:

* A library, written in Go, for building business documents,
* A repository of global taxes and validation rules, especially for invoices, and,
* A JSON Schema, for being able to easily create and share with others.

Conventional standards focus on defining a schema, usually in XML, and leave the implementation and localisation aspects to other libraries and local state agencies. Given the flexibility of name-spacing and regional extensions, the results are usually difficult and time consuming to use and implement. They also usually lack tax definitions and validation rules, so it's up to developers to ensure what they're creating will actually work, usually through trial and error.

The initial focus of GOBL is in the area of Electronic Invoicing. We believe it is taking too long for electronic invoice to become common place, and part of the reason for that is the complexity of current standards. GOBL was created for developers, by developers, with the aim of making it easy to convert sales into fiscally valid invoices.

Our aim is to leverage open source practices to create a single global public library that consolidates all the existing public information on taxes and validation rules.
