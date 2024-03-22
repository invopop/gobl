//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

func main() {
	if error := generate(); error != nil {
		panic(error)
	}
}

func generate() error {
	for _, r := range tax.AllRegimes() {
		doc, err := schema.NewObject(r)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return err
		}
		n := string(r.Country)
		if r.Zone != "" {
			n = n + "_" + string(r.Zone)
		}
		n = strings.ToLower(n)
		f := filepath.Join("data", "regimes", n+".json")
		if err := os.WriteFile(f, data, 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}
