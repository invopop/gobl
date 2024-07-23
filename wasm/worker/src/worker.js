import "./wasm_exec.js";

const goblCliVersion = GOBL_CLI_VERSION;

const wasmUrl = `https://cdn.gobl.org/cli/gobl.${goblCliVersion}.wasm`;

// Polyfill instantiateStreaming for browsers missing it.
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// Initialize the Go WASM glue.
const go = new Go();
WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
  .then((result) => {
    go.run(result.instance);
  })
  .catch((err) => {
    console.error("Failed to run GOBL WASM instance: ", err);
  });
