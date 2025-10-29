# 🇸🇪 GOBL Sweden Tax Regime

Sweden uses the PEPPOL BIS Billing 3.0 format (CIUS based on EN 16931) for their e-invoicing system.

Find example SEGOBL files in the [`examples`](../../examples/se) (uncalculated documents) and [`examples/out`](../../examples/se/out) (calculated envelopes) subdirectories.

## Public Documentation

- [PEPPOL BIS Billing 3.0 Specification](https://docs.peppol.eu/poacc/billing/3.0/)
- [Agency for Digital Government (DIGG) - E-invoicing](https://www.digg.se/e-handel-och-e-faktura/obligatorisk-e-fakturering-i-offentlig-sektor)
- [Swedish Tax Agency (Skatterverket)](https://www.skatteverket.se/foretag.4.76a43be412206334b89800052908.html)
  - [English version](https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers.4.12815e4f14a62bc048f5159.html)
- [Swedish Tax Agency (Skatteverket) - VAT Rules](https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momssatspavarorochtjanster.4.58d555751259e4d66168000409.html)
- [2025 Official Skatterverket VAT Guide](https://www.skatteverket.se/download/18.7be5268414bea0646946f3e/1428566850726/552B14.pdf)

## E-invoicing Requirements

Since April 1, 2019, all suppliers to Swedish public authorities must send and receive invoices electronically via the PEPPOL network using the PEPPOL BIS Billing 3.0 profile. While not mandatory for B2B transactions, e-invoicing is highly recommended and effectively a de-facto standard when trading with larger enterprises.

A compliant invoice must include all Core Invoice elements as specified by EN 16931 Business Rules and the PEPPOL profile, including:

1. Date of issue of the invoice.
2. A unique serial number for each invoice based on one or more series (unique and sequential per fiscal year).
3. The seller's VAT registration number.
4. The buyer's VAT registration number if the buyer is liable for payment for the purchase, so-called reverse charge.
5. Name and address of the seller and buyer.
6. The quantity and nature of the goods or the scope and nature of the services.
7. The date on which the sale of the goods or services was made or completed or the date on which the advance or on-account payment was made, if such a date can be determined and it is different from the invoice date.
8. The tax base for each VAT rate or exemption, the unit price excluding VAT, and any price reduction or discount not included in the unit price.
9. Applied VAT rate.
9. The amount of VAT to be paid. If the seller uses profit margin taxation on the transaction, the VAT amount should not be stated on the invoice.

[Source](https://www.skatteverket.se/foretag/moms/saljavarorochtjanster/momslagensregleromfakturering.4.58d555751259e4d66168000403.html#fakturansinnehall)

## Sweden-specific Requirements

[Source](https://www4.skatteverket.se/rattsligvagledning/edition/2025.2/394342.html)

### Identification Numbers

**All identification numbers are 10 digits long**. Individuals and businesses have different formats.

The **supplier must always be identified by their tax identification number**.

If the customer is liable to pay VAT (they are registered for VAT or the operation is a reverse charge), their VAT number must also be present. Otherwise name and address are sufficient.

> **VAT ID format:** `SE` + 10-digit organization or personal number + `01` as check digits

#### Individuals:

- **Social security number (Personnummer)** - For individuals registered in Sweden.
  - Birthdate in format YYMMDD, "-" if less than 100 years old, "+" otherwise, 3-digit birthnumber (last digit odd for biological males, even for females) and a checksum.
  - Example: `990101-1230` Male born in 1999-01-01
- **Coordination number (Samordningsnummer)** - For individuals not registered in Sweden.
  - Birthdate in format YYMMDD (day + 60), "-" if less than 100 years old, "+" otherwise, 3-digit birthnumber (last digit odd for biological males, even for females) and a checksum.
  - Example: `990161+1229` Female born in 1899-01-01

#### Businesses:

- **Organization number (Organisationsnummer, `Org.nr`)**
  - First digit (the "group number") may identify the type of entity (but not necessarily):
    <details>
      <summary>Group numbers:</summary>
      <ul>
        <li>1 - Death certificate</li>
        <li>2 or 8 - Religious denominations</li>
        <li>20 - Government agencies (assigned by the Statistics Sweden)</li>
        <li>3 - Foreign companies engaged in business activities or own real estate in Sweden</li>
        <li>5 - limited liability companies (Aktiebolag), branches, banks, insurance companies and European companies</li>
        <li>6 - Single company</li>
        <li>7 or 8 - Tenant-owner associations, economic associations, non-profit associations, housing associations, cooperative tenancy associations, European cooperatives and European groupings for territorial cooperation</li>
        <li>9 - Partnerships and limited partnerships</li>
      </ul>
    </details>

  - The last digit is a checksum.
  - Sole proprietorships use the owner's personal number.
  - Example: `556036-0793` Private limited company

#### Checksum

<details>
<summary>The checksum is calculated using the <a href="https://stripe.com/en-es/resources/more/how-to-use-the-luhn-algorithm-a-guide-in-applications-for-businesses">Luhn algorithm</a>:</summary>

- Start with the payload digits. Moving from right to left, double every second digit, starting from the last digit. If doubling a digit results in a value > 9, subtract 9 from it (or sum its digits).
- Sum all the resulting digits (including the ones that were not doubled).
- The check digit is then calculated using the formula $(10 - (s \pmod{10})) \pmod{10}$, where $s$ is the sum from the previous step. This yields the smallest non-negative number which, when added to $s$, results in a multiple of 10.

</details>

### VAT Rates

In Sweden, VAT is called "Moms" (Mervärdesskatt). The following rates (Skattesatser) are used in Sweden:

| Rate            | Swedish Term          | Percentage | Description                                                                                                              |
| --------------- | --------------------- | ---------- | ------------------------------------------------------------------------------------------------------------------------ |
| Standard        | Normalskattesats      | 25%        | Most goods and services                                                                                                  |
| Reduced         | Skattesats 12 procent | 12%        | Food products, hotel accommodations, restaurant and catering services, shoe repair, leather goods, clothing, bicycles... |
| Heavily reduced | Skattesats 6 procent  | 6%         | Passenger transport, intellectual property, cultural services (except cinema), books, newspapers...                      |
| Exempt          | Momsfri               | 0%         | Exports, intra-community supplies, pharmaceuticals, certain financial and healthcare services                            |

### F-Tax

When a business is registered for F-tax, their customers do not have to deduct taxes on payments made to them for work performed in Sweden.

This is usually done by domesticsole proprietorships and foreign companies.

### Reverse Charge

Cases where reverse charge applies:

- **Domestic Transactions**: Applicable to specific sectors like construction services, trading of certain metals, waste and scrap materials, emission rights, and services related to real estate.

- **Cross-Border Transactions**:
  - Intra-Community Supplies: When goods or services are supplied between EU member states, and the customer is VAT-registered in another member state.
  - Services from Abroad: When services are provided by a supplier not established in Sweden to a VAT-registered customer in Sweden.

Implications:

- VAT rate = 0
- VAT category code "Reverse charge"
- Include both supplier and customer VAT IDs
- Add note or exemption reason "Reverse charge"
- Customer must account for VAT

### PEPPOL Technical Requirements

Invoices must reference the PEPPOL BIS identifiers:
- Specification Identifier: `urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0`
- Business Process Specified Document Context Parameter: `urn:fdc:peppol.eu:2017:poacc:billing:01:1.0`

No empty XML elements are allowed, and all mandatory fields must be present to pass validation.


## Decisions

<details>
<summary>Decision-making process and steps taken while implementing the Sweden regime:</summary>

> [!NOTE]
> These were written in a journal style, you may find them useful to understand the way the module was planned and implemented. If redundant, simply remove them.

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

- Tax Categories ([`tax_categories.go`](tax_categories.go))
- Validators and normalizers
  - Tax Identities ([`tax_identities.go`](tax_identities.go), [`tax_identities_test.go`](tax_identities_test.go))
  - Organization Identities ([`org_identities.go`](org_identities.go), [`org_identities_test.go`](org_identities_test.go))
  - Invoices ([`invoices.go`](invoices.go), [`invoices_test.go`](invoices_test.go))

And, potentially (since most regimes did include them, but not all):
- ~~Scenarios ([`scenarios.go`](scenarios.go))~~
- ~~Corrections ([`corrections.go`](corrections.go))~~

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

- How to handle the [F-Tax exemption](https://www.skatteverket.se/servicelankar/otherlanguages/inenglishengelska/businessesandemployers/startingandrunningaswedishbusiness/registeringabusiness/approvalforftax.4.676f4884175c97df4192308.html)? Is allowing users to simply select the 0% tax rate enough? Do they need an extension for it? A tax category?
- [Reverse charges have specific rules](https://marosavat.com/manual/vat/sweden/reverse-charge/), how should we implement them?
- Do I need `Corrections`?
- Benchmark the performance of the couple of `PERF` tags I added.
- Not 100% sure on what validations are required (for example, some regimes validate `Invoice`s, others don't). I assume this would become clearer had I dug into their specific documentation, but it would've massively increased the time required. For now, I added them to tax identities, org identities and invoices.
- Check the generation is actually correct (just because it worked needn't mean it meets expectations).
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
- Used `any` instead of `interface{}` since it's the "preferred" type for Go nowadays, more readable and familiar to people coming from other languages.
- I originally used an exponent of 2 for the `MakePercentage`s, but to keep it consistent with the rest of the codebase, I changed it to 3.
- For the VAT Rate dates, I didn't see any go past 1990, so I started from there.
- Parallelised the tests (`t.Parallel()`) whenever possible when they didn't have any potential race conditions or other parallelism issues. Free performance wins!
- Improved the implementation of countries so we can easily get all the data for a country from a centralised location and avoiding certain hardcodings (country name, ISO code, etc).
- Ported over the code from https://github.com/luhnmod10/go instead of adding a new dependency as it was easy enough to understand and the algorithm won't change.
- There's a [project/library to validate Swedish organisation numbers](https://organisationsnummer.dev/). We could use it to ensure greater correctness.


I heavily relied on the Polish and Spanish regimes, since the former use the Peppol BIS Billing 3.0 format and the latter I'm more familiar with. The Indian, Brazilian, Saudi and Italian ones have also been useful as reference.
</details>

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
