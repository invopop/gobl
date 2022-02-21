package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/imdario/mergo"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/iotools"
)

// BuildOptions are the options to pass to the Build function.
type BuildOptions struct {
	Data      io.Reader
	SetYAML   map[string]string
	SetString map[string]string
	SetFile   map[string]string
}

// Build builds and validates a GOBL document from opts.
func Build(ctx context.Context, opts BuildOptions) (*gobl.Envelope, error) {
	values, err := parseSets(opts)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(iotools.CancelableReader(ctx, opts.Data))
	var intermediate map[string]interface{}
	if err := dec.Decode(&intermediate); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := mergo.Merge(&intermediate, values, mergo.WithOverride); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	encoded, err := json.Marshal(intermediate)
	if err != nil {
		return nil, err
	}
	env := new(gobl.Envelope)
	if err := json.Unmarshal(encoded, env); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = reInsertDoc(env)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return env, nil
}

func reInsertDoc(env *gobl.Envelope) error {
	if env.Document == nil {
		return errors.New("no document included")
	}
	doc, err := extractDoc(env)
	if err != nil {
		return err
	}
	return env.Insert(doc)
}

func extractDoc(env *gobl.Envelope) (interface{}, error) {
	if env.Document == nil {
		return nil, errors.New("no document found")
	}
	if env.Document.Schema == "" {
		return nil, errors.New("missing document schema")
	}
	doc := env.Document.Schema.Interface()
	if doc == nil {
		return nil, fmt.Errorf("unrecognized document schema %q", env.Document.Schema)
	}
	err := env.Document.Extract(doc)
	return doc, err
}

func parseSets(opts BuildOptions) (map[string]interface{}, error) {
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
			return nil, err
		}
		if err := setValue(&values, k, parsed); err != nil {
			return nil, err
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
