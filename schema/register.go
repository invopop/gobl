package schema

import "reflect"

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

func (r *registry) add(id ID, obj interface{}) error {
	return r.addType(id, baseTypeOf(obj))
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

// Register
func Register(id ID, obj interface{}) error {
	return schemas.add(id, obj)
}

// RegisterIn will determine the anchor and add it to the base schema before
// adding to the global registry.
func RegisterIn(base ID, obj interface{}) error {
	if err := schemas.addWithAnchor(base, obj); err != nil {
		return err
	}
	return nil
}

// RegisterAllIn takes the base schema ID and adds all the provided objects as
// anchored entries in the base.
func RegisterAllIn(base ID, objs []interface{}) error {
	for _, obj := range objs {
		if err := RegisterIn(base, obj); err != nil {
			return err
		}
	}
	return nil
}

// Lookup finds the objects schema ID, if set
func Lookup(obj interface{}) ID {
	return schemas.lookup(obj)
}

// Type provides the type from a matching registered schema.
func Type(id ID) reflect.Type {
	return schemas.typeFor(id)
}
