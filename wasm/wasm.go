//go:build js

// Package main provides a WASM interface to GOBL.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall/js"

	"github.com/invopop/gobl/internal/cli"
)

func main() {
	r, w := io.Pipe()
	js.Global().Call("addEventListener", "message", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsonEvent := js.Global().Get("JSON").Call("stringify", args[0].Get("data")).String()
		if _, err := fmt.Fprintln(w, jsonEvent); err != nil {
			fmt.Fprintln(os.Stderr, "failed to send event: %s", err)
		}
		return nil
	}))

	js.Global().Call("postMessage", map[string]any{"ready": true})

	processMessages(r)

	fmt.Println("exiting")
}

func processMessages(r io.Reader) {
	bulkOpts := &cli.BulkOptions{
		In: r,
	}
	for result := range cli.Bulk(context.TODO(), bulkOpts) {
		postMessage(result)
	}
}

func postMessage(result *cli.BulkResponse) {
	response := js.Global().Get("Object").New()
	if result.ReqID != "" {
		response.Set("req_id", result.ReqID)
	}
	response.Set("seq_id", result.SeqID)
	if len(result.Payload) > 0 {
		response.Set("payload", string(result.Payload))
	}
	if result.Error != nil {
		ed, err := json.Marshal(result.Error)
		if err != nil {
			response.Set("error", fmt.Sprintf("failed to marshal error: %s", err.Error()))
		}
		response.Set("error", string(ed))
	}
	if result.IsFinal {
		response.Set("is_final", true)
	}
	js.Global().Call("postMessage", response)
}
