#!/bin/bash

GOOS=js GOARCH=wasm go build -o gobl.wasm

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

go run ./serve
