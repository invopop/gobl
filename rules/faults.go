package rules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fault represents a single rule assertion failure identified by a code and path.
// Fault is *not* designed to be instantiated directly, and will be created as part
// of the validation processes from defined rules.
//
//nolint:errname
type Fault struct {
	path    string
	code    Code
	message string
}

func newFault(path string, id Code, message string) *Fault {
	return &Fault{path: path, code: id, message: message}
}

// Path returns the JSON Path (RFC 6901) location where this fault occurred.
// Returns "$" if the fault is at the root.
func (f *Fault) Path() string {
	return publicPath(f.path)
}

// Code returns the assertion code that produced this fault.
func (f *Fault) Code() Code {
	return f.code
}

// Message returns the human-readable message associated with this fault.
func (f *Fault) Message() string {
	return f.message
}

// Error implements the error interface.
func (f *Fault) Error() string {
	msg := f.message
	if f.path != "" {
		msg = "(" + publicPath(f.path) + ") " + msg
	}
	return fmt.Sprintf("[%s] %s", f.code, msg)
}

// MarshalJSON encodes the fault as a JSON object with path, code, and message fields.
func (f *Fault) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path    string `json:"path"`
		Code    Code   `json:"code"`
		Message string `json:"message"`
	}{publicPath(f.path), f.code, f.message})
}

// Faults is the interface for a collection of validation faults.
// A nil value indicates no faults were found.
type Faults interface {
	error
	HasPath(path string) bool
	HasCode(code Code) bool
	// Len returns the number of faults in the collection.
	Len() int
	// First returns the first fault in the collection.
	First() *Fault
	// Last returns the last fault in the collection.
	Last() *Fault
	// At returns the fault at position i.
	At(i int) *Fault
	// List returns the underlying slice of faults.
	List() []*Fault
}

// faultList is the concrete slice-based implementation of Faults.
//
//nolint:errname
type faultList []*Fault

// newFaults wraps a set of faults in the Faults interface, returning nil when empty.
func newFaults(faults ...*Fault) Faults {
	if len(faults) == 0 {
		return nil
	}
	return faultList(faults)
}

// Error implements the error interface.
func (fs faultList) Error() string {
	var b strings.Builder
	for i, f := range fs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(f.Error())
	}
	return b.String()
}

// MarshalJSON encodes the faults as a JSON array.
func (fs faultList) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*Fault(fs))
}

// HasPath reports whether any fault has exactly the given JSON Path.
func (fs faultList) HasPath(path string) bool {
	for _, f := range fs {
		if publicPath(f.path) == path {
			return true
		}
	}
	return false
}

// HasCode reports whether any fault has exactly the given code.
func (fs faultList) HasCode(code Code) bool {
	for _, f := range fs {
		if f.code == code {
			return true
		}
	}
	return false
}

// Len returns the number of faults in the collection.
func (fs faultList) Len() int {
	return len(fs)
}

// First returns the first fault in the collection.
func (fs faultList) First() *Fault {
	return fs[0]
}

// Last returns the last fault in the collection.
func (fs faultList) Last() *Fault {
	return fs[len(fs)-1]
}

// At returns the fault at position i.
func (fs faultList) At(i int) *Fault {
	return fs[i]
}

// List returns the underlying slice of faults.
func (fs faultList) List() []*Fault {
	return []*Fault(fs)
}

// prependPath returns a new slice with prefix prepended to each fault's path.
func prependPath(prefix string, faults []*Fault) []*Fault {
	if prefix == "" || len(faults) == 0 {
		return faults
	}
	result := make([]*Fault, len(faults))
	for i, f := range faults {
		result[i] = &Fault{
			path:    joinPath(prefix, f.path),
			code:    f.code,
			message: f.message,
		}
	}
	return result
}

// publicPath converts an internal relative path to JSON Path notation.
// An empty internal path (root-level fault) returns "$".
func publicPath(internal string) string {
	if internal == "" {
		return "$"
	}
	if internal[0] == '[' {
		return "$" + internal
	}
	return "$." + internal
}

func joinPath(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	if suffix == "" {
		return prefix
	}
	if len(suffix) > 0 && suffix[0] == '[' {
		return prefix + suffix
	}
	return prefix + "." + suffix
}
