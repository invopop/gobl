package gobl

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/validation"
)

// Error provides a structure to better be able to make error comparisons.
// The contents can also be serialised as JSON ready to send to a client
// if needed.
type Error struct {
	Key   cbc.Key `json:"key"`
	Cause error   `json:"cause"`
}

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
	return &Error{Key: key}
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
	if _, ok := err.(validation.Errors); ok {
		return ErrValidation.WithCause(err)
	}
	return ErrInternal.WithCause(err)
}

// Error provides a string representation of the error.
func (e *Error) Error() string {
	if e.Cause != nil {
		cause := e.Cause.Error()
		if reflect.TypeOf(e.Cause).Kind() == reflect.Map {
			cause = fmt.Sprintf("(%s).", cause)
		}
		return fmt.Sprintf("%s: %s", e.Key.String(), cause)
	}
	return e.Key.String()
}

// WithCause is used to copy and add an underlying error to this one.
func (e *Error) WithCause(err error) *Error {
	ne := e.copy()
	ne.Cause = err
	return ne
}

// WithReason returns an error with a reason attached.
func (e *Error) WithReason(msg string) *Error {
	ne := e.copy()
	ne.Cause = errors.New(msg)
	return ne
}

// WithErrorf wraps around the `fmt.Errorf` call to provide a more meaningful
// error in the context.
func (e *Error) WithErrorf(format string, a ...interface{}) *Error {
	ne := e.copy()
	ne.Cause = fmt.Errorf(format, a...)
	return ne
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
		return errors.Is(e.Cause, target)
	}
	return e.Key == t.Key
}
