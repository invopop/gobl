package legal_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/legal"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullContractArchitecture(t *testing.T) {
	contract := fullContract()
	norm.Normalize(contract)
	require.NoError(t, rules.Validate(contract))
	assert.Equal(t, legal.ContractKindAgreement, contract.Kind)
	assert.Equal(t, contract.UUID, contract.Agreement)
	assert.Equal(t, int32(1), contract.Chapters[0].Index)
}

func TestContractExecutionReferences(t *testing.T) {
	contract := fullContract()
	contract.Execution.Requirements[0].Signatory = "missing"
	norm.Normalize(contract)
	faults := rules.Validate(contract)
	require.Error(t, faults)
	assert.True(t, faults.HasCode("GOBL-LEGAL-CONTRACT-13"))
}

func TestContractAllowsOpenPartyRole(t *testing.T) {
	contract := fullContract()
	contract.Parties[1].Entity = nil
	contract.Parties[1].Description = "Any business that accepts the published terms."
	norm.Normalize(contract)
	require.NoError(t, rules.Validate(contract))
}

func TestIncorporatedResourceRequiresDigest(t *testing.T) {
	contract := fullContract()
	contract.Resources[0].Digest = nil
	norm.Normalize(contract)
	faults := rules.Validate(contract)
	require.Error(t, faults)
	assert.True(t, faults.HasCode("GOBL-LEGAL-RESOURCE-20"))
	assert.True(t, faults.HasPath("$.resources[0].digest"))
}

func TestAssentDocument(t *testing.T) {
	contract := fullContract()
	assent := &legal.Assent{
		Contract:  contractReference(contract),
		Party:     "provider",
		Signatory: "provider-director",
		Intent:    legal.IntentExecute,
		Capacity:  " Director ",
	}
	obj, err := schema.NewObject(assent)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())
	assert.Equal(t, "Director", assent.Capacity)
	assert.NoError(t, rules.Validate(assent))
	assert.Equal(t, schema.GOBL.Add(legal.ShortSchemaAssent), schema.Lookup(assent))
}

func TestAnalysisDocument(t *testing.T) {
	contract := fullContract()
	analysis := &legal.Analysis{
		Contract:   contractReference(contract),
		ProducedAt: cal.TimestampOf(time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)),
		Method:     " example-llm-v1 ",
		Annotations: []*legal.Annotation{
			{
				Anchor:  "services",
				Kind:    legal.AnnotationObligation,
				Subject: "provider",
				Object:  "customer",
				Summary: " Provider must perform the described services. ",
			},
		},
	}
	obj, err := schema.NewObject(analysis)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())
	assert.Equal(t, "example-llm-v1", analysis.Method.String())
	assert.Equal(t, "Provider must perform the described services.", analysis.Annotations[0].Summary)
	assert.NoError(t, rules.Validate(analysis))
	assert.Equal(t, schema.GOBL.Add(legal.ShortSchemaAnalysis), schema.Lookup(analysis))
}

func fullContract() *legal.Contract {
	contract := validContract()
	contract.Effect = &legal.Effect{Trigger: "when all required assents have been delivered"}
	contract.Recitals = []*legal.Recital{
		{Anchor: "purpose", Content: "The Customer wishes to obtain services from the Provider."},
	}
	contract.Definitions = []*legal.Definition{
		{Anchor: "services-definition", Term: "Services", Meaning: "The services described in Schedule 1."},
	}
	contract.Chapters[0].Sections[0].Anchor = "services"
	contract.Law = &legal.GoverningLaw{
		Country:      l10n.GB.ISO(),
		Jurisdiction: "England and Wales",
	}
	contract.Disputes = &legal.DisputeResolution{
		Method: legal.DisputeMethodCourt,
		Forum:  "courts of England and Wales",
	}
	contract.Parties[0].Signatories = []*legal.Signatory{
		{
			Key: "provider-director",
			Person: &org.Person{
				Name: &org.Name{Given: "Ada", Surname: "Example"},
			},
			Capacity: "Director",
		},
	}
	contract.Parties = append(contract.Parties, &legal.Party{
		Key:       "customer",
		Role:      "customer",
		DefinedAs: "Customer",
		Entity:    &org.Party{Name: "Example Customer Ltd"},
	})
	contract.Execution = &legal.Execution{
		Method:       "electronic-signature",
		Counterparts: true,
		Requirements: []*legal.SignatureRequirement{
			{Party: "provider", Signatory: "provider-director", Intent: legal.IntentExecute},
			{Party: "customer", Intent: legal.IntentExecute},
		},
	}
	contract.Resources = []*legal.Resource{
		{
			Key:          "schedule-1",
			Title:        "Schedule 1 – Services",
			URL:          "https://example.com/contracts/schedule-1.pdf",
			MIME:         "application/pdf",
			Digest:       testDigest(),
			Incorporated: true,
		},
	}
	return contract
}

func contractReference(contract *legal.Contract) *legal.Reference {
	return &legal.Reference{
		Relation: legal.ReferenceRelationAgreement,
		UUID:     contract.UUID,
		Title:    contract.Title,
		Digest:   testDigest(),
	}
}

func testDigest() *dsig.Digest {
	return &dsig.Digest{
		Algorithm: dsig.DigestSHA256,
		Value:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
}
