# @invopop/gobl-worker

Node.js package to help use the GOBL CLI WASM binary distributed via CDN. This package at the moment won't try to embed the `.wasm` as an asset.

Versions of this package always match the GOBL version which was uploaded to the CDN.

It turns out providing workers from npm packages is actually quite hard. This library uses Vite to embed the `worker.js` inside the main `gobl-worker.js` file as Base64 data to avoid issues loading external JS assets.

## Usage

Install the GOBL worker into your project:

```bash
npm install @invopop/gobl-worker
```

Grab namespace and run operations:

```typescript
import * as gobl from "@invopop/gobl-worker";

const data = await gobl.ping();
const result = JSON.parse(data);

if (result.pong) {
  console.log("Success!");
}
```

## Notice

This library is still in experimental stages and requires a refactor to make interfaces much easier to use.

More tests and consolidation of the `/wasm` project also including in this repository is required.
