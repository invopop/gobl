## Added

- `org`: `Item` supports the EN 16931 price details (BG-29): `list`, the list
  price before the price discount (BT-148); `discount`, the amount deducted
  from the list price (BT-147); and `per`, the number of units the price
  applies to (BT-149). When a list price is set, the price is calculated as
  the list price less the discount.
- `org`: `Item.UnitPrice` provides the item's price for a single unit, taking
  `per` into account.
- `bill`: line sums divide the item's price by `per`, so a price per 100 kg
  applied to 250 kg gives a sum of 2.5 times the price.

## Changed

- `org`: the `Item.Price` schema title is "Net Price".
