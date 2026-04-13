package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStatusMinimal(t *testing.T) *bill.Status {
	t.Helper()
	return &bill.Status{
		Type:      bill.StatusTypeResponse,
		IssueDate: cal.MakeDate(2025, 3, 15),
		Series:    "S-1",
		Code:      "001",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B-98602642",
			},
		},
	}
}

func testStatusFull(t *testing.T) *bill.Status {
	t.Helper()
	st := testStatusMinimal(t)
	st.Customer = &org.Party{
		Name: "Test Customer",
		TaxID: &tax.Identity{
			Country: "ES",
			Code:    "54387763P",
		},
	}
	st.Issuer = &org.Party{
		Name: "Test Issuer",
	}
	st.Recipient = &org.Party{
		Name: "Test Recipient",
	}
	st.Ordering = &bill.Ordering{
		Code: "PO-123",
	}
	st.Lines = []*bill.StatusLine{
		{
			Key: bill.StatusEventAccepted,
			Doc: &org.DocumentRef{
				Series:    "F1",
				Code:      "0001",
				IssueDate: cal.NewDate(2025, 3, 10),
			},
			Description: "Invoice accepted for payment",
			Reasons: []*bill.Reason{
				{
					Key:         bill.ReasonKeyNone,
					Description: "All good",
				},
			},
			Actions: []*bill.Action{
				{
					Key:         bill.ActionKeyNone,
					Description: "No action needed",
				},
			},
		},
		{
			Key: bill.StatusEventRejected,
			Doc: &org.DocumentRef{
				Series: "F1",
				Code:   "0002",
			},
			Reasons: []*bill.Reason{
				{
					Key:         bill.ReasonKeyReferences,
					Description: "Missing PO reference",
					Conditions: []*bill.Condition{
						{
							Code:    "ERR-001",
							Paths:   []string{"ordering.code"},
							Message: "PO reference is required",
						},
					},
				},
			},
			Actions: []*bill.Action{
				{
					Key:         bill.ActionKeyReissue,
					Description: "Please reissue with correct PO reference",
				},
			},
		},
	}
	st.Notes = []*org.Note{
		{
			Key:  org.NoteKeyGeneral,
			Text: "Batch processing complete",
		},
	}
	return st
}

func TestStatusCalculate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		assert.Equal(t, "ES", st.Regime.Country.String())
		assert.Equal(t, "B98602642", st.Supplier.TaxID.Code.String(), "should normalize tax ID")
		assert.Equal(t, cbc.Code("S-1"), st.Series, "should normalize series code")
		assert.Equal(t, cbc.Code("001"), st.Code)
	})

	t.Run("regime from supplier", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		assert.False(t, st.Regime.IsEmpty())
		assert.Equal(t, "ES", st.Regime.Country.String())
	})

	t.Run("missing supplier", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Supplier = nil
		assert.NotPanics(t, func() {
			// Without supplier, regime can't be determined but should not panic
			_ = st.Calculate()
		})
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Supplier.TaxID = nil
		assert.NotPanics(t, func() {
			_ = st.Calculate()
		})
	})

	t.Run("without issue date", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.IssueDate = cal.Date{}
		require.NoError(t, st.Calculate())
		tn := cal.TodayIn(st.RegimeDef().TimeLocation())
		assert.Equal(t, tn, st.IssueDate)
		assert.Nil(t, st.IssueTime)
	})

	t.Run("with empty issue time", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.IssueDate = cal.Date{}
		st.IssueTime = new(cal.Time)
		require.NoError(t, st.Calculate())
		tn := cal.ThisSecondIn(st.RegimeDef().TimeLocation())
		assert.Equal(t, tn.Date().String(), st.IssueDate.String())
		assert.Equal(t, tn.Time().Hour, st.IssueTime.Hour)
		assert.Equal(t, tn.Time().Minute, st.IssueTime.Minute)
		assert.Equal(t, tn.Time().Second, st.IssueTime.Second)
	})

	t.Run("with preset issue time", func(t *testing.T) {
		st := testStatusMinimal(t)
		it := cal.MakeTime(10, 30, 0)
		st.IssueTime = &it
		require.NoError(t, st.Calculate())
		// Date and time should remain as-is
		assert.Equal(t, "2025-03-15", st.IssueDate.String())
		assert.Equal(t, 10, st.IssueTime.Hour)
		assert.Equal(t, 30, st.IssueTime.Minute)
	})

	t.Run("line indexing", func(t *testing.T) {
		st := testStatusFull(t)
		require.NoError(t, st.Calculate())
		assert.Equal(t, 1, st.Lines[0].Index)
		assert.Equal(t, 2, st.Lines[1].Index)
	})

	t.Run("line indexing with nil entries", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{Key: bill.StatusEventIssued},
			nil,
			{Key: bill.StatusEventAccepted},
		}
		require.NoError(t, st.Calculate())
		assert.Equal(t, 1, st.Lines[0].Index)
		assert.Equal(t, 3, st.Lines[2].Index)
	})

	t.Run("full status", func(t *testing.T) {
		st := testStatusFull(t)
		require.NoError(t, st.Calculate())
		assert.Equal(t, "ES", st.Regime.Country.String())
		assert.Equal(t, 1, st.Lines[0].Index)
		assert.Equal(t, 2, st.Lines[1].Index)
		assert.Equal(t, cbc.Code("PO-123"), st.Ordering.Code, "should normalize ordering code")
	})

	t.Run("normalize parties", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Customer = &org.Party{
			Name: "Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "5438-7763-P",
			},
		}
		require.NoError(t, st.Calculate())
		assert.Equal(t, "54387763P", st.Customer.TaxID.Code.String())
	})

	t.Run("with nil array entries", func(t *testing.T) {
		st := testStatusFull(t)
		st.Lines = append(st.Lines, nil)
		st.Notes = append(st.Notes, nil)
		st.Complements = append(st.Complements, nil)
		require.NoError(t, st.Calculate())
	})

	t.Run("no lines", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = nil
		require.NoError(t, st.Calculate())
	})

	t.Run("with addon", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Addons.SetAddons(tbai.V1)
		require.NoError(t, st.Calculate())
	})
}

