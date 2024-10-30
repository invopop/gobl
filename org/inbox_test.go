package org_test

import (
	"testing"

	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestAddInbox(t *testing.T) {
	key := cbc.Key("test-inbox")
	st := struct {
		Inboxes []*org.Inbox
	}{
		Inboxes: []*org.Inbox{
			{
				Key:  key,
				Code: "BAR",
			},
		},
	}
	st.Inboxes = org.AddInbox(st.Inboxes, &org.Inbox{
		Key:  key,
		Code: "BARDOM",
	})
	assert.Len(t, st.Inboxes, 1)
	assert.Equal(t, "BARDOM", st.Inboxes[0].Code.String())
}

func TestInboxNormalize(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		var id *org.Inbox
		assert.NotPanics(t, func() {
			id.Normalize(nil)
		})
	})
	t.Run("missing extensions", func(t *testing.T) {
		id := &org.Inbox{
			Key:  cbc.Key("inbox"),
			Code: "BAR",
			Ext:  tax.Extensions{},
		}
		id.Normalize(nil)
		assert.Equal(t, "inbox", id.Key.String())
		assert.Nil(t, id.Ext)
	})
	t.Run("with extension", func(t *testing.T) {
		id := &org.Inbox{
			Code: "BAR",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0004",
			},
		}
		id.Normalize(nil)
		assert.Equal(t, "BAR", id.Code.String())
		assert.Equal(t, "0004", id.Ext[iso.ExtKeySchemeID].String())
	})
}

func TestInboxValidate(t *testing.T) {
	t.Run("with basics", func(t *testing.T) {
		id := &org.Inbox{
			Code: "BAR",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0004",
			},
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("with both key", func(t *testing.T) {
		id := &org.Inbox{
			Key:  "fiscal-code",
			Code: "1234567890",
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("missing code", func(t *testing.T) {
		id := &org.Inbox{
			Key: "fiscal-code",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "code: cannot be blank without url")
	})
	t.Run("with URL", func(t *testing.T) {
		id := &org.Inbox{
			URL: "https://inbox.example.com",
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("with code and URL", func(t *testing.T) {
		id := &org.Inbox{
			Code: "FOOO",
			URL:  "https://inbox.example.com",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "url: mutually exclusive with code")
	})
}
