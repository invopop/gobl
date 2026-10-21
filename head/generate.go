//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

func main() {
	if err := generate(); err != nil {
		panic(err)
	}
}

// generate outputs the list of registered stamp definitions to the data
// directory.
func generate() error {
	data, err := json.MarshalIndent(head.AllStampDefs(), "", "  ")
	if err != nil {
		return err
	}
	f := filepath.Join("data", "stamps.json")
	if err := os.WriteFile(f, data, 0644); err != nil {
		return err
	}
	fmt.Printf("Processed %v\n", f)
	return nil
}
