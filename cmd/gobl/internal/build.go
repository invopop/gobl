package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/internal/iotools"
)

// BuildOptions are the options to pass to the Build function.
type BuildOptions struct {
	Data io.Reader
}

// Build builds and validates a GOBL document from opts.
func Build(ctx context.Context, opts BuildOptions) (*gobl.Envelope, error) {
	dec := yaml.NewDecoder(iotools.CancelableReader(ctx, opts.Data))
	var intermediate map[string]interface{}
	if err := dec.Decode(&intermediate); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
	if err := env.Insert(doc); err != nil {
		return err
	}
	return nil
}

func extractDoc(env *gobl.Envelope) (gobl.Document, error) {
	switch env.Head.Type {
	case bill.InvoiceType:
		doc := new(bill.Invoice)
		err := env.Extract(doc)
		return doc, err
	default:
		return nil, fmt.Errorf("unrecognized document type: %s", env.Head.Type)
	}
}