func TestStatusValidate(t *testing.T) {
	t.Run("valid minimal", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("valid full", func(t *testing.T) {
		st := testStatusFull(t)
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing type", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		st.Type = ""
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status type is required")
	})

	t.Run("invalid type", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		st.Type = "invalid"
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status type is not valid")
	})

	t.Run("missing issue date", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.IssueDate = cal.MakeDate(2025, 3, 15)
		require.NoError(t, st.Calculate())
		st.IssueDate = cal.Date{}
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status issue date is required")
	})

	t.Run("missing code", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		st.Code = ""
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status code is required")
	})

	t.Run("missing supplier", func(t *testing.T) {
		st := testStatusMinimal(t)
		require.NoError(t, st.Calculate())
		st.Supplier = nil
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status supplier is required")
	})

	t.Run("all status types valid", func(t *testing.T) {
		for _, st := range []cbc.Key{bill.StatusTypeResponse, bill.StatusTypeUpdate, bill.StatusTypeSystem} {
			s := testStatusMinimal(t)
			s.Type = st
			require.NoError(t, s.Calculate())
			require.NoError(t, rules.Validate(s), "type %s should be valid", st)
		}
	})
}

func TestStatusLineValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{Key: bill.StatusEventAccepted},
		}
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{Key: ""},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status line key is required")
	})

	t.Run("invalid key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{Key: "invalid-event"},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "status line key is not valid")
	})

	t.Run("all event keys valid", func(t *testing.T) {
		events := []cbc.Key{
			bill.StatusEventIssued,
			bill.StatusEventAcknowledged,
			bill.StatusEventProcessing,
			bill.StatusEventQuerying,
			bill.StatusEventRejected,
			bill.StatusEventAccepted,
			bill.StatusEventPaid,
			bill.StatusEventError,
		}
		for _, ev := range events {
			st := testStatusMinimal(t)
			st.Lines = []*bill.StatusLine{
				{Key: ev},
			}
			require.NoError(t, st.Calculate())
			require.NoError(t, rules.Validate(st), "event key %s should be valid", ev)
		}
	})
}

func TestReasonValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Reasons: []*bill.Reason{
					{Key: bill.ReasonKeyReferences},
				},
			},
		}
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Reasons: []*bill.Reason{
					{Key: ""},
				},
			},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "reason key is required")
	})

	t.Run("invalid key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Reasons: []*bill.Reason{
					{Key: "bogus"},
				},
			},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "reason key is not valid")
	})

	t.Run("all reason keys valid", func(t *testing.T) {
		reasons := []cbc.Key{
			bill.ReasonKeyNone,
			bill.ReasonKeyReferences,
			bill.ReasonKeyLegal,
			bill.ReasonKeyUnknownReceiver,
			bill.ReasonKeyQuality,
			bill.ReasonKeyDelivery,
			bill.ReasonKeyPrices,
			bill.ReasonKeyQuantity,
			bill.ReasonKeyItems,
			bill.ReasonKeyPaymentTerms,
			bill.ReasonKeyNotRecognized,
			bill.ReasonKeyFinanceTerms,
			bill.ReasonKeyPartial,
			bill.ReasonKeyOther,
		}
		for _, rk := range reasons {
			st := testStatusMinimal(t)
			st.Lines = []*bill.StatusLine{
				{
					Key:     bill.StatusEventRejected,
					Reasons: []*bill.Reason{{Key: rk}},
				},
			}
			require.NoError(t, st.Calculate())
			require.NoError(t, rules.Validate(st), "reason key %s should be valid", rk)
		}
	})
}

func TestActionValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Actions: []*bill.Action{
					{Key: bill.ActionKeyReissue},
				},
			},
		}
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Actions: []*bill.Action{
					{Key: ""},
				},
			},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "action key is required")
	})

	t.Run("invalid key", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Actions: []*bill.Action{
					{Key: "bogus"},
				},
			},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "action key is not valid")
	})

	t.Run("all action keys valid", func(t *testing.T) {
		actions := []cbc.Key{
			bill.ActionKeyNone,
			bill.ActionKeyProvide,
			bill.ActionKeyReissue,
			bill.ActionKeyCreditFull,
			bill.ActionKeyCreditPartial,
			bill.ActionKeyCreditAmount,
			bill.ActionKeyOther,
		}
		for _, ak := range actions {
			st := testStatusMinimal(t)
			st.Lines = []*bill.StatusLine{
				{
					Key:     bill.StatusEventRejected,
					Actions: []*bill.Action{{Key: ak}},
				},
			}
			require.NoError(t, st.Calculate())
			require.NoError(t, rules.Validate(st), "action key %s should be valid", ak)
		}
	})
}

func TestConditionValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Reasons: []*bill.Reason{
					{
						Key: bill.ReasonKeyReferences,
						Conditions: []*bill.Condition{
							{
								Code:    "ERR-001",
								Paths:   []string{"ordering.code"},
								Message: "PO reference required",
							},
						},
					},
				},
			},
		}
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing code", func(t *testing.T) {
		st := testStatusMinimal(t)
		st.Lines = []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Reasons: []*bill.Reason{
					{
						Key: bill.ReasonKeyReferences,
						Conditions: []*bill.Condition{
							{Code: ""},
						},
					},
				},
			},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "condition code is required")
	})
}

func TestStatusLineNormalize(t *testing.T) {
	t.Run("nil status line", func(t *testing.T) {
		var sl *bill.StatusLine
		assert.NotPanics(t, func() {
			sl.Normalize(nil)
		})
	})
}

func TestReasonNormalize(t *testing.T) {
	t.Run("nil reason", func(t *testing.T) {
		var r *bill.Reason
		assert.NotPanics(t, func() {
			r.Normalize(nil)
		})
	})

	t.Run("normalizes conditions", func(t *testing.T) {
		r := &bill.Reason{
			Key: bill.ReasonKeyReferences,
			Conditions: []*bill.Condition{
				{Code: "ERR-001"},
			},
		}
		assert.NotPanics(t, func() {
			r.Normalize(nil)
		})
		assert.Equal(t, cbc.Code("ERR-001"), r.Conditions[0].Code)
	})
}

func TestActionNormalize(t *testing.T) {
	t.Run("nil action", func(t *testing.T) {
		var a *bill.Action
		assert.NotPanics(t, func() {
			a.Normalize(nil)
		})
	})

	t.Run("non-nil action", func(t *testing.T) {
		a := &bill.Action{Key: bill.ActionKeyReissue}
		assert.NotPanics(t, func() {
			a.Normalize(nil)
		})
	})
}

