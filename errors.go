package gobl

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/validation"
)

// An Error provides a structure to better be able to make error comparisons.
// The contents can also be serialised as JSON ready to send to a client
// if needed, see [MarshalJSON] method.
type Error struct {
	key   cbc.Key
	cause error // the underlying error
}

// FieldErrors is a map of field names to errors, which is provided when it
// was possible to determine that the error was related to a specific field.
type FieldErrors map[string]error

var (
	// ErrNoDocument is provided when the envelope does not contain a
	// document payload.
	ErrNoDocument = NewError("no-document")

	// ErrValidation is used when a document fails a validation request.
	ErrValidation = NewError("validation")

	// ErrCalculation wraps around errors that we're generated during a
	// call to perform calculations on a document.
	ErrCalculation = NewError("calculation")

	// ErrMarshal is provided when there has been a problem attempting to encode
	// or marshal an object, usually into JSON.
	ErrMarshal = NewError("marshal")

	// ErrUnmarshal is used when that has been a problem attempting to read the
	// source data.
	ErrUnmarshal = NewError("unmarshal")

	// ErrSignature identifies an issue related to signatures.
	ErrSignature = NewError("signature")

	// ErrDigest identifies an issue related to the digest.
	ErrDigest = NewError("digest")

	// ErrInternal is a "catch-all" for errors that are not expected.
	ErrInternal = NewError("internal")

	// ErrUnknownSchema is provided when we attempt to determine the schema for an object
	// or from an ID and cannot find a match.
	ErrUnknownSchema = NewError("unknown-schema")
)

// NewError provides a new error with a code that is meant to provide
// a context.
func NewError(key cbc.Key) *Error {
	return &Error{key: key}
}

// wrapError is used to ensure that errors are wrapped around the GOBL standard
// error so they can be output in a consistent manner.
func wrapError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*Error); ok {
		return err // nothing to do
	}
	if errors.Is(err, schema.ErrUnknownSchema) {
		return ErrUnknownSchema
	}
	switch err.(type) {
	case validation.Errors:
		return ErrValidation.WithCause(err)
	}
	return ErrInternal.WithCause(err)
}

// Error provides a string representation of the error.
func (e *Error) Error() string {
	if e.cause != nil {
		msg := e.cause.Error()
		if e.Fields() != nil {
			msg = fmt.Sprintf("(%s).", msg)
		}
		return fmt.Sprintf("%s: %s", e.key, msg)
	}
	return e.key.String()
}

// WithCause is used to copy and add an underlying error to this one,
// unless the errors is already of type [*Error], in which case it will
// be returned as is.
func (e *Error) WithCause(err error) *Error {
	ne := e.copy()
	switch te := err.(type) {
	case *Error:
		return te
	case validation.Errors:
		ne.cause = fieldErrorsFromValidation(te)
	default:
		ne.cause = err
	}
	return ne
}

// WithReason returns the error with a specific reason.
func (e *Error) WithReason(msg string, a ...interface{}) *Error {
	ne := e.copy()
	ne.cause = fmt.Errorf(msg, a...)
	return ne
}

// Fields returns the errors that are associated with specific fields
// or nil if there are no field errors available.
func (e *Error) Fields() FieldErrors {
	if fe, ok := e.cause.(FieldErrors); ok {
		return fe
	}
	return nil
}

// Message returns a string representation of the error if it cannot
// be serialised as a map of field names to sub-errors, see also the
// [Fields] method.
func (e *Error) Message() string {
	if e.cause == nil {
		return ""
	}
	if e.Fields() != nil {
		return ""
	}
	return e.cause.Error()
}

func (e *Error) copy() *Error {
	ne := new(Error)
	*ne = *e
	return ne
}

// Is checks to see if the target error matches the current error or
// part of the chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.cause, target)
	}
	return e.key == t.key
}

// MarshalJSON converts the Error into a valid JSON, correctly
// handling mashaling of cause objects that might not have a
// valid MarhsalJSON method.
func (e *Error) MarshalJSON() ([]byte, error) {
	err := struct {
		Key     cbc.Key     `json:"key"`
		Fields  FieldErrors `json:"fields,omitempty"`
		Message string      `json:"message,omitempty"`
	}{
		Key:     e.key,
		Fields:  e.Fields(),
		Message: e.Message(),
	}
	return json.Marshal(err)
}

// Error returns the error string of FieldErrors. This is based on the
// implementation of [validation.Errors].
func (fe FieldErrors) Error() string {
	if len(fe) == 0 {
		return ""
	}

	keys := make([]string, len(fe))
	i := 0
	for key := range fe {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}
		switch errs := fe[key].(type) {
		case FieldErrors, validation.Errors:
			_, _ = fmt.Fprintf(&s, "%v: (%v)", key, errs)
		default:
			_, _ = fmt.Fprintf(&s, "%v: %v", key, fe[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
}

// MarshalJSON converts the FieldErrors into a valid JSON. Based on the
// implementation of [validation.Errors].
func (fe FieldErrors) MarshalJSON() ([]byte, error) {
	errs := map[string]any{}
	for key, err := range fe {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}

func fieldErrorsFromValidation(errs validation.Errors) FieldErrors {
	fe := make(FieldErrors)
	for key, err := range errs {
		fe[key] = err
	}
	return fe
}
