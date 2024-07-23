package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/schema"
)

// BulkRequest represents a single request in the stream of bulk requests.
type BulkRequest struct {
	// Action is the action to perform on the payload.
	Action string `json:"action"`
	// ReqID is an opaque request ID, which is returned with the associated
	// response.
	ReqID string `json:"req_id"`
	// Payload is the payload upon which to perform the action.
	Payload json.RawMessage `json:"payload"`
	// When true, responses are indented for easier human consumption
	Indent bool `json:"indent"`
}

// BulkResponse represents a single response in the stream of bulk responses.
type BulkResponse struct {
	// ReqID is an exact copy of the value provided in the request, if any.
	ReqID string `json:"req_id,omitempty"`
	// SeqID is the sequence ID of the request this response correspond to,
	// starting at 1.
	SeqID int64 `json:"seq_id"`
	// Payload is the response payload.
	Payload json.RawMessage `json:"payload,omitempty"`
	// Error represents an error processing a request item.
	Error *Error `json:"error"`
	// IsFinal will be true once the end of the request input stream has been
	// reached, or an unrecoverable error has occurred.
	IsFinal bool `json:"is_final"`
}

// BulkOptions are the options used for processing a stream of bulk requests.
type BulkOptions struct {
	// In is the input stream of requests
	In io.Reader
	// DefaultPrivateKey is the default private key to use with sign requests
	DefaultPrivateKey *dsig.PrivateKey
}

// VerifyRequest is the payload for a verification request.
type VerifyRequest struct {
	Data      []byte          `json:"data"`
	PublicKey *dsig.PublicKey `json:"publickey"`
}

// VerifyResponse is the response to a verification request.
type VerifyResponse struct {
	OK bool `json:"ok"`
}

// ValidateResponse is the response to a validate request.
type ValidateResponse struct {
	OK bool `json:"ok"`
}

// BuildRequest is the payload for a build request.
type BuildRequest struct {
	Template []byte `json:"template"`
	Data     []byte `json:"data"`
	DocType  string `json:"type"`
	Envelop  bool   `json:"envelop"`
}

// SignRequest is the payload for a sign request.
type SignRequest struct {
	Template   []byte           `json:"template"`
	Data       []byte           `json:"data"`
	PrivateKey *dsig.PrivateKey `json:"privatekey"`
	DocType    string           `json:"type"`
	Envelop    bool             `json:"envelop"`
}

// ValidateRequest is the payload for a validate request.
type ValidateRequest struct {
	Data []byte `json:"data"`
}

// KeygenResponse is the payload for a key generation response.
type KeygenResponse struct {
	Private *dsig.PrivateKey `json:"private"`
	Public  *dsig.PublicKey  `json:"public"`
}

// CorrectRequest is the payload used to generate a corrected document.
// If the schema option is true, the options data is ignored.
type CorrectRequest struct {
	Data    []byte `json:"data"`
	Options []byte `json:"options"`
	Schema  bool   `json:"schema"`
}

// ReplicateRequest defines the payload used to generate a replicated document.
type ReplicateRequest struct {
	Data []byte `json:"data"`
}

// SchemaRequest defines a body used to request a specific JSON schema
type SchemaRequest struct {
	Path string `json:"path"`
}

// RegimeRequest defines a body used to request the definition of a Tax Regime.
type RegimeRequest struct {
	Code string `json:"code"`
}

// Bulk processes a stream of bulk requests.
func Bulk(ctx context.Context, opts *BulkOptions) <-chan *BulkResponse {
	dec := json.NewDecoder(opts.In)
	resCh := make(chan *BulkResponse, 1)
	wg := &sync.WaitGroup{}
	go func() {
		var seq int64
		defer close(resCh)
		for {
			seq := atomic.AddInt64(&seq, 1)
			var req BulkRequest
			err := dec.Decode(&req)
			if err != nil {
				wg.Wait()
				res := &BulkResponse{
					ReqID:   req.ReqID,
					SeqID:   seq,
					IsFinal: true,
				}
				if err != io.EOF {
					res.Error = wrapError(StatusUnprocessableEntity, err)
				}
				resCh <- res
				return
			}
			wg.Add(1)
			go func() {
				resCh <- processRequest(ctx, req, seq, opts)
				wg.Done()
			}()
		}
	}()
	return resCh
}

