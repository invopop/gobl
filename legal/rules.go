package legal

import (
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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
		rules.Field("agreement",
			rules.Assert("01", "contract agreement UUID is required", is.Present),
		),
		rules.Field("kind",
			rules.AssertIfPresent("02", "contract kind is not valid",
				is.In(ContractKindAgreement, ContractKindDeed, ContractKindAmendment, ContractKindTerms),
			),
		),
		rules.Field("title",
			rules.Assert("03", "contract title is required", is.Present),
		),
		rules.Field("language",
			rules.Assert("06", "contract primary language is required", is.Present),
		),
		rules.Field("parties",
			rules.Assert("04", "contract requires at least one party", is.Present),
		),
		rules.Field("chapters",
			rules.Assert("05", "contract requires at least one chapter", is.Present),
		),
		rules.Assert("10", "anchors must be unique throughout the contract",
			is.Func("unique anchors", contractAnchorsUnique),
		),
		rules.Assert("11", "local references must resolve to an anchor in the contract",
			is.Func("resolved local references", contractLocalRefsResolve),
		),
		rules.Assert("12", "contract party and signatory keys must be unique",
			is.Func("unique party and signatory keys", contractPartyKeysUnique),
		),
		rules.Assert("13", "execution requirements must reference contract parties and signatories",
			is.Func("valid execution references", contractExecutionReferencesValid),
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

func partyRules() *rules.Set {
	return rules.For(new(Party),
		rules.Field("key",
			rules.Assert("01", "contract party key is required", is.Present),
		),
		rules.Field("role",
			rules.Assert("02", "contract party role is required", is.Present),
		),
		rules.Field("defined_as",
			rules.Assert("03", "contract party defined name is required", is.Present),
		),
		rules.Assert("04", "contract party requires an entity or description",
			is.Expr(`Entity != nil || Description != ""`),
		),
	)
}

func signatoryRules() *rules.Set {
	return rules.For(new(Signatory),
		rules.Field("key",
			rules.Assert("01", "signatory key is required", is.Present),
		),
		rules.Field("person",
			rules.Assert("02", "signatory person is required", is.Present),
		),
	)
}

func recitalRules() *rules.Set {
	return rules.For(new(Recital),
		rules.Field("$anchor",
			rules.AssertIfPresent("01", "recital anchor must be a valid plain-name anchor",
				is.Matches(`^`+anchorPattern+`$`),
			),
		),
		rules.Field("content",
			rules.Assert("02", "recital content is required", is.Present),
		),
	)
}

func definitionRules() *rules.Set {
	return rules.For(new(Definition),
		rules.Field("$anchor",
			rules.AssertIfPresent("01", "definition anchor must be a valid plain-name anchor",
				is.Matches(`^`+anchorPattern+`$`),
			),
		),
		rules.Field("term",
			rules.Assert("02", "defined term is required", is.Present),
		),
		rules.Field("meaning",
			rules.Assert("03", "definition meaning is required", is.Present),
		),
	)
}

func effectRules() *rules.Set {
	return rules.For(new(Effect),
		rules.Assert("10", "contract effect requires a start date, end date, or trigger",
			is.Expr(`Start != nil || End != nil || Trigger != ""`),
		),
	)
}

func governingLawRules() *rules.Set {
	return rules.For(new(GoverningLaw),
		rules.Assert("10", "governing law requires a country, jurisdiction, or legal instrument",
			is.Func("governing law identified", governingLawIdentified),
		),
	)
}

func disputeResolutionRules() *rules.Set {
	return rules.For(new(DisputeResolution),
		rules.Field("method",
			rules.Assert("01", "dispute resolution method is required", is.Present),
			rules.Assert("02", "dispute resolution method is not valid",
				is.In(DisputeMethodCourt, DisputeMethodArbitration, DisputeMethodMediation, DisputeMethodNegotiation),
			),
		),
	)
}

func executionRules() *rules.Set {
	return rules.For(new(Execution),
		rules.Field("requirements",
			rules.Assert("01", "execution requires at least one signature requirement", is.Present),
		),
	)
}

func signatureRequirementRules() *rules.Set {
	return rules.For(new(SignatureRequirement),
		rules.Field("party",
			rules.Assert("01", "signature requirement party is required", is.Present),
		),
		rules.Field("intent",
			rules.Assert("02", "signature requirement intent is required", is.Present),
			rules.Assert("03", "signature requirement intent is not valid",
				is.In(IntentExecute, IntentApprove, IntentAcknowledge, IntentWitness),
			),
		),
	)
}

func referenceRules() *rules.Set {
	return rules.For(new(Reference),
		rules.Field("relation",
			rules.Assert("01", "legal reference relation is required", is.Present),
		),
		rules.Field("url",
			rules.AssertIfPresent("02", "legal reference URL must be valid", is.URL),
		),
		rules.Assert("10", "legal reference requires a UUID, code, or URL",
			is.Func("legal reference identified", referenceIdentified),
		),
	)
}

func resourceRules() *rules.Set {
	return rules.For(new(Resource),
		rules.Field("key",
			rules.Assert("01", "legal resource key is required", is.Present),
		),
		rules.Field("title",
			rules.Assert("02", "legal resource title is required", is.Present),
		),
		rules.Field("url",
			rules.Assert("03", "legal resource URL is required", is.Present),
			rules.Assert("04", "legal resource URL must be valid", is.URL),
		),
		rules.When(is.Expr(`Incorporated`),
			rules.Field("digest",
				rules.Assert("20", "incorporated legal resource requires a digest", is.Present),
			),
		),
	)
}

func assentRules() *rules.Set {
	return rules.For(new(Assent),
		rules.Field("contract",
			rules.Assert("01", "assent contract reference is required", is.Present),
			rules.Field("digest",
				rules.Assert("02", "assent contract reference requires a digest", is.Present),
			),
		),
		rules.Field("party",
			rules.Assert("03", "assent party is required", is.Present),
		),
		rules.Field("intent",
			rules.Assert("04", "assent intent is required", is.Present),
			rules.Assert("05", "assent intent is not valid",
				is.In(IntentExecute, IntentApprove, IntentAcknowledge, IntentWitness),
			),
		),
	)
}

func analysisRules() *rules.Set {
	return rules.For(new(Analysis),
		rules.Field("contract",
			rules.Assert("01", "analysis contract reference is required", is.Present),
			rules.Field("digest",
				rules.Assert("02", "analysis contract reference requires a digest", is.Present),
			),
		),
		rules.Field("produced_at",
			rules.Assert("03", "analysis production time is required", cal.TimestampNotZero()),
		),
		rules.Field("method",
			rules.Assert("04", "analysis method is required", is.Present),
		),
		rules.Field("annotations",
			rules.Assert("05", "analysis requires at least one annotation", is.Present),
		),
	)
}

func annotationRules() *rules.Set {
	return rules.For(new(Annotation),
		rules.Field("anchor",
			rules.Assert("01", "annotation anchor is required", is.Present),
			rules.Assert("02", "annotation anchor must be a valid plain-name anchor",
				is.Matches(`^`+anchorPattern+`$`),
			),
		),
		rules.Field("kind",
			rules.Assert("03", "annotation kind is required", is.Present),
		),
		rules.Field("summary",
			rules.Assert("04", "annotation summary is required", is.Present),
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

func contractPartyKeysUnique(value any) bool {
	c := contractFrom(value)
	if c == nil {
		return false
	}
	parties := make(map[cbc.Key]struct{})
	for _, party := range c.Parties {
		if party == nil {
			continue
		}
		if _, exists := parties[party.Key]; exists {
			return false
		}
		parties[party.Key] = struct{}{}
		signatories := make(map[cbc.Key]struct{})
		for _, signatory := range party.Signatories {
			if signatory == nil {
				continue
			}
			if _, exists := signatories[signatory.Key]; exists {
				return false
			}
			signatories[signatory.Key] = struct{}{}
		}
	}
	return true
}

func contractExecutionReferencesValid(value any) bool {
	c := contractFrom(value)
	if c == nil {
		return false
	}
	if c.Execution == nil {
		return true
	}
	parties := make(map[cbc.Key]map[cbc.Key]struct{})
	for _, party := range c.Parties {
		if party == nil {
			continue
		}
		signatories := make(map[cbc.Key]struct{})
		for _, signatory := range party.Signatories {
			if signatory != nil {
				signatories[signatory.Key] = struct{}{}
			}
		}
		parties[party.Key] = signatories
	}
	for _, requirement := range c.Execution.Requirements {
		if requirement == nil {
			continue
		}
		signatories, exists := parties[requirement.Party]
		if !exists {
			return false
		}
		if requirement.Signatory != "" {
			if _, exists := signatories[requirement.Signatory]; !exists {
				return false
			}
		}
	}
	return true
}

func governingLawIdentified(value any) bool {
	var law *GoverningLaw
	switch v := value.(type) {
	case *GoverningLaw:
		law = v
	case GoverningLaw:
		law = &v
	}
	return law != nil && (law.Country != "" || law.Jurisdiction != "" || law.Instrument != "")
}

func referenceIdentified(value any) bool {
	var ref *Reference
	switch v := value.(type) {
	case *Reference:
		ref = v
	case Reference:
		ref = &v
	}
	return ref != nil && (ref.UUID != "" || ref.Code != "" || ref.URL != "")
}

func walkContract(value any, visit func(anchor, ref string)) bool {
	c := contractFrom(value)
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
	for _, recital := range c.Recitals {
		if recital != nil {
			visit(recital.Anchor, "")
		}
	}
	for _, definition := range c.Definitions {
		if definition != nil {
			visit(definition.Anchor, "")
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

func contractFrom(value any) *Contract {
	switch c := value.(type) {
	case *Contract:
		return c
	case Contract:
		return &c
	}
	return nil
}
