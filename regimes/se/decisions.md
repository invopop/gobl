# Reasoning and decisions for Sweden

## Initial steps

My approach was first to try to understand how other countries were implemented, distill the commonalities and then apply them to Sweden. This would allow me to understand what code I had to actually write, and what could be reused, as well as how the project was structured, what things I could forego focusing on and ensure consistency with the rest of the project.

A quick script to get the last commit date for each regime (that way I could see which ones were most recently added):
```sh
for d in ./regimes/*/; do
    [ -d "$d" ] || continue
    date=$(git log --diff-filter=A --format="%aI" -- "$d" | tail -1)
    printf "%-30s %s\n" "$d" "$date"
done | sort -k2 -r
```

I also looked at previous PRs to see what other users had implemented, such as [#440](https://github.com/invopop/gobl/pull/440/files), which in turn came from [#433](https://github.com/invopop/gobl/pull/433), to understand what feedback they had received and what they had done to address it.

Obviously, I also started gathering information on Sweden's tax system on the side (which was quite a bit). I also came at it from the assumption that at least one of the two involved parties is registered for VAT in Sweden (most likely the supplier). Otherwise, they would issue invoices in some other country's format.

AI came in handy to get useful sources to look at (I couldn't blindly trust the research since LLM models are known to hallucinate) and for an initial project breakdown. I manually checked certain folders, starting with `regimes` to get a sense of the codebase and start getting familiarised.

My ultimate goal was to break down the problem into smaller, manageable chunks, so I could deliver a working implementation that resembled existing ones and passed the tests. After that, I could refactor and add any additional features as we saw fit based on feedback.

Once I had an initial blueprint (started with the Polish implementation), I stripped away whatever I could to work on the MVP. I also noticed everything had tests, so knowing if the things I chose to implement at least worked as expected would be as simple as making the tests pass and thus I could apply a TDD-esque approach.

The identified requirements were:

- Tax Categories ([tax_categories.go](tax_categories.go))
- Validators and normalizers
  - Tax Identities ([tax_identity.go](tax_identity.go))
  - Organization Identities ([org_identities.go](org_identities.go))
  - Invoices ([invoices.go](invoices.go), [invoices_test.go](invoices_test.go))

And, potentially (since most regimes did include them, but not all):
- Scenarios ([scenarios.go](scenarios.go))
- Corrections ([corrections.go](corrections.go))

This list doesn't include ensuring that the new regime was included everywhere it needed to be (e.g. [`regimes.go`](../regimes.go), adding examples...).

With a starting set of features to work on and a foundation of knowledge, I started working on it.

## Implementation

Many fields, such as `Tags` or `Addons`, are optional by definition and according to the GOBL JSON Schema, so I decided to not add any apart from the default, if any at all.

### Tax Identities

When normalising Tax IDs (and since they are different from Identification numbers, but directly derived from them by adding a prefix and suffix), I kept them with the prefix and suffix. Also, [Peppol BIS Billing 3.0 requires it](https://docs.peppol.eu/poacc/billing/3.0/rules/ubl-peppol/SE-R-001/). However, IDs are normalized to only numeric characters. This allows direct conversion between the two if needed.

For tax code validation, since Swedish ones are so simple, I decided not to use RegExp and instead simply check the length, prefix (`SE`), suffix (`01`) and 10-number identifier, as it would be easier to understand and highly likely faster too (benchmark pending, obviously, as performance should always be measured - hunches and educated guesses are not enough).

### Organization Identities


For personal ID numbers, since the `-`/`+` is significant (indicates if a person is over 100 years old), it must be preserved during normalisation. Therefore I couldn't use the common functions for normalisation either.

Also, if no separator is present but we have the right number of digits, I decided to insert a hyphen at the right position, since it's the most statistically likely separator and that was better than failing validation later due to not being able to normalise the ID.

### Others

For reverse charges, I went with [Stripe's guide on the matter](https://stripe.com/en-es/guides/introduction-to-eu-vat-and-european-vat-oss), which provides a [handy flowchart](https://images.stripeassets.com/fzn2n1nzq965/4W51jI9ssA42jiX43IYjq6/1816da1e38f7f7e01d212c980cd4fdd7/Tax_guide_phyiscal_vs_digital_GB.png?w=2550&q=80&fm=webp), as well as this [other guide by Marosa](https://marosavat.com/manual/vat/sweden/reverse-charge/). The summary is present in the [README](./README.md).

## Questions/TO-DOs

- Generate the [`examples/se/out`](../examples/se/out/) directory.
- Is my handling of the Tax IDs' `SE` prefix and `01` suffix correct? I wasn't entirely sure [leaving the `SE` prefix and `01` suffix](#tax-identities) was the correct move, so I'd like to get feedback on that.
- How to handle the [F-Tax exemption](https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/registeringabusiness/approvalforftax.4.676f4884175c97df4192308.html)? Is allowing users to simply select the 0% tax rate enough? Do they need an extension for it? A tax category?
- [Reverse charges have specific rules](https://marosavat.com/manual/vat/sweden/reverse-charge/), how should we implement them?
- Difference between a tax identity and an organisation identity? I have a rough idea on their definition: tax is specifically for invoicing, whereas organisation can be more general. Am I in the right ballpark?
- Do I need `Corrections`?
- Benchmark the performance of the couple of `PERF` tags I added.
- Not 100% sure on what validations are required (for example, some regimes validate `Invoice`s, others don't). I assume this would become clearer had I dug into their specific documentation, but it would've massively increased the time required. For now, I added them to tax identities, org identities and invoices.
- Check the generation is actually correct (just because it worked needn't mean it meets expectations).
- Create a sister PR for docs.
- Add extra validation to invoices:
  - Check the invoice ID is unique for the year and has a monotonically increasing number (not even sure if this is possible, plus it may be done already).
  - Ensure the generated invoice is compliant with PEPPOL BIS Billing 3.0 and can actually be delivered as an e-invoice -> [Bare minimum](https://docs.peppol.eu/poacc/billing/3.0/rules/ubl-peppol/).
- Revisit scenarios and invoices and their tests to ensure we're not missing nor duplicating any checks.
- Add localisation for Swedish.
- Add a check where the total of a [simplified invoice cannot exceed 4000 SEK (page 3)](https://www.skatteverket.se/download/18.7be5268414bea0646946f3e/1428566850726/552B14.pdf).

## Miscellanea

- I did my best to adhere to existing patterns and conventions, and maintain consistency with the rest of the codebase.
- I did my best to document all the code and decisions in a brief but comprehensive manner, avoiding redundant comments.
- I followed the Boy Scout rule and applied corrections and improvements wherever I saw fit. Any changes substantial enough to add noise are added in a separate PR.
  - I added a [`CONTRIBUTING.md`](CONTRIBUTING.md) file to the repo to make it easier for others to contribute.
  - I added a `tools.go` file as well as including the golangci-lint version in the `go.mod` file, to ensure the same setup is used everywhere. However, including tools in the `go.mod` is not the best practice, so open to moving it to its own separate module.
- Used `any` instead of `interface{}` since it's the "preferred" type for Go nowadays, more readable and familiar to people coming from other languages.
- I originally used an exponent of 2 for the `MakePercentage`s, but to keep it consistent with the rest of the codebase, I changed it to 3.
- For the VAT Rate dates, I didn't see any go past 1990, so I started from there.
- Parallelised the tests (`t.Parallel()`) whenever possible when they didn't have any potential race conditions or other paralellisim issues. Free performance wins!
- Improved the implementation of countries so we can easily get all the data for a country from a centralised location and avoiding certain hardcodings (country name, ISO code, etc).
- Ported over the code from https://github.com/luhnmod10/go instead of adding a new dependency as it was easy enough to understand and the algorithm won't change.

## Suggestions

- Adding a template `regime` with plenty of comments and documentation explaining what things should be added and with examples, would also be a great help for others to contribute (albeit maybe painful to maintain?).
- For each regime, depend more on the already existing `l10n` package.
- For each country data in [l10n/countries.go](../../l10n/countries.go), add basic data such as the country's name in the official langugage, the timezone, currency...
- This *may* be a faster way to convert a `Code` into pure digits:
```go
digitsOnly := ""
for _, c := range code.String() {
	if c >= '0' && c <= '9' {
		digitsOnly += string(c)
	}
}
```
- So long as the checklist for the PRs is a manual "yes I have done this", you won't ever truly know if users have actually done those actions. I would expect CI/CD pipelines to provide those assurances.
- There's a [project/library to validate Swedish organisation numbers](https://organisationsnummer.dev/). We could use it to ensure greater correctness.
---

# Sources

Non-exhaustive list of sources used to gather information on Sweden's invoicing system:

- [PEPPOL BIS Billing 3.0](https://docs.peppol.eu/poacc/billing/3.0/)
- https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/321516.html
- https://marosavat.com/manual/vat/sweden/
- https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/declaringtaxesbusinesses/vat/vatratesandvatexemption.4.676f4884175c97df419255d.html
- https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momssatspavarorochtjanster.4.58d555751259e4d66168000409.html
- https://www.oecd.org/content/dam/oecd/en/topics/policy-sub-issues/consumption-tax-trends/consumption-tax-trends-sweden.pdf
- https://vatapp.net/vat-rates
- https://www.fonoa.com/countries/sweden
- https://www.skatteverket.se/privat/folkbokforing/samordningsnummer.4.5c281c7015abecc2e201130b.html
- https://en.wikipedia.org/wiki/Personal_identity_number_(Sweden)
- https://org-id.guide/list/SE-ON
- https://www.skatteverket.se/foretag/drivaforetag/startaochregistrera/organisationsnummer.4.361dc8c15312eff6fd235d1.html
- https://fakturasolid.se/en/invoice-sweden
- https://stripe.com/en-es/guides/introduction-to-eu-vat-and-european-vat-oss
- https://marosavat.com/manual/vat/sweden/reverse-charge/
- https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momslagensregleromfakturering.4.58d555751259e4d66168000403.html
- [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm)
- [How to use the Luhn algorithm](https://stripe.com/en-es/resources/more/how-to-use-the-luhn-algorithm-a-guide-in-applications-for-businesses)
- https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/registeringabusiness/approvalforftax.4.676f4884175c97df4192308.html
- https://www.vatcalc.com/sweden/sweden-vat-country-guide/
- [Official VAT Brochure by Skatteverket](https://www.skatteverket.se/download/18.7be5268414bea0646946f3e/1428566850726/552B14.pdf)
- [Official guide for contents of an invoice](https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/394342.html)
- [Requirement to include address of both supplier and customer](https://www4.skatteverket.se/rattsligvagledning/415583.html?date=2022-09-20)
- [Information to identify supplier in simplified invoice](https://www4.skatteverket.se/rattsligvagledning/28101.html?date=2005-06-27)


I also heavily relied on the Polish and Spanish regimes, since the former use the Peppol BIS Billing 3.0 format and the latter I'm more familiar with. The Indian, Brazilian, Saudi and Italian ones have also been useful as reference.
