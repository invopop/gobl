package gobl_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/invopop/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
)

var testKey = dsig.NewES256Key()

const testMessageContent = "This is test content."

func ExampleNewEnvelope_complete() {
	// Prepare a new Envelope with a region
	env := gobl.NewEnvelope()
	env.Head.UUID = uuid.MustParse("871c1e6a-8b5c-11ec-af5f-3e7e00ce5635")

	// Prepare a payload and insert
	msg := &note.Message{
		Content: "sample message content",
	}
	msg.SetUUID(uuid.MustParse("e8c70516-0098-11ef-92c8-0242ac120002"))
	if err := env.Insert(msg); err != nil {
		panic(err.Error())
	}
	if err := env.Validate(); err != nil {
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
	// 			"val": "6854b999501883c478f0dbcb929ea1cb33e0e738fd0e74ac8194d1e5b7991980"
	// 		}
	// 	},
	// 	"doc": {
	// 		"$schema": "https://gobl.org/draft-0/note/message",
	//		"uuid": "e8c70516-0098-11ef-92c8-0242ac120002",
	// 		"content": "sample message content"
	// 	}
	// }
}

func TestEnvelop(t *testing.T) {
	msg := &note.Message{
		Content: testMessageContent,
	}
	msg.UUID = uuid.MustParse("871c1e6a-8b5c-11ec-af5f-3e7e00ce5635")
	e, err := gobl.Envelop(msg)
	require.NoError(t, err)
	if assert.NotNil(t, e) {
		assert.Equal(t, "cf75a55f8f00e57201685aebfa5765c908c1d22520858024610bbc2f6a494824", e.Head.Digest.Value)
	}
}

func TestEnvelopeDocument(t *testing.T) {
	m := testNoteExample()

	e := gobl.NewEnvelope()
	if assert.NotNil(t, e.Head) {
		assert.NotEmpty(t, e.Head.UUID, "empty header uuid")
	}
	assert.NotNil(t, e.Document)

	if err := e.Insert(m); err != nil {
		t.Errorf("failed to insert payload: %v", err)
		return
	}

	if assert.NotNil(t, e.Head.Digest) {
		assert.Equal(t, e.Head.Digest.Algorithm, dsig.DigestSHA256, "unexpected digest algorithm")
		assert.Equal(t, "54eb5ac433e82575b554dc21a8e53b291479dab188dffaabc97e8141d1cdfc65", e.Head.Digest.Value, "digest should be the same")
	}

	assert.Empty(t, e.Signatures)
	assert.NoError(t, e.Sign(testKey), "signing envelope")
	assert.NotEmpty(t, e.Signatures, "expected a signature")

	assert.NoError(t, e.Validate(), "did not expect validation error")

	nm, ok := e.Extract().(*note.Message)
	require.True(t, ok, "unrecognized content")
	assert.Equal(t, m.Content, nm.Content, "content mismatch")
}

func TestEnvelopeExtract(t *testing.T) {
	e := &gobl.Envelope{}
	obj := e.Extract()
	assert.Nil(t, obj)
}

func TestEnvelopeInsert(t *testing.T) {
	m := testNoteExample()

	t.Run("missing head", func(t *testing.T) {
		e := new(gobl.Envelope)
		err := e.Insert(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing head")
	})

	t.Run("no document", func(t *testing.T) {
		e := gobl.NewEnvelope()
		err := e.Insert(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no-document")
	})
}

func TestEnvelopeCalculate(t *testing.T) {
	m := testNoteExample()

	t.Run("basics", func(t *testing.T) {
		e := gobl.NewEnvelope()
		require.NoError(t, e.Insert(m))
		err := e.Calculate()
		assert.NoError(t, err)
	})

	t.Run("handle stamps", func(t *testing.T) {
		e := gobl.NewEnvelope()
		require.NoError(t, e.Insert(m))
		e.Head.AddStamp(&head.Stamp{Provider: cbc.Key("test"), Value: "test"})
		err := e.Calculate()
		assert.NoError(t, err)
		require.NoError(t, e.Sign(testKey))
		assert.NotEmpty(t, e.Head.Stamps)

		// remove signatures
		e.Signatures = nil
		err = e.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stamps: must be blank.")
		err = e.Calculate()
		assert.NoError(t, err)
		assert.Len(t, e.Head.Stamps, 1)
		/*
			// Removed for now as we prefer to just validate.
			assert.Empty(t, e.Head.Stamps)
			err = e.Validate()
			assert.NoError(t, err)
		*/
	})
}

func TestEnvelopeComplete(t *testing.T) {
	e := new(gobl.Envelope)

	data, err := os.ReadFile("./regimes/es/examples/invoice-es-es.env.yaml")
	require.NoError(t, err)
	err = yaml.Unmarshal(data, e)
	require.NoError(t, err)

	err = e.Calculate()
	require.NoError(t, err)

	inv, ok := e.Extract().(*bill.Invoice)
	require.True(t, ok)
	require.NoError(t, err)

	assert.Equal(t, "1210.00", inv.Totals.Payable.String())
}

func TestEnvelopeCompleteErrors(t *testing.T) {
	t.Run("missing document", func(t *testing.T) {
		e := new(gobl.Envelope)
		err := e.Calculate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, gobl.ErrNoDocument)
	})
	t.Run("missing document payload", func(t *testing.T) {
		e := gobl.NewEnvelope()
		err := e.Calculate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, gobl.ErrNoDocument)
	})
}

func TestEnvelopeValidate(t *testing.T) {
	key := dsig.NewES256Key()
	tests := []struct {
		name string
		env  func() *gobl.Envelope
		want string
	}{
		{
			name: "empty envelope",
			env: func() *gobl.Envelope {
				return &gobl.Envelope{}
			},
			want: "validation: ($schema: cannot be blank; doc: cannot be blank; head: cannot be blank.).",
		},
		{
			name: "missing message body, draft",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				require.NoError(t, env.Insert(&note.Message{}))
				return env
			},
			want: "validation: (doc: (content: cannot be blank.).).",
		},
		{
			name: "missing sig, draft",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				require.NoError(t, env.Insert(&note.Message{Content: "foo"}))
				return env
			},
		},
		{
			name: "with sig, not draft",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				require.NoError(t, env.Insert(&note.Message{Content: "foo"}))
				assert.NoError(t, env.Sign(key))
				return env
			},
		},
		{
			name: "with sig, not draft, modified",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				require.NoError(t, env.Insert(&note.Message{Content: "foo"}))
				assert.NoError(t, env.Sign(key))
				msg := env.Extract().(*note.Message)
				msg.Content = "bar"
				return env
			},
			want: "digest: mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := tt.env()
			err := env.Validate()
			if tt.want == "" && err == nil {
				return
			}
			assert.EqualError(t, err, tt.want)
		})
	}
}

