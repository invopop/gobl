## Added

- `org`: `Item.Pricing` describes how an item's price was determined (EN 16931
  BG-29): `per`, the number of units the price applies to (BT-149); `gross`,
  the price before the discount (BT-148); and `discount`, the amount deducted
  from the gross price (BT-147). When a gross price is set, the item's price
  is calculated as the gross price less the discount.
- `org`: `Item.UnitPrice` provides the item's price for a single unit, taking
  `pricing.per` into account.
- `bill`: line sums divide the item's price by `pricing.per`, so a price per
  100 kg applied to 250 kg gives a sum of 2.5 times the price.
