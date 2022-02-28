package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_serve_build(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		err  string
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
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`{"data":{}}`))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			err: `code=422, message=no document included`,
		},
		{
			name: "success",
			req: func() *http.Request {
				data, err := ioutil.ReadFile("testdata/success.json")
				if err != nil {
					t.Fatal(err)
				}
				body, err := json.Marshal(map[string]interface{}{
					"data": json.RawMessage(data),
				})
				if err != nil {
					t.Fatal(err)
				}
				req, _ := http.NewRequest(http.MethodPost, "/build", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
		},
		{
			name: "invalid data",
			req: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "/build", strings.NewReader(`{"data":"not an object"}`))
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
			err: "code=400, message=yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `not an ...` into map[string]interface {}",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(tt.req, rec)

			err := serve().build()(c)
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
			err: "code=400, message=error unmarshaling JSON: json: cannot unmarshal string into Go value of type gobl.Envelope",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(tt.req, rec)

			err := serve().verify()(c)
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
