package tax

import "reflect"

// Normalizer is used for functions that will normalize the provided object
// ensuring that all the data is aligned with expected values, and adding
// any additional data tha may be required.
//
// Normalizer cannot fail by design, they should always be designed to fail
// silently in case of issues and depend on the Validator to pick up on
// any issues.
type Normalizer func(doc any)

// Normalizers defines a list of normalizer methods with some helpers for
// execution.
type Normalizers []Normalizer

type regimeImpl interface {
	RegimeDef() *RegimeDef
}

type addonsImpl interface {
	normalizeAddons()
	AddonDefs() []*AddonDef
}

// ExtractNormalizers will extract the normalizers from the provided object
// that is using either the regime or addons.
func ExtractNormalizers(obj any) Normalizers {
	if obj == nil {
		return nil
	}
	normalizers := make(Normalizers, 0)
	if n, ok := obj.(regimeImpl); ok {
		if r := n.RegimeDef(); r != nil {
			normalizers = normalizers.Append(r.Normalizer)
		}
	}
	if n, ok := obj.(addonsImpl); ok {
		n.normalizeAddons()
		for _, a := range n.AddonDefs() {
			normalizers = normalizers.Append(a.Normalizer)
		}
	}
	return normalizers
}

type normalizeImpl interface {
	Normalize(Normalizers)
}

type normalizeSimpleImpl interface {
	Normalize()
}

// Each will run a simple loop over the normalizers on the provided object.
func (ns Normalizers) Each(doc any) {
	if doc == nil || ns == nil {
		return
	}
	for _, n := range ns {
		n(doc)
	}
}

// Append adds the normalizer, but only if it is not nil.
func (ns Normalizers) Append(n Normalizer) Normalizers {
	if n == nil {
		return ns
	}
	return append(ns, n)
}

// Normalize will either run the "Normalize" method on the provided object,
// or directly go through the list of normalizers on the object.
// This supports arrays and slices, and will automatically normalize each
// element in the list.
func Normalize(list Normalizers, doc any) {
	if doc == nil {
		return
	}
	if n, ok := doc.(normalizeImpl); ok {
		n.Normalize(list)
	} else if n, ok := doc.(normalizeSimpleImpl); ok {
		n.Normalize()
		list.Each(doc)
	} else {
		switch reflect.TypeOf(doc).Kind() {
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(doc)
			for i := 0; i < s.Len(); i++ {
				d := s.Index(i).Interface()
				Normalize(list, d)
			}
		default:
			list.Each(doc)
		}
	}
}
