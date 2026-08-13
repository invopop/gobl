//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/rules"
)

func main() {
	if err := generate(); err != nil {
		panic(err)
	}
}

func generate() error {
	for _, s := range rules.AllSets() {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(s); err != nil {
			return err
		}
		n := s.Package
		f := filepath.Join("data", "rules", n+".json")
		if err := os.WriteFile(f, buf.Bytes(), 0644); err != nil {
			return err
		}
		fmt.Printf("Processed %v\n", f)
	}
	return nil
}
