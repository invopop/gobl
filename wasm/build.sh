#!/bin/bash

GOOS=js GOARCH=wasm go build -o gobl.wasm

WASM_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"
if [ ! -f "$WASM_EXEC" ]; then
    WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
fi
cp "$WASM_EXEC" .

go run ./serve
