# Currency Reference Data

Currency support in GOBL aims to cover most of the currencies in the world and provide methods to format and work with them.

List is based on ISO 4217.

Much of the inspiration for this package and source data in the `./data` directory originally came from the excellent and widely used [money gem in ruby](https://rubymoney.github.io/money/). A few alterations to source data have been made.

Currencies around the world change more often than expected, please [send a PR](https://github.com/invopop/gobl/pulls) if you spot anything that is out of date, along with a link that references the source of the change.

In order to keep application size low, GOBL will only load and unmarshall currency data if requested.
