## Changed

- `addons/pl/favat`: **breaking**: the Polish KSeF FA_VAT (`pl-favat-v3`) addon moved to the standalone [`github.com/invopop/gobl.pl.ksef`](https://github.com/invopop/gobl.pl.ksef) module, alongside the KSeF converter that consumes it. Add a blank import (`_ "github.com/invopop/gobl.pl.ksef/addon"`) to keep using the `pl-favat-v3` addon key. The key itself is unchanged, and remains a valid `$addons` value through the approved external addon list.
