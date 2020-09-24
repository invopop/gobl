# GoBL

Go Business Language. 

## Introduction

GoBL (pronounced "gobble") is a hybrid solution between a document standard, a software library, and a database. On one side it defines a language for business documents in JSON, while at the same time providing a library and server written in Go to help build, validate, sign, and localise.

Traditionally software business document standards consist of a set of definitions layered into different namespaces, followed by an independent set of implementations. The theory is that a well written standard can be implemented anywhere and be compatible with any existing library, in practice however, and especially for XML-base standards, dealing with multiple namespaces adds a lot of complexity. This tends to be reflected in lower quality libraries.

For GoBL a different approach is being tested. The code and library implementation is in itself the standard. Rather than trying to be super-flexible at the cost of complexity, GoBL includes everything needed for digital signatures, validation, and local implementations including tax definitions, all maintained under an open source license.

In our opinion, Go is an ideal language for this type of project due to its simple and concise syntax, performance, testing capabilities, and portability. We understand however that Go won't be everyone's cup of tea, so the project is designed to offer a server component (you could call it a microserivce) to be launched in a docker container inside your own infrastructure. Using gPRC, a simplified API helps you build and validate your business documents in your language of choice.

## Serialization

### Amounts

Marshalling numbers can be problematic with JSON as the standard dictates that numbers should be represented as integers or floats, without any tailing 0s. This is fine for maths problems, but not useful when trying to convey monetary values or rates with a specific level of accuracy. GoBL will always serialise Amounts as strings to avoid any potential issues with number conversion.