func TestEnvelopeSign(t *testing.T) {
	t.Run("will sign", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Foooo"}))
		err := env.Sign(testKey)
		assert.NoError(t, err)
		assert.Len(t, env.Signatures, 1)
	})
	t.Run("cannot sign invalid document", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{})) // missing msg content
		err := env.Sign(testKey)
		assert.ErrorContains(t, err, "validation: (doc: (content: cannot be blank.).).")
	})
	t.Run("sign valid document", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message"}))
		err := env.Sign(testKey)
		assert.NoError(t, err)
		assert.True(t, env.Signed())
	})
}
func TestEnvelopeCorrect(t *testing.T) {
	t.Run("correct invoice", func(t *testing.T) {
		env := gobl.NewEnvelope()

		data, err := os.ReadFile("./regimes/es/examples/invoice-es-es.env.yaml")
		require.NoError(t, err)
		err = yaml.Unmarshal(data, env)
		require.NoError(t, err)
		require.NoError(t, env.Calculate())

		_, err = env.Correct(
			bill.Corrective,
			bill.WithExtension(facturae.ExtKeyCorrection, "01"),
		)
		require.NoError(t, err)

		doc := env.Extract().(*bill.Invoice)
		assert.Equal(t, doc.Type, bill.InvoiceTypeStandard, "should not update in place")

		e2, err := env.Correct(
			bill.Corrective,
			bill.WithExtension(facturae.ExtKeyCorrection, "02"),
		)
		require.NoError(t, err)
		doc = e2.Extract().(*bill.Invoice)
		assert.Equal(t, doc.Type, bill.InvoiceTypeCorrective, "corrected")
	})
}

