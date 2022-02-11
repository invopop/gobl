package gobl_test

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/uuid"
)

func ExampleNewEnvelope_complete() {
	// Prepare a new Envelope with a region
	env := gobl.NewEnvelope(region.ES)
	env.Head.UUID = uuid.MustParse("871c1e6a-8b5c-11ec-af5f-3e7e00ce5635")

	// Prepare a payload and insert
	msg := &note.Message{
		Content: "sample message content",
	}
	if err := env.Insert(msg); err != nil {
		panic(err.Error())
	}

	data, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v\n", string(data))
	// Output:
	// {
	// 	"$schema": "https://gobl.org/draft0/envelope",
	// 	"head": {
	// 		"uuid": "871c1e6a-8b5c-11ec-af5f-3e7e00ce5635",
	// 		"rgn": "ES",
	// 		"dig": {
	// 			"alg": "sha256",
	// 			"val": "d1310846172b493492c5aa1bee607d461a1878b5e535e9fb39e819d5a4125554"
	// 		}
	// 	},
	// 	"doc": {
	// 		"$schema": "https://gobl.org/draft0/note#Message",
	// 		"content": "sample message content"
	// 	},
	// 	"sigs": []
	// }
}
