package legal

import (
	"strings"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// anchorPattern follows the JSON Schema plain-name syntax used by $anchor.
const anchorPattern = `[A-Za-z_][-A-Za-z0-9._]*`

var referenceFormat = is.AnyOf(
	is.RequestURI,
	is.Matches(`^#`+anchorPattern+`$`),
)

func contractRules() *rules.Set {
	return rules.For(new(Contract),
		rules.Field("title",
			rules.Assert("01", "contract title is required", is.Present),
		),
		rules.Field("chapters",
			rules.Assert("02", "contract requires at least one chapter", is.Present),
		),
		rules.Assert("03", "anchors must be unique throughout the contract",
			is.Func("unique anchors", contractAnchorsUnique),
		),
		rules.Assert("04", "local references must resolve to an anchor in the contract",
			is.Func("resolved local references", contractLocalRefsResolve),
		),
	)
}

func chapterRules() *rules.Set {
	return rules.For(new(Chapter),
		rules.Field("idx",
			rules.Assert("01", "chapter index is required and must be positive", is.Present, is.Min(int32(1))),
		),
		rules.Field("title",
			rules.Assert("02", "chapter title is required", is.Present),
		),
		rules.Field("$anchor",
			rules.AssertIfPresent("04", "chapter anchor must be a valid plain-name anchor",
				is.Matches(`^`+anchorPattern+`$`),
			),
		),
		rules.Field("$ref",
			rules.AssertIfPresent("05", "chapter reference must be an absolute URI, absolute path, or local anchor",
				referenceFormat,
			),
		),
		rules.Assert("03", "chapter requires sections or a reference",
			is.Expr(`len(Sections) > 0 || Ref != ""`),
		),
	)
}

func sectionRules() *rules.Set {
	return rules.For(new(Section),
		rules.Field("idx",
			rules.Assert("01", "section index is required and must be positive", is.Present, is.Min(int32(1))),
		),
		rules.Assert("02", "section requires content, sub-sections, or a reference",
			is.Expr(`Content != "" || len(Sections) > 0 || Ref != ""`),
		),
		rules.Field("$anchor",
			rules.AssertIfPresent("03", "section anchor must be a valid plain-name anchor",
				is.Matches(`^`+anchorPattern+`$`),
			),
		),
		rules.Field("$ref",
			rules.AssertIfPresent("04", "section reference must be an absolute URI, absolute path, or local anchor",
				referenceFormat,
			),
		),
	)
}

func contractAnchorsUnique(value any) bool {
	anchors := make(map[string]struct{})
	valid := true
	return walkContract(value, func(anchor, _ string) {
		if anchor == "" {
			return
		}
		if _, exists := anchors[anchor]; exists {
			valid = false
		}
		anchors[anchor] = struct{}{}
	}) && valid
}

func contractLocalRefsResolve(value any) bool {
	anchors := make(map[string]struct{})
	refs := make([]string, 0)
	if !walkContract(value, func(anchor, ref string) {
		if anchor != "" {
			anchors[anchor] = struct{}{}
		}
		if strings.HasPrefix(ref, "#") {
			refs = append(refs, strings.TrimPrefix(ref, "#"))
		}
	}) {
		return false
	}
	for _, ref := range refs {
		if _, exists := anchors[ref]; !exists {
			return false
		}
	}
	return true
}

func walkContract(value any, visit func(anchor, ref string)) bool {
	var c *Contract
	switch v := value.(type) {
	case *Contract:
		c = v
	case Contract:
		c = &v
	default:
		return false
	}
	if c == nil {
		return false
	}
	var walkSections func([]*Section)
	walkSections = func(sections []*Section) {
		for _, section := range sections {
			if section == nil {
				continue
			}
			visit(section.Anchor, section.Ref)
			walkSections(section.Sections)
		}
	}
	for _, chapter := range c.Chapters {
		if chapter == nil {
			continue
		}
		visit(chapter.Anchor, chapter.Ref)
		walkSections(chapter.Sections)
	}
	return true
}