func TestConditionNormalize(t *testing.T) {
	t.Run("nil condition", func(t *testing.T) {
		var c *bill.Condition
		assert.NotPanics(t, func() {
			c.Normalize(nil)
		})
	})

	t.Run("normalizes code", func(t *testing.T) {
		c := &bill.Condition{Code: "  ERR--001  "}
		c.Normalize(nil)
		assert.Equal(t, cbc.Code("ERR-001"), c.Code)
	})
}

func TestStatusDefinitions(t *testing.T) {
	t.Run("status types count", func(t *testing.T) {
		assert.Len(t, bill.StatusTypes, 3)
	})

	t.Run("status events count", func(t *testing.T) {
		assert.Len(t, bill.StatusEvents, 8)
	})

	t.Run("reason keys count", func(t *testing.T) {
		assert.Len(t, bill.ReasonKeys, 14)
	})

	t.Run("action keys count", func(t *testing.T) {
		assert.Len(t, bill.ActionKeys, 7)
	})
}

func TestStatusJSONSchemaExtend(t *testing.T) {
	eg := `{
		"properties": {
			"$regime": {
				"$ref": "https://gobl.org/draft-0/tax/regime-code",
				"title": "Tax Regime"
			},
			"type": {
				"$ref": "https://gobl.org/draft-0/cbc/key",
				"title": "Type"
			},
			"series": {
				"$ref": "https://gobl.org/draft-0/cbc/code",
				"title": "Series"
			},
			"lines": {
				"type": "array",
				"title": "Lines"
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	st := bill.Status{}
	st.JSONSchemaExtend(js)

	t.Run("types", func(t *testing.T) {
		prop, ok := js.Properties.Get("type")
		require.True(t, ok)
		require.Len(t, prop.OneOf, len(bill.StatusTypes))
		assert.Equal(t, bill.StatusTypes[0].Key.String(), prop.OneOf[0].Const)
		assert.Equal(t, bill.StatusTypes[0].Name.String(), prop.OneOf[0].Title)
		assert.Equal(t, bill.StatusTypes[0].Desc.String(), prop.OneOf[0].Description)
	})

	t.Run("recommended", func(t *testing.T) {
		assert.Len(t, js.Extras[schema.Recommended], 3)
	})
}

func schemaWithKeyProp(t *testing.T) *jsonschema.Schema {
	t.Helper()
	eg := `{
		"properties": {
			"key": {
				"$ref": "https://gobl.org/draft-0/cbc/key",
				"title": "Key"
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))
	return js
}

func TestStatusLineJSONSchemaExtend(t *testing.T) {
	js := schemaWithKeyProp(t)

	sl := bill.StatusLine{}
	sl.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	require.True(t, ok)
	require.Len(t, prop.OneOf, len(bill.StatusEvents))
	assert.Equal(t, bill.StatusEvents[0].Key.String(), prop.OneOf[0].Const)
	assert.Equal(t, bill.StatusEvents[0].Name.String(), prop.OneOf[0].Title)
	assert.Equal(t, bill.StatusEvents[0].Desc.String(), prop.OneOf[0].Description)
}

func TestReasonJSONSchemaExtend(t *testing.T) {
	js := schemaWithKeyProp(t)

	r := bill.Reason{}
	r.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	require.True(t, ok)
	require.Len(t, prop.OneOf, len(bill.ReasonKeys))
	assert.Equal(t, bill.ReasonKeys[0].Key.String(), prop.OneOf[0].Const)
	assert.Equal(t, bill.ReasonKeys[0].Name.String(), prop.OneOf[0].Title)
	assert.Equal(t, bill.ReasonKeys[0].Desc.String(), prop.OneOf[0].Description)
}

func TestActionJSONSchemaExtend(t *testing.T) {
	js := schemaWithKeyProp(t)

	a := bill.Action{}
	a.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("key")
	require.True(t, ok)
	require.Len(t, prop.OneOf, len(bill.ActionKeys))
	assert.Equal(t, bill.ActionKeys[0].Key.String(), prop.OneOf[0].Const)
	assert.Equal(t, bill.ActionKeys[0].Name.String(), prop.OneOf[0].Title)
	assert.Equal(t, bill.ActionKeys[0].Desc.String(), prop.OneOf[0].Description)
}
