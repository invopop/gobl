package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func TestBulk(t *testing.T) { //nolint:gocyclo
	type tt struct {
		opts *BulkOptions
		want []*BulkResponse
	}

	tests := testy.NewTable()
	tests.Add("invalid input", tt{
		opts: &BulkOptions{
			In: strings.NewReader("this ain't json"),
		},
		want: []*BulkResponse{
			{
				SeqID: 1,
				Error: &Error{
					Code:    422,
					Message: "invalid character 'h' in literal true (expecting 'r')",
				},
				IsFinal: true,
			},
		},
	})
	tests.Add("no input", tt{
		opts: &BulkOptions{
			In: strings.NewReader(""),
		},
		want: []*BulkResponse{
			{
				SeqID:   1,
				IsFinal: true,
			},
		},
	})
	tests.Add("one verification", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "verify",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":      base64.StdEncoding.EncodeToString(payload),
				"publickey": publicKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					Payload: []byte(`{"ok":true}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("two verifications", func(_ *testing.T) interface{} {
		req1, _ := json.Marshal(map[string]interface{}{
			"action":  "sleep",
			"req_id":  "abc",
			"payload": "10ms",
		})
		req2, _ := json.Marshal(map[string]interface{}{
			"action":  "sleep",
			"req_id":  "def",
			"payload": "50ms",
		})
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req1), bytes.NewReader(req2)),
			},
			want: []*BulkResponse{
				{
					ReqID:   "abc",
					SeqID:   1,
					Payload: []byte(`{"sleep":"done"}`),
					IsFinal: false,
				},
				{
					ReqID:   "def",
					SeqID:   2,
					Payload: []byte(`{"sleep":"done"}`),
					IsFinal: false,
				},
				{
					SeqID:   3,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("success then failure", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "verify",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":      base64.StdEncoding.EncodeToString(payload),
				"publickey": publicKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req), strings.NewReader("not json")),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					Payload: []byte(`{"ok":true}`),
					IsFinal: false,
				},
				{
					SeqID: 2,
					Error: &Error{
						Code:    422,
						Message: "invalid character 'o' in literal null (expecting 'u')",
					},
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("non-fatal payload error", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action":  "verify",
			"req_id":  "asdf",
			"payload": "not an object",
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req)),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
					Error: &Error{
						Code:    422,
						Message: "json: cannot unmarshal string into Go value of type cli.VerifyRequest",
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("non-fatal data error", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action": "verify",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":      json.RawMessage(`"oink"`),
				"publickey": publicKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req)),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
					Error: &Error{
						Code:    400,
						Message: `error converting YAML to JSON: yaml: invalid leading UTF-8 octet`,
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("one build, already signed", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "build",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":       base64.StdEncoding.EncodeToString(payload),
				"privatekey": privateKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Payload: json.RawMessage(`{
						"$schema": "https://gobl.org/draft-0/envelope"
					}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("one build, field errors", func(t *testing.T) interface{} {
		payload := []byte(`{
			"$schema":"https://gobl.org/draft-0/note/message",
			"title":"This is a title"
		}`)
		req, err := json.Marshal(map[string]interface{}{
			"action": "build",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(payload),
				// "privatekey": privateKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
					Error: &Error{
						Code: 422,
						Key:  cbc.Key("validation"),
						Fields: gobl.FieldErrors{
							"content": errors.New("cannot be blank"),
						},
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("one build, success", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/nosig.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "build",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":       base64.StdEncoding.EncodeToString(payload),
				"privatekey": privateKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Payload: json.RawMessage(`{
						"$schema": "https://gobl.org/draft-0/envelope"
					}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("build, invalid doc type", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/nosig.json")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := json.Marshal(map[string]interface{}{
			"action": "build",
			"payload": map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(payload),
				"type": "chicken",
			},
		})
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					SeqID: 1,
					Error: &Error{
						Code:    400,
						Message: `unrecognized doc type: "chicken"`,
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("build, invalid template", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/nosig.json")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := json.Marshal(map[string]interface{}{
			"action": "build",
			"payload": map[string]interface{}{
				"data":     base64.StdEncoding.EncodeToString(payload),
				"template": "chicken",
			},
		})
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					SeqID: 1,
					Error: &Error{
						Code:    422,
						Message: "invalid payload: illegal base64 data at input byte 4",
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("non-fatal payload error, build", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action":  "build",
			"req_id":  "asdf",
			"payload": "not an object",
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req)),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
					Error: &Error{
						Code:    422,
						Message: `invalid payload: json: cannot unmarshal string into Go value of type cli.BuildRequest`,
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("non-fatal data error, build", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action": "build",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":      base64.StdEncoding.EncodeToString([]byte(`"oink"`)),
				"publickey": publicKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req)),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
					Error: &Error{
						Code:    400,
						Message: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `oink` into map[string]interface {}",
					},
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("correct, success", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "correct",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":    base64.StdEncoding.EncodeToString(payload),
				"options": []byte(`{"type":"credit-note","ext":{"es-facturae-correction":"01"}}`),
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Payload: json.RawMessage(`{
						"$schema": "https://gobl.org/draft-0/envelope"
					}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("correct, options", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "correct",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":   base64.StdEncoding.EncodeToString(payload),
				"schema": true,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Payload: json.RawMessage(`{
						"$id":"https://gobl.org/draft-0/bill/correction-options?tax_regime=es"
					}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("replicate, success", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/success.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "replicate",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(payload),
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Payload: json.RawMessage(`{
						"$schema": "https://gobl.org/draft-0/envelope"
					}`),
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("unknown action", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action": "frobnicate",
			"req_id": "asdf",
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: io.MultiReader(bytes.NewReader(req)),
			},
			want: []*BulkResponse{
				{
					ReqID: "asdf",
					SeqID: 1,
					Error: &Error{
						Code:    400,
						Message: "unrecognized action: 'frobnicate'",
					},
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("keygen", func(t *testing.T) interface{} {
		req, err := json.Marshal(map[string]interface{}{
			"action": "keygen",
			"req_id": "asdf",
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("ping", tt{
		opts: &BulkOptions{
			In: strings.NewReader(`{"action":"ping"}`),
		},
		want: []*BulkResponse{
			{SeqID: 1},
			{SeqID: 2, IsFinal: true},
		},
	})
	tests.Add("sleep", tt{
		opts: &BulkOptions{
			In: strings.NewReader(`{"action":"sleep","payload":"10ms"}`),
		},
		want: []*BulkResponse{
			{SeqID: 1},
			{SeqID: 2, IsFinal: true},
		},
	})
	tests.Add("schemas", tt{
		opts: &BulkOptions{
			In: strings.NewReader(`{"action":"schemas"}`),
		},
		want: []*BulkResponse{
			{
				SeqID: 1,
				// Following raw message is copied and pasted! (sorry!)
				Payload: json.RawMessage(`{
					"list": [
						"https://gobl.org/draft-0/bill/correction-options", "https://gobl.org/draft-0/bill/invoice", "https://gobl.org/draft-0/cal/date", "https://gobl.org/draft-0/cal/date-time", "https://gobl.org/draft-0/cal/period", "https://gobl.org/draft-0/cbc/code", "https://gobl.org/draft-0/cbc/code-definition", "https://gobl.org/draft-0/cbc/code-map", "https://gobl.org/draft-0/cbc/key", "https://gobl.org/draft-0/cbc/key-definition", "https://gobl.org/draft-0/cbc/meta", "https://gobl.org/draft-0/cbc/note", "https://gobl.org/draft-0/currency/amount", "https://gobl.org/draft-0/currency/code", "https://gobl.org/draft-0/currency/exchange-rate", "https://gobl.org/draft-0/dsig/digest", "https://gobl.org/draft-0/dsig/signature", "https://gobl.org/draft-0/envelope", "https://gobl.org/draft-0/head/header", "https://gobl.org/draft-0/head/link", "https://gobl.org/draft-0/head/stamp", "https://gobl.org/draft-0/i18n/string", "https://gobl.org/draft-0/l10n/code", "https://gobl.org/draft-0/l10n/country-code", "https://gobl.org/draft-0/note/message", "https://gobl.org/draft-0/num/amount", "https://gobl.org/draft-0/num/percentage", "https://gobl.org/draft-0/org/address", "https://gobl.org/draft-0/org/coordinates", "https://gobl.org/draft-0/org/email", "https://gobl.org/draft-0/org/identity", "https://gobl.org/draft-0/org/image", "https://gobl.org/draft-0/org/inbox", "https://gobl.org/draft-0/org/item", "https://gobl.org/draft-0/org/name", "https://gobl.org/draft-0/org/party", "https://gobl.org/draft-0/org/person", "https://gobl.org/draft-0/org/registration", "https://gobl.org/draft-0/org/telephone", "https://gobl.org/draft-0/org/unit", "https://gobl.org/draft-0/org/website", "https://gobl.org/draft-0/pay/advance", "https://gobl.org/draft-0/pay/instructions", "https://gobl.org/draft-0/pay/terms", "https://gobl.org/draft-0/regimes/mx/food-vouchers", "https://gobl.org/draft-0/regimes/mx/fuel-account-balance", "https://gobl.org/draft-0/schema/object", "https://gobl.org/draft-0/tax/extensions", "https://gobl.org/draft-0/tax/identity", "https://gobl.org/draft-0/tax/regime", "https://gobl.org/draft-0/tax/set", "https://gobl.org/draft-0/tax/total"
					]
				}`),
				IsFinal: false,
			},
			{SeqID: 2, IsFinal: true},
		},
	})
	tests.Add("schema", func(_ *testing.T) interface{} {
		return tt{
			opts: &BulkOptions{
				In: strings.NewReader(`{"action":"schema","payload":{"path":"head/stamp"}}`),
			},
			want: []*BulkResponse{
				{
					SeqID: 1,
					// Following raw message is copied and pasted on failures! (sorry!)
					Payload: json.RawMessage(`{
						"$schema": "https://json-schema.org/draft/2020-12/schema",
						"$id": "https://gobl.org/draft-0/head/stamp",
						"$ref": "#/$defs/Stamp",
						"$defs": {
						  "Stamp": {
							"properties": {
							  "prv": {
								"$ref": "https://gobl.org/draft-0/cbc/key",
								"title": "Provider",
								"description": "Identity of the agency used to create the stamp usually defined by each region."
							  },
							  "val": {
								"type": "string",
								"title": "Value",
								"description": "The serialized stamp value generated for or by the external agency"
							  }
							},
							"type": "object",
							"required": [
							  "prv",
							  "val"
							],
							"description": "Stamp defines an official seal of approval from a third party like a governmental agency or intermediary and should thus be included in any official envelopes."
						  }
						}
					  }`),
					IsFinal: false,
				},
				{SeqID: 2, IsFinal: true},
			},
		}
	})
	tests.Add("regime", func(_ *testing.T) interface{} {
		return tt{
			opts: &BulkOptions{
				In: strings.NewReader(`{"action":"regime","payload":{"code":"es"}}`),
			},
			want: []*BulkResponse{
				{
					SeqID: 1,
					// A small sample from the Spanish regime
					Payload: json.RawMessage(`{
						"$schema": "https://gobl.org/draft-0/tax/regime",
						"name": {
							"en": "Spain",
							"es": "España"
						}
					  }`),
					IsFinal: false,
				},
				{SeqID: 2, IsFinal: true},
			},
		}
	})
	tests.Add("out of order", tt{
		opts: &BulkOptions{
			In: strings.NewReader(`{"action":"sleep","payload":"100ms","req_id":"sleep"} {"action":"ping","req_id":"ping"}`),
		},
		want: []*BulkResponse{
			{SeqID: 2, ReqID: "ping"},
			{SeqID: 1, ReqID: "sleep"},
			{SeqID: 3, IsFinal: true},
		},
	})
	tests.Add("sign, explicit key given", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/nosig.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "sign",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data":       base64.StdEncoding.EncodeToString(payload),
				"privatekey": privateKey,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In: bytes.NewReader(req),
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})
	tests.Add("sign, default key", func(t *testing.T) interface{} {
		payload, err := os.ReadFile("testdata/nosig.json")
		if err != nil {
			t.Fatal(err)
		}
		req, err := json.Marshal(map[string]interface{}{
			"action": "sign",
			"req_id": "asdf",
			"payload": map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(payload),
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return tt{
			opts: &BulkOptions{
				In:                bytes.NewReader(req),
				DefaultPrivateKey: privateKey,
			},
			want: []*BulkResponse{
				{
					ReqID:   "asdf",
					SeqID:   1,
					IsFinal: false,
				},
				{
					SeqID:   2,
					IsFinal: true,
				},
			},
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		ch := Bulk(context.Background(), tt.opts)
		results := []*BulkResponse{}
		for res := range ch {
			results = append(results, res)
		}
		if d := cmp.Diff(tt.want, results, cmpopts.IgnoreFields(BulkResponse{}, "Payload", "Error")); d != "" {
			t.Error(d)
		}
		for i, row := range results {
			if tt.want[i].Payload != nil {
				var got map[string]interface{}
				if err := json.Unmarshal(row.Payload, &got); err != nil {
					t.Errorf("row %d: %v", i, err)
					continue
				}
				var want map[string]interface{}
				if err := json.Unmarshal(tt.want[i].Payload, &want); err != nil {
					t.Errorf("row %d: %v", i, err)
					continue
				}
				assert.Subset(t, got, want)
			}
			// Errors can get very complicated, so we resort to JSON comparisons
			if tt.want[i].Error != nil {
				got, err := json.Marshal(row.Error)
				if err != nil {
					t.Errorf("row %d: %v", i, err)
					continue
				}
				want, err := json.Marshal(tt.want[i].Error)
				if err != nil {
					t.Errorf("row %d: %v", i, err)
					continue
				}
				assert.JSONEq(t, string(got), string(want))
			}
		}
	})
}
