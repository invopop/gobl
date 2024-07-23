package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/imdario/mergo"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/iotools"
	"github.com/invopop/gobl/schema"
	"gopkg.in/yaml.v3"
)

// ParseOptions are the options used for parsing incoming GOBL data.
type ParseOptions struct {
	Template  io.Reader
	Input     io.Reader
	DocType   string
	SetYAML   map[string]string
	SetString map[string]string
	SetFile   map[string]string

	// When set to `true`, the parsed data is wrapped in an envelope (if needed).
	Envelop bool
}

// decodeInto unmarshals in as YAML, then merges it into dest.
func decodeInto(ctx context.Context, dest *map[string]interface{}, in io.Reader) error {
	var intermediate map[string]interface{}
	dec := yaml.NewDecoder(iotools.CancelableReader(ctx, in))
	if err := dec.Decode(&intermediate); err != nil {
		return wrapError(StatusBadRequest, err)
	}
	if err := mergo.Merge(dest, intermediate, mergo.WithOverride); err != nil {
		return wrapError(StatusUnprocessableEntity, err)
	}
	return nil
}

func parseGOBLData(ctx context.Context, opts *ParseOptions) (interface{}, error) {
	var intermediate map[string]interface{}

	values, err := parseSets(opts)
	if err != nil {
		return nil, err
	}

	if opts.Template != nil {
		if err = decodeInto(ctx, &intermediate, opts.Template); err != nil {
			return nil, err
		}
	}

	if err = decodeInto(ctx, &intermediate, opts.Input); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&intermediate, values, mergo.WithOverride); err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	if opts.DocType != "" {
		// We want to infer if the incoming data is an envelope, but don't want
		// to call `gobl.Parse` just yet (because we want to first merge based
		// on `opts.DocType`. Thus, we (somewhat hacky) infer based on the
		// intermediate map ourselves.
		schemaDataFunc := docSchemaData
		if schema, ok := intermediate["$schema"].(string); ok && schema == string(gobl.EnvelopeSchema) {
			schemaDataFunc = docInEnvelopeSchemaData
		}
		schema := FindType(opts.DocType)
		if schema == "" {
			return nil, wrapError(StatusBadRequest, fmt.Errorf("unrecognized doc type: %q", opts.DocType))
		}
		if err := mergo.Merge(&intermediate, schemaDataFunc(schema)); err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}
	}

	// Encode intermediate to JSON for usage with `gobl.Parse`.
	intermediateJSON, err := json.Marshal(intermediate)
	if err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	// Parse the JSON encoded intermediate, so we can figure out if the incoming data
	// is already an envelope.
	obj, err := gobl.Parse(intermediateJSON)
	if err != nil {
		return nil, wrapError(StatusBadRequest, err)
	}

	// If the incoming data was parsed as an envelope, we can simply return
	// without the need for wrapping.
	env, isEnvelope := obj.(*gobl.Envelope)
	if isEnvelope {
		return env, nil
	}

	var doc *schema.Object
	if d, ok := obj.(*schema.Object); ok {
		doc = d
	} else {
		var err error
		doc, err = schema.NewObject(obj)
		if err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}
	}

	if !opts.Envelop {
		return doc, nil
	}

	// Wrap the parsed document in an envelope.
	// Note: We don't use `gobl.Envelop()`, because it (indirectly) calculates
	// the document, which we don't want to do here. Any calculations should
	// be incurred from the call site of this function.
	env = gobl.NewEnvelope()
	env.Document = doc
	// Set envelope as draft, so it can be rebuilt over time, and eventually
	// signed using the separate `sign` command.
	env.Head.Draft = true

	return env, nil
}

func docInEnvelopeSchemaData(schema schema.ID) map[string]interface{} {
	return map[string]interface{}{
		"doc": docSchemaData(schema),
	}
}

func docSchemaData(schema schema.ID) map[string]interface{} {
	return map[string]interface{}{
		"$schema": schema,
	}
}

func parseSets(opts *ParseOptions) (map[string]interface{}, error) {
	values := map[string]interface{}{}
	keys := make([]string, 0, len(opts.SetYAML))
	for k := range opts.SetYAML {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetYAML[k]
		var parsed interface{}
		if err := yaml.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}
		if err := setValue(&values, k, parsed); err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}
	}

	keys = make([]string, 0, len(opts.SetString))
	for k := range opts.SetString {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetString[k]
		if err := setValue(&values, k, v); err != nil {
			return nil, err
		}
	}

	keys = make([]string, 0, len(opts.SetFile))
	for k := range opts.SetFile {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := opts.SetFile[k]
		f, err := os.Open(v)
		if err != nil {
			return nil, err
		}
		defer f.Close() // nolint:errcheck
		dec := yaml.NewDecoder(f)
		var val interface{}
		if err := dec.Decode(&val); err != nil {
			return nil, err
		}
		if err := setValue(&values, k, val); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func setValue(values *map[string]interface{}, key string, value interface{}) error {
	key = strings.ReplaceAll(key, `\.`, "\x00")

	// If the key starts with '.', we treat that as the root of the
	// target object
	if key == "." {
		return mergo.Merge(values, value, mergo.WithOverride)
	}
	if len(key) > 1 && key[0] == '.' {
		key = key[1:]
	}

	for {
		i := strings.LastIndex(key, ".")
		if i == -1 {
			break
		}
		value = map[string]interface{}{
			strings.ReplaceAll(key[i+1:], "\x00", "."): value,
		}
		key = key[:i]
	}
	return mergo.Merge(values, map[string]interface{}{
		strings.ReplaceAll(key, "\x00", "."): value,
	}, mergo.WithOverride)
}
