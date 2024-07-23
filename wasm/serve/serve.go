// Package main is used to launch a demo site for the wasm.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Println("Navigate to http://localhost:9999/")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