func processRequest(ctx context.Context, req BulkRequest, seq int64, bulkOpts *BulkOptions) *BulkResponse { //nolint:gocyclo
	marshal := json.Marshal
	if req.Indent {
		marshal = func(i interface{}) ([]byte, error) {
			return json.MarshalIndent(i, "", "\t")
		}
	}
	res := &BulkResponse{
		ReqID: req.ReqID,
		SeqID: seq,
	}
	switch req.Action {
	case "verify":
		vrfy := &VerifyRequest{}
		if err := json.Unmarshal(req.Payload, vrfy); err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		err := Verify(ctx, bytes.NewReader(vrfy.Data), vrfy.PublicKey)
		if err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(VerifyResponse{OK: true})
	case "validate":
		valReq := &ValidateRequest{}
		if err := json.Unmarshal(req.Payload, valReq); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		data := bytes.NewReader(valReq.Data)
		if err := Validate(ctx, data); err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(ValidateResponse{OK: true})
	case "build":
		bld := &BuildRequest{}
		if err := json.Unmarshal(req.Payload, bld); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		opts := &BuildOptions{
			ParseOptions: &ParseOptions{
				DocType: bld.DocType,
				Input:   bytes.NewReader(bld.Data),
				Envelop: bld.Envelop,
			},
		}
		if len(bld.Template) > 0 {
			opts.Template = bytes.NewReader(bld.Template)
		}
		env, err := Build(ctx, opts)
		if err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(env)
	case "sign":
		bld := &SignRequest{}
		if err := json.Unmarshal(req.Payload, bld); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		opts := &SignOptions{
			ParseOptions: &ParseOptions{
				DocType: bld.DocType,
				Input:   bytes.NewReader(bld.Data),
			},
			PrivateKey: bld.PrivateKey,
		}
		if len(bld.Template) > 0 {
			opts.Template = bytes.NewReader(bld.Template)
		}
		if opts.PrivateKey == nil {
			opts.PrivateKey = bulkOpts.DefaultPrivateKey
		}
		env, err := Sign(ctx, opts)
		if err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(env)
	case "correct":
		bld := &CorrectRequest{}
		if err := json.Unmarshal(req.Payload, bld); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		opts := &CorrectOptions{
			ParseOptions: &ParseOptions{
				Input: bytes.NewReader(bld.Data),
			},
			OptionsSchema: bld.Schema,
			Data:          bld.Options,
		}
		env, err := Correct(ctx, opts)
		if err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(env)
	case "replicate":
		rep := &ReplicateRequest{}
		if err := json.Unmarshal(req.Payload, rep); err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		opts := &ReplicateOptions{
			ParseOptions: &ParseOptions{
				Input: bytes.NewReader(rep.Data),
			},
		}
		env, err := Replicate(ctx, opts)
		if err != nil {
			res.Error = wrapError(StatusUnprocessableEntity, err)
			return res
		}
		res.Payload, _ = marshal(env)
	case "keygen":
		key := dsig.NewES256Key()

		res.Payload, _ = marshal(KeygenResponse{
			Private: key,
			Public:  key.Public(),
		})
	case "ping":
		res.Payload, _ = marshal(map[string]interface{}{
			"pong": true,
		})
	case "sleep":
		var delay string
		if err := json.Unmarshal(req.Payload, &delay); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		dur, err := time.ParseDuration(delay)
		if err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid duration: %w", err)
			return res
		}
		time.Sleep(dur)
		res.Payload, _ = marshal(map[string]interface{}{
			"sleep": "done",
		})

	case "schemas":
		list := schema.List()
		items := make([]string, len(list))
		for i, v := range list {
			items[i] = v.String()
		}
		// sorting makes comparisons easier
		sort.Strings(items)
		res.Payload, _ = marshal(map[string]interface{}{
			"list": items,
		})
	case "schema":
		sch := new(SchemaRequest)
		if err := json.Unmarshal(req.Payload, sch); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		ext := filepath.Ext(sch.Path)
		if ext == "" {
			sch.Path = sch.Path + ".json"
		}
		sch.Path = path.Join("schemas", sch.Path)
		data, err := data.Content.ReadFile(sch.Path)
		if err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid schema: %w", err)
			return res
		}
		res.Payload = data
	case "regime":
		reg := new(RegimeRequest)
		if err := json.Unmarshal(req.Payload, reg); err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid payload: %w", err)
			return res
		}
		p := path.Join("regimes", strings.ToLower(reg.Code)+".json")
		data, err := data.Content.ReadFile(p)
		if err != nil {
			res.Error = wrapErrorf(StatusUnprocessableEntity, "invalid regime: %w", err)
			return res
		}
		res.Payload = data
	default:
		res.Error = wrapErrorf(StatusBadRequest, "unrecognized action: '%s'", req.Action)
	}
	return res
}
