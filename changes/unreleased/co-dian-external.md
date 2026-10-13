## Changed

- `co-dian-v2`: **breaking**: the Colombia DIAN addon moved to the standalone
  [`github.com/invopop/gobl.co.dian`](https://github.com/invopop/gobl.co.dian)
  module. Add a blank import (`_ "github.com/invopop/gobl.co.dian/addon"`) to
  keep using the `co-dian-v2` addon key. The key itself is unchanged, and
  remains a valid `$addons` value through the approved external addon list. The
  Colombian tax regime (`regimes/co`) stays in GOBL core.
