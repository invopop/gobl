package schema

import (
	"path/filepath"
	"reflect"
)

// registry contains all the schemas that we can possibly know about from either
// inside or outside GOBL.
type registry struct {
	entries []*entry
}

type entry struct {
	id  ID
	typ reflect.Type
}

var schemas *registry

func newRegistry() *registry {
	return &registry{
		entries: make([]*entry, 0),
	}
}

func (r *registry) addType(id ID, typ reflect.Type) error {
	e := &entry{
		id:  id,
		typ: typ,
	}
	r.entries = append(r.entries, e)
	return nil
}

func (r *registry) add(base ID, obj interface{}) error {
	typ := baseTypeOf(obj)
	id := base.Add(typ.Name())
	return r.addType(id, typ)
}

func (r *registry) addWithAnchor(base ID, obj interface{}) error {
	typ := baseTypeOf(obj)
	id := base.Anchor(typ.Name())
	return r.addType(id, typ)
}

func (r *registry) lookup(obj interface{}) ID {
	typ := baseTypeOf(obj)
	for _, e := range r.entries {
		if typ == e.typ {
			return e.id
		}
	}
	return UnknownID
}

func (r *registry) typeFor(id ID) reflect.Type {
	for _, e := range r.entries {
		if id == e.id {
			return e.typ
		}
	}
	return nil
}

func (r *registry) ids() []ID {
	ids := make([]ID, len(r.entries))
	for i, e := range r.entries {
		ids[i] = e.id
	}
	return ids
}

// baseTypeOf removes the pointer and ensures we have a base type.
func baseTypeOf(obj interface{}) reflect.Type {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

// Register adds a new link between a schema ID and object to the global schema
// registry. This should be called for all GOBL models that will be included
// inside schema documents or included in an envelope document payload. The name
// of the object will be determined from the type of the object provided.
func Register(base ID, obj interface{}) {
	if err := schemas.add(base, obj); err != nil {
		panic(err)
	}
}

// RegisterIn will determine the anchor and add it to the base schema before
// adding to the global registry.
func RegisterIn(base ID, obj interface{}) {
	if err := schemas.addWithAnchor(base, obj); err != nil {
		panic(err)
	}
}

// RegisterAll takes an array of objects to register as additional schema using the
// ID#Add method. Se `RegistereAllIn` for registering schema using anchors.
func RegisterAll(base ID, objs []interface{}) {
	for _, obj := range objs {
		Register(base, obj)
	}
}

// RegisterAllIn takes the base schema ID and adds all the provided objects as
// anchored entries in the base.
func RegisterAllIn(base ID, objs []interface{}) {
	for _, obj := range objs {
		RegisterIn(base, obj)
	}
}

// Lookup finds the objects schema ID, if set
func Lookup(obj interface{}) ID {
	return schemas.lookup(obj)
}

// Type provides the type from a matching registered schema.
func Type(id ID) reflect.Type {
	return schemas.typeFor(id)
}

// Types provides a complete map of types to schema IDs that have been registered.
func Types() map[reflect.Type]ID {
	l := make(map[reflect.Type]ID)
	for _, e := range schemas.entries {
		l[e.typ] = e.id
	}
	return l
}

// List of known schema IDs. Mainly used for debugging.
func List() []ID {
	l := make([]ID, len(schemas.entries))
	for i, e := range schemas.entries {
		l[i] = e.id
	}
	return l
}

// Find searches the registered schemas for a match against term, and returns
// the best match along with a confidence value (1.0 = perfect match).
func Find(term string) (float32, ID) {
	return schemas.find(term)
}

func (r *registry) find(term string) (float32, ID) {
	for _, entry := range r.entries {
		if term == string(entry.id) {
			return 1, entry.id
		}
		if term == entry.typ.Name() {
			return 1, entry.id
		}
		if term == filepath.Base(entry.typ.PkgPath())+"."+entry.typ.Name() {
			return 1, entry.id
		}
	}
	return 0, ""
}
