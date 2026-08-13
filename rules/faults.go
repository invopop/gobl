package rules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fault represents a single rule assertion failure identified by a code and one or
// more paths. When multiple paths share the same code and message, they are merged
// into a single Fault. Fault is *not* designed to be instantiated directly, and
// will be created as part of the validation processes from defined rules.
//
//nolint:errname
type Fault struct {
	paths   []string
	code    Code
	message string
}

func newFault(path string, id Code, message string) *Fault {
	return &Fault{paths: []string{path}, code: id, message: message}
}

// Paths returns the JSON Path (RFC 6901) locations where this fault occurred.
func (f *Fault) Paths() []string {
	result := make([]string, len(f.paths))
	for i, p := range f.paths {
		result[i] = publicPath(p)
	}
	return result
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
	var pathPart string
	switch {
	case len(f.paths) == 1 && f.paths[0] != "":
		pathPart = "(" + publicPath(f.paths[0]) + ") "
	case len(f.paths) > 1:
		parts := make([]string, len(f.paths))
		for i, p := range f.paths {
			parts[i] = publicPath(p)
		}
		pathPart = "(" + strings.Join(parts, ", ") + ") "
	}
	return fmt.Sprintf("[%s] %s", f.code, pathPart+f.message)
}

// MarshalJSON encodes the fault as a JSON object with paths, code, and message fields.
func (f *Fault) MarshalJSON() ([]byte, error) {
	paths := make([]string, len(f.paths))
	for i, p := range f.paths {
		paths[i] = publicPath(p)
	}
	return json.Marshal(struct {
		Code    Code     `json:"code"`
		Paths   []string `json:"paths"`
		Message string   `json:"message"`
	}{f.code, paths, f.message})
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
// Faults with the same code and message are merged, combining their paths.
func newFaults(faults ...*Fault) Faults {
	if len(faults) == 0 {
		return nil
	}
	return faultList(mergeFaults(faults))
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
		for _, p := range f.paths {
			if publicPath(p) == path {
				return true
			}
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

// mergeFaults combines faults that share the same (code, message) pair,
// concatenating their paths into a single Fault.
func mergeFaults(faults []*Fault) []*Fault {
	type key struct {
		code    Code
		message string
	}
	seen := make(map[key]int) // index in result
	result := make([]*Fault, 0, len(faults))
	for _, f := range faults {
		k := key{f.code, f.message}
		if idx, ok := seen[k]; ok {
			result[idx].paths = append(result[idx].paths, f.paths...)
		} else {
			seen[k] = len(result)
			result = append(result, &Fault{
				paths:   append([]string(nil), f.paths...),
				code:    f.code,
				message: f.message,
			})
		}
	}
	return result
}

// prependPath returns a new slice with prefix prepended to each fault's paths.
func prependPath(prefix string, faults []*Fault) []*Fault {
	if prefix == "" || len(faults) == 0 {
		return faults
	}
	result := make([]*Fault, len(faults))
	for i, f := range faults {
		newPaths := make([]string, len(f.paths))
		for j, p := range f.paths {
			newPaths[j] = joinPath(prefix, p)
		}
		result[i] = &Fault{
			paths:   newPaths,
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
