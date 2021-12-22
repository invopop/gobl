---
description: Dealing with numbers in business documents reliably.
---

# Numbers

Marshalling numbers can be problematic. Business documents, and especially Invoices, need to take into account the potential of rounding errors when making calculations.

The [JSON RFC](https://datatracker.ietf.org/doc/html/rfc7159.html#section-6) dictates that numbers should be represented as integers or floats, without any tailing 0s, which matches expectations with most programming languages. This is fine for mathematical problems, but can very easily lead to rounding errors when trying to convey monetary values or rates with a specific level of accuracy.

GOBL will always parse and serialise amounts and percentages as **strings** to avoid any potential issues with number conversion. In effect, numbers are not allowed in GOBL documents except for specific use-cases such as counters or indexes.

Amounts therefore:

* Contain a single optional decimal place always represented with `.`&#x20;
* In libraries will be dealt with as integers, never floats.
* Have tailing `0`s to determine the significant digits, e.g. `1.000` implies we're dealing with an integer value of `1000`, and thus an accuracy of three decimal places.
* When performing arithmetic will always convert the incoming amount's value to have the same exponential as the base.
* Does not include support for exponential values in output, like `2.0e10`, as amounts are typically used for consumption by humans and not complex mathematical calculations with large numbers.

Percentages use the same base as amounts, but:

* Require a `%` symbol at the end of the serialised value, for example `10.0%`
* Will always be converted with a factor of 100, so `16.0%` implies an underlying value of `0.160` or `160` as an integer, in order to make calculations.