func TestEnvelopeReplicate(t *testing.T) {
	t.Run("replicate invoice", func(t *testing.T) {
		env := gobl.NewEnvelope()

		data, err := os.ReadFile("./regimes/es/examples/invoice-es-es.env.yaml")
		require.NoError(t, err)
		err = yaml.Unmarshal(data, env)
		require.NoError(t, err)
		require.NoError(t, env.Calculate())

		_, err = env.Replicate()
		require.NoError(t, err)

		doc := env.Extract().(*bill.Invoice)
		assert.Equal(t, "SAMPLE-001", doc.Code, "should not update in place")

		e2, err := env.Replicate()
		require.NoError(t, err)
		doc = e2.Extract().(*bill.Invoice)
		assert.Empty(t, doc.Code)
	})
}

func TestDocument(t *testing.T) {
	msg := testNoteExample()
	env := gobl.NewEnvelope()
	err := env.Insert(msg)
	require.NoError(t, err)
	doc := env.Document

	id := schema.Lookup(&note.Message{})
	assert.Contains(t, id.String(), "https://gobl.org/")
	assert.Contains(t, id.String(), "/note/message")

	dig := "54eb5ac433e82575b554dc21a8e53b291479dab188dffaabc97e8141d1cdfc65"
	assert.Equal(t, id, doc.Schema)
	sha, err := env.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)
	assert.Equal(t, doc.Instance(), msg)

	data, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.Equal(t, `{"$schema":"`+id.String()+`","uuid":"e8c70516-0098-11ef-92c8-0242ac120002","content":"This is test content."}`, string(data))

	doc = new(schema.Object)
	err = json.Unmarshal(data, doc)
	require.NoError(t, err)

	assert.Equal(t, doc.Schema, id)
	sha, err = env.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)

	obj, ok := doc.Instance().(*note.Message)
	assert.True(t, ok)
	assert.Equal(t, msg.Content, obj.Content)
}

func TestDocumentValidation(t *testing.T) {
	msg := &note.Message{}

	doc, err := schema.NewObject(msg)
	require.NoError(t, err)

	err = doc.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "content: cannot be blank")
	}

	doc = new(schema.Object)
	data, err := os.ReadFile("./regimes/es/examples/invoice-es-es.yaml")
	require.NoError(t, err)
	err = yaml.Unmarshal(data, doc)
	require.NoError(t, err)

	inv := doc.Instance().(*bill.Invoice)
	inv.Code = "" // blank, which will not be accepted if not a draft
	require.NoError(t, doc.Calculate())
	assert.NoError(t, doc.Validate())
	inv.IssueDate = cal.Date{}
	err = doc.Validate()
	if assert.Error(t, err) {
		// Double check to make sure validation working
		assert.Contains(t, err.Error(), "issue_date: required")
	}
}

func TestDocumentValidationOutput(t *testing.T) {
	msg := &note.Message{}

	doc, err := schema.NewObject(msg)
	require.NoError(t, err)

	err = doc.Validate()
	data, err := json.Marshal(err)
	require.NoError(t, err)
	assert.Equal(t, `{"content":"cannot be blank"}`, string(data))

	env := gobl.NewEnvelope()
	require.NoError(t, env.Insert(msg))
	err = env.Validate()
	data, err = json.Marshal(err)
	require.NoError(t, err)
	assert.Equal(t, `{"key":"validation","fields":{"doc":{"content":"cannot be blank"}}}`, string(data))
}

func TestEnvelopeVerify(t *testing.T) {
	t.Run("invalid situations", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message"}))
		err := env.Verify()
		assert.ErrorContains(t, err, "no signatures to verify")
	})

	t.Run("valid signature", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message"}))
		err := env.Sign(testKey)
		require.NoError(t, err)
		err = env.Verify()
		assert.NoError(t, err)
		err = env.Verify(testKey.Public())
		assert.NoError(t, err)
		rk := dsig.NewES256Key()
		err = env.Verify(rk.Public(), testKey.Public())
		assert.NoError(t, err)
		err = env.Verify(rk.Public())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signatures: (0: no key match found.)")
	})

	t.Run("changes", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message"}))
		err := env.Sign(testKey)
		require.NoError(t, err)
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message 2"}))
		err = env.Verify()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signatures: (0: header mismatch.)")
		err = env.Verify(testKey.Public())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signatures: (0: header mismatch.)")

		rk := dsig.NewES256Key()
		err = env.Verify(rk.Public())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signatures: (0: no key match found.)")
	})

}

func testNoteExample() *note.Message {
	m := new(note.Message)
	m.Content = testMessageContent
	m.UUID = uuid.MustParse("e8c70516-0098-11ef-92c8-0242ac120002")
	return m
}
