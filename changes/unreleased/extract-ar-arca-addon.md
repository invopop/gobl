## Changed

- `addons/ar/arca`: **breaking**: the Argentina ARCA (`ar-arca-v4`) addon moved to the standalone [`github.com/invopop/gobl.ar.arca`](https://github.com/invopop/gobl.ar.arca) module. Add a blank import (`_ "github.com/invopop/gobl.ar.arca/addon"`) to keep using the `ar-arca-v4` addon key. The key itself is unchanged, and remains a valid `$addons` value through the approved external addon list.
