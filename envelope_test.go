package gobl_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/invopop/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl"
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
	// 			"val": "7d539c46ca03a4ecb1fcc4cb00d2ada34275708ee326caafee04d9dcfed862ee"
	// 		},
	//		"draft": true
	// 	},
	// 	"doc": {
	// 		"$schema": "https://gobl.org/draft-0/note/message",
	// 		"content": "sample message content"
	// 	}
	// }
}

func TestEnvelop(t *testing.T) {
	msg := &note.Message{Content: testMessageContent}
	e, err := gobl.Envelop(msg)
	require.NoError(t, err)
	if assert.NotNil(t, e) {
		assert.Equal(t, "c6a5148ce90f70c24ebfe6de1abed0d0aafde4323a9bcf47cc4a5d544af9ea19", e.Head.Digest.Value)
	}
}

func TestEnvelopeDocument(t *testing.T) {
	m := new(note.Message)
	m.Content = testMessageContent

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
		assert.Equal(t, "c6a5148ce90f70c24ebfe6de1abed0d0aafde4323a9bcf47cc4a5d544af9ea19", e.Head.Digest.Value, "digest should be the same")
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
	m := new(note.Message)
	m.Content = testMessageContent

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
	m := new(note.Message)
	m.Content = testMessageContent

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
		assert.NotEmpty(t, e.Head.Stamps)
		e.Head.Draft = true
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
		assert.ErrorIs(t, err, schema.ErrNoDocument)
	})
	t.Run("missing document payload", func(t *testing.T) {
		e := gobl.NewEnvelope()
		err := e.Calculate()
		assert.Error(t, err)
		assert.ErrorIs(t, err, schema.ErrNoDocument)
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
			want: "$schema: cannot be blank; doc: cannot be blank; head: cannot be blank.",
		},
		{
			name: "missing message body, draft",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				env.Head.Draft = true
				require.NoError(t, env.Insert(&note.Message{}))
				return env
			},
			want: "doc: (content: cannot be blank.).",
		},
		{
			name: "missing sig, draft",
			env: func() *gobl.Envelope {
				env := gobl.NewEnvelope()
				env.Head.Draft = true
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
			want: "document: digest mismatch",
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
	t.Run("will sign draft", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Foooo"}))
		assert.True(t, env.Head.Draft)
		err := env.Sign(testKey)
		assert.NoError(t, err)
		assert.False(t, env.Head.Draft)
	})
	t.Run("cannot sign invalid document", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{})) // missing msg content
		err := env.Sign(testKey)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation: doc: (content: cannot be blank.).")
	})
	t.Run("sign valid document", func(t *testing.T) {
		env := gobl.NewEnvelope()
		require.NoError(t, env.Insert(&note.Message{Content: "Test Message"}))
		err := env.Sign(testKey)
		assert.NoError(t, err)
		assert.Len(t, env.Signatures, 1)
	})
}
func TestEnvelopeCorrect(t *testing.T) {
	t.Run("correct invoice", func(t *testing.T) {
		env := gobl.NewEnvelope()

		data, err := ioutil.ReadFile("./regimes/es/examples/invoice-es-es.env.yaml")
		require.NoError(t, err)
		err = yaml.Unmarshal(data, env)
		require.NoError(t, err)
		require.NoError(t, env.Calculate())

		_, err = env.Correct()
		require.NoError(t, err)

		doc := env.Extract().(*bill.Invoice)
		assert.Equal(t, doc.Type, bill.InvoiceTypeStandard, "no change")

		e2, err := env.Correct()
		require.NoError(t, err)
		doc = e2.Extract().(*bill.Invoice)
		assert.Equal(t, doc.Type, bill.InvoiceTypeCorrective, "corrected")
	})
}

func TestDocument(t *testing.T) {
	msg := &note.Message{
		Content: "test message",
	}
	env := gobl.NewEnvelope()
	err := env.Insert(msg)
	require.NoError(t, err)
	doc := env.Document

	id := schema.Lookup(&note.Message{})
	assert.Contains(t, id.String(), "https://gobl.org/")
	assert.Contains(t, id.String(), "/note/message")

	dig := "82a5cddc56f069ff17705f310161dd17cd8b00d94728e6be3fafdad980522a27"
	assert.Equal(t, id, doc.Schema())
	sha, err := env.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)
	assert.Equal(t, doc.Instance(), msg)

	data, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.Equal(t, `{"$schema":"`+id.String()+`","content":"test message"}`, string(data))
	digest := dsig.NewSHA256Digest(data) // this works as the JSON is very simple!
	assert.Equal(t, dig, digest.Value)

	doc = new(schema.Document)
	err = json.Unmarshal(data, doc)
	require.NoError(t, err)

	assert.Equal(t, doc.Schema(), id)
	sha, err = env.Digest()
	require.NoError(t, err)
	assert.Equal(t, dig, sha.Value)

	obj, ok := doc.Instance().(*note.Message)
	assert.True(t, ok)
	assert.Equal(t, msg.Content, obj.Content)
}

func TestDocumentValidation(t *testing.T) {
	msg := &note.Message{}

	doc, err := schema.NewDocument(msg)
	require.NoError(t, err)

	err = doc.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "content: cannot be blank")
	}

	doc = new(schema.Document)
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
