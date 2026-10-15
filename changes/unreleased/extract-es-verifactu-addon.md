## Changed

- `addons/es/verifactu`: **breaking**: the Spanish VERI*FACTU (`es-verifactu-v1`) addon moved to the standalone [`github.com/invopop/gobl.es.verifactu`](https://github.com/invopop/gobl.es.verifactu) module, alongside the VeriFactu converter that consumes it. Add a blank import (`_ "github.com/invopop/gobl.es.verifactu/addon"`) to keep using the `es-verifactu-v1` addon key. The key itself is unchanged, and remains a valid `$addons` value through the approved external addon list.
