package gobl_test

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/uuid"
)

func ExampleNewEnvelope_complete() {
	// Prepare a new Envelope with a region
	env := gobl.NewEnvelope()
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
	// 	"$schema": "https://gobl.org/draft-0/envelope",
	// 	"head": {
	// 		"uuid": "871c1e6a-8b5c-11ec-af5f-3e7e00ce5635",
	// 		"dig": {
	// 			"alg": "sha256",
	// 			"val": "7d539c46ca03a4ecb1fcc4cb00d2ada34275708ee326caafee04d9dcfed862ee"
	// 		}
	// 	},
	// 	"doc": {
	// 		"$schema": "https://gobl.org/draft-0/note/message",
	// 		"content": "sample message content"
	// 	},
	// 	"sigs": []
	// }
}
