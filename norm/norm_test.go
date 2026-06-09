package norm

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/stretchr/testify/assert"
)

// --- mutation through nesting, slices and pointers ---

type mItem struct {
	Code string
}

type mSub struct {
	Val string
}

type mRoot struct {
	Name  string
	Items []*mItem
	Vals  []mSub
	Sub   *mSub
}

func init() {
	Register(
		For(func(r *mRoot) { r.Name = strings.ToUpper(r.Name) }),
		For(func(i *mItem) { i.Code = strings.ToUpper(i.Code) }),
		For(func(s *mSub) { s.Val = strings.ToUpper(s.Val) }),
	)
}

func TestNormalizeMutation(t *testing.T) {
	r := &mRoot{
		Name:  "root",
		Items: []*mItem{{Code: "a"}, {Code: "b"}},
		Vals:  []mSub{{Val: "x"}, {Val: "y"}},
		Sub:   &mSub{Val: "z"},
	}
	Normalize(r)

	assert.Equal(t, "ROOT", r.Name)
	assert.Equal(t, "A", r.Items[0].Code)
	assert.Equal(t, "B", r.Items[1].Code)
	assert.Equal(t, "X", r.Vals[0].Val, "value-slice elements are addressable and mutated")
	assert.Equal(t, "Y", r.Vals[1].Val)
	assert.Equal(t, "Z", r.Sub.Val, "pointer field is dereferenced and mutated")
}

func TestNormalizeNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { Normalize(nil) })
	assert.NotPanics(t, func() { Normalize((*mRoot)(nil)) })
	assert.NotPanics(t, func() { Normalize(mRoot{}) }, "non-pointer root is a no-op")
}

// --- post-order traversal: children before parent ---

type oChild struct{ order *[]string }
type oParent struct {
	Child oChild
	order *[]string
}

func init() {
	Register(
		For(func(p *oParent) { *p.order = append(*p.order, "parent") }),
		For(func(c *oChild) { *c.order = append(*c.order, "child") }),
	)
}

func TestNormalizePostOrder(t *testing.T) {
	var order []string
	p := &oParent{Child: oChild{order: &order}, order: &order}
	Normalize(p)
	assert.Equal(t, []string{"child", "parent"}, order, "child normalized before parent")
}

// --- unguarded normalizers run before guarded ones ---

type ordRoot struct {
	Keys  []string
	trace *[]string
}

func (r ordRoot) RulesContext() rules.WithContext {
	return func(rc *rules.Context) {
		for _, k := range r.Keys {
			rc.Set(rules.ContextKey(k), k)
		}
	}
}

func ctxHas(key string) rules.Test {
	return is.InContext(is.Func("has "+key, func(v any) bool {
		s, _ := v.(string)
		return s == key
	}))
}

func init() {
	RegisterWithGuard(ctxHas("guard"),
		For(func(r *ordRoot) { *r.trace = append(*r.trace, "guarded") }),
	)
	Register(
		For(func(r *ordRoot) { *r.trace = append(*r.trace, "intrinsic") }),
	)
}

func TestNormalizeUnguardedFirst(t *testing.T) {
	var trace []string
	r := &ordRoot{Keys: []string{"guard"}, trace: &trace}
	Normalize(r)
	assert.Equal(t, []string{"intrinsic", "guarded"}, trace,
		"intrinsic (unguarded) normalizer runs before guarded one regardless of registration order")
}

// --- guard gating by context ---

type gRoot struct {
	Keys []string
	hits *int
}

func (r gRoot) RulesContext() rules.WithContext {
	return func(rc *rules.Context) {
		for _, k := range r.Keys {
			rc.Set(rules.ContextKey(k), k)
		}
	}
}

func init() {
	RegisterWithGuard(ctxHas("on"),
		For(func(r *gRoot) { *r.hits++ }),
	)
}

func TestNormalizeGuardGating(t *testing.T) {
	on := 0
	Normalize(&gRoot{Keys: []string{"on"}, hits: &on})
	assert.Equal(t, 1, on, "guarded normalizer runs when its key is in context")

	off := 0
	Normalize(&gRoot{Keys: []string{"other"}, hits: &off})
	assert.Equal(t, 0, off, "guarded normalizer is skipped when its key is absent")
}

