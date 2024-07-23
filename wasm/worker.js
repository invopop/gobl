self.importScripts("wasm_exec.js");

// Polyfill instantiateStreaming for browsers missing it
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// initialize the Go WASM glue
const go = new self.Go();
WebAssembly.instantiateStreaming(fetch("gobl.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
  })
  .catch((err) => {
    console.error("Failed to run GOBL WASM instance: ", err);
  });

// console.log("worker.js");
// addEventListener('message', async (e) => {
//     console.log("worker.js handler");

//     // tell the main thread we are done
//     postMessage({
//         done: true
//     });
// }, false);
