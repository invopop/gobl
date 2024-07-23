package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

const signingKeyText = `{"use":"sig","kty":"EC","kid":"b7cee60f-204e-438b-a88f-021d28af6991","crv":"P-256","alg":"ES256","x":"wLez6TfqNReD3FUUyVP4Q7HAGdokmAfE6LwfcM28DlQ","y":"CIxURqWtiFIu9TaatRa85NkNsw1LZHw_ZQ9A45GW_MU","d":"xNx9MxONcuLk8Ai6s2isqXMZaDi3HNGLkFX-qiNyyeo"}`

func Test_serve_build(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		err     string
		replace []testy.Replacement
	}{
		{
			name: "wrong content type",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", nil)
				req.Header.Set("Content-Type", "text/plain")
				return req
			}(),
			err: "code=415, message=Unsupported Media Type",
		},
		{
			name: "invalid json payload",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`invalid`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value`,
		},
		{
			name: "missing payload",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`{}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=400, message=no payload`,
		},
		{
			name: "missing doc",
			req: func() *http.Request {
				body, err := json.Marshal(map[string]interface{}{
					"data":       base64.StdEncoding.EncodeToString([]byte(`{"$schema":"https://gobl.org/draft-0/envelope"}`)),
					"privatekey": json.RawMessage(signingKeyText),
				})
				if err != nil {
					t.Fatal(err)
				}
				req, _ := http.NewRequest(http.MethodPost, "/build", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=422, message=no-document`,
		},
		{
			name: "invalid data",
			req: func() *http.Request {
				body, err := json.Marshal(map[string]interface{}{
					"data":       base64.StdEncoding.EncodeToString([]byte("not an object")),
					"privatekey": json.RawMessage(signingKeyText),
				})
				if err != nil {
					t.Fatal(err)
				}
				req, _ := http.NewRequest(http.MethodPost, "/build", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `not an ...` into map[string]interface {}",
		},
		{
			name: "template",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`{"template":"not an object","data":{}}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: "code=400, message=illegal base64 data at input byte 3, internal=illegal base64 data at input byte 3",
		},
		{
			name: "envelop success",
			req: func() *http.Request {
				data, err := os.ReadFile("testdata/message.json")
				if err != nil {
					t.Fatal(err)
				}
				body, _ := json.Marshal(map[string]interface{}{
					"data":       base64.StdEncoding.EncodeToString(data),
					"type":       "note.Message",
					"privatekey": json.RawMessage(signingKeyText),
				})

				req, _ := http.NewRequest(http.MethodPost, "/build", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			replace: []testy.Replacement{
				{
					Regexp:      regexp.MustCompile(`"uuid":.?"[^\"]+"`),
					Replacement: `"uuid":"00000000-0000-0000-0000-000000000000"`,
				},
			},
		},
		{
			name: "envelop success indented",
			req: func() *http.Request {
				data, err := os.ReadFile("testdata/message.json")
				if err != nil {
					t.Fatal(err)
				}
				body, _ := json.Marshal(map[string]interface{}{
					"data":       base64.StdEncoding.EncodeToString(data),
					"type":       "note.Message",
					"privatekey": json.RawMessage(signingKeyText),
				})

				req, _ := http.NewRequest(http.MethodPost, "/build?indent=true", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			replace: []testy.Replacement{
				{
					Regexp:      regexp.MustCompile(`"uuid":.?"[^\"]+"`),
					Replacement: `"uuid":"00000000-0000-0000-0000-000000000000"`,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(tt.req, rec)

			err := serve().build(c)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if tt.replace == nil {
				tt.replace = make([]testy.Replacement, 0)
			}
			if d := testy.DiffHTTPResponse(testy.Snapshot(t), rec.Result(), tt.replace...); d != nil {
				t.Error(d)
			}
		})
	}
}

func Test_serve_verify(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		err  string
	}{
		{
			name: "wrong content type",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/verify", nil)
				req.Header.Set("Content-Type", "text/plain")
				return req
			}(),
			err: "code=415, message=Unsupported Media Type",
		},
		{
			name: "invalid json payload",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/verify", strings.NewReader(`invalid`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value`,
		},
		{
			name: "validation failure",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/verify", strings.NewReader(`{}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=422, message=$schema: cannot be blank; doc: cannot be blank; head: cannot be blank.`,
		},
		{
			name: "invalid data",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`{"data":"not an object"}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: "code=400, message=illegal base64 data at input byte 3, internal=illegal base64 data at input byte 3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(tt.req, rec)

			err := serve().verify(c)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if d := testy.DiffHTTPResponse(testy.Snapshot(t), rec.Result()); d != nil {
				t.Error(d)
			}
		})
	}
}

func Test_serve_keygen(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/key", nil)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = serve().keygen(c)
	if err != nil {
		t.Fatal(err)
	}

	if d := testy.DiffHTTPResponse(testy.Snapshot(t), rec.Result(), jwkREs...); d != nil {
		t.Error(d)
	}
}

func Test_serve_bulk(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/bulk", strings.NewReader(`{"action":"oink"}`))
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = serve().bulk(c)
	if err != nil {
		t.Fatal(err)
	}

	if d := testy.DiffHTTPResponse(testy.Snapshot(t), rec.Result(), jwkREs...); d != nil {
		t.Error(d)
	}
}