// --- meta-addon: a key appended mid-run activates a guarded normalizer ---

type metaChild struct{ count *int }
type metaRoot struct {
	Keys  []string
	Child metaChild
}

func (r metaRoot) RulesContext() rules.WithContext {
	return func(rc *rules.Context) {
		for _, k := range r.Keys {
			rc.Set(rules.ContextKey(k), k)
		}
	}
}

func init() {
	// Unguarded: on first sight, append key "B" (idempotent).
	Register(
		For(func(r *metaRoot) {
			for _, k := range r.Keys {
				if k == "B" {
					return
				}
			}
			r.Keys = append(r.Keys, "B")
		}),
	)
	// Guarded by "B": only runs once "B" is present in context.
	RegisterWithGuard(ctxHas("B"),
		For(func(c *metaChild) { *c.count++ }),
	)
}

func TestNormalizeMetaAddonPass(t *testing.T) {
	count := 0
	r := &metaRoot{Keys: []string{"A"}, Child: metaChild{count: &count}}
	Normalize(r)
	assert.Contains(t, r.Keys, "B", "meta normalizer appended the addon key")
	assert.Equal(t, 1, count,
		"normalizer guarded by the appended key ran on the follow-up pass")
}

// --- pass cap: keys that never stabilise stop at maxPasses ---

type runawayRoot struct {
	Keys   []string
	passes *int
}

func (r runawayRoot) RulesContext() rules.WithContext {
	return func(rc *rules.Context) {
		for _, k := range r.Keys {
			rc.Set(rules.ContextKey(k), k)
		}
	}
}

func init() {
	Register(
		For(func(r *runawayRoot) {
			*r.passes++
			// Always add a brand-new key, so the key set never stabilises.
			r.Keys = append(r.Keys, strings.Repeat("k", len(r.Keys)+1))
		}),
	)
}

func TestNormalizePassCap(t *testing.T) {
	passes := 0
	Normalize(&runawayRoot{passes: &passes})
	assert.Equal(t, maxPasses, passes, "runaway key growth is bounded by maxPasses")
}

// --- walk reaches maps, interfaces, arrays and pointers; When guards by value ---

type wLeaf struct{ Name string }

type wContainer struct {
	ByKey map[string]*wLeaf // map with pointer values
	Iface any               // interface holding a pointer
	Arr   [2]wLeaf          // array of struct values
	Ptr   *wLeaf            // pointer (may be nil)

	guarded bool // set by a When-guarded normalizer
}

func init() {
	Register(
		For(func(l *wLeaf) { l.Name = strings.ToUpper(l.Name) }),
	)
	// When with a per-value (non-contextual) guard, exercising the plain
	// Test.Check path in runTest.
	Register(
		When(is.Func("is container", func(v any) bool {
			_, ok := v.(*wContainer)
			return ok
		}),
			For(func(c *wContainer) { c.guarded = true }),
		),
	)
}

func TestNormalizeWalksContainers(t *testing.T) {
	c := &wContainer{
		ByKey: map[string]*wLeaf{"a": {Name: "x"}},
		Iface: &wLeaf{Name: "y"},
		Arr:   [2]wLeaf{{Name: "p"}, {Name: "q"}},
		Ptr:   &wLeaf{Name: "z"},
	}
	Normalize(c)
	assert.Equal(t, "X", c.ByKey["a"].Name, "pointer values in maps are normalized")
	assert.Equal(t, "Y", c.Iface.(*wLeaf).Name, "pointers behind interfaces are normalized")
	assert.Equal(t, "P", c.Arr[0].Name, "array elements are normalized")
	assert.Equal(t, "Q", c.Arr[1].Name)
	assert.Equal(t, "Z", c.Ptr.Name, "pointer fields are normalized")
	assert.True(t, c.guarded, "When-guarded normalizer ran")
}

func TestNormalizeWalksNilMembers(t *testing.T) {
	c := &wContainer{} // nil map, nil interface, nil pointer, zero array
	assert.NotPanics(t, func() { Normalize(c) })
	assert.True(t, c.guarded, "node normalizer still runs with nil members")
}
