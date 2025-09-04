package org_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddInbox(t *testing.T) {
	t.Run("duplicate key", func(t *testing.T) {
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
	})
	t.Run("with append nil", func(t *testing.T) {
		st := struct {
			Inboxes []*org.Inbox
		}{
			Inboxes: []*org.Inbox{
				{
					Key:  "test-inbox",
					Code: "BAR",
				},
			},
		}
		st.Inboxes = org.AddInbox(st.Inboxes, nil)
		assert.Len(t, st.Inboxes, 1)
	})
	t.Run("with nil list", func(t *testing.T) {
		inboxes := org.AddInbox(nil, &org.Inbox{Code: "foo"})
		assert.Len(t, inboxes, 1)
		assert.Equal(t, "foo", inboxes[0].Code.String())
	})
	t.Run("with new inbox", func(t *testing.T) {
		inboxes := []*org.Inbox{
			{
				Key:  "other",
				Code: "BAR",
			},
		}
		inboxes = org.AddInbox(inboxes, &org.Inbox{Key: "test", Code: "FOO"})
		require.Len(t, inboxes, 2)
		assert.Equal(t, "BAR", inboxes[0].Code.String())
		assert.Equal(t, "FOO", inboxes[1].Code.String())
	})
}

func TestInboxNormalize(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		var id *org.Inbox
		assert.NotPanics(t, func() {
			id.Normalize()
		})
	})
	t.Run("with scheme", func(t *testing.T) {
		id := &org.Inbox{
			Scheme: " 0004 ",
			Code:   " BAR ",
		}
		id.Normalize()
		assert.Equal(t, "BAR", id.Code.String())
		assert.Equal(t, "0004", id.Scheme.String())
	})
	t.Run("with email in code", func(t *testing.T) {
		id := &org.Inbox{
			Code: "dev@invopop.com",
		}
		id.Normalize()
		assert.Empty(t, id.Code.String())
		assert.Empty(t, id.URL)
		assert.Equal(t, "dev@invopop.com", id.Email)
	})
	t.Run("with url in code", func(t *testing.T) {
		id := &org.Inbox{
			Code: "https://inbox.example.com",
		}
		id.Normalize()
		assert.Empty(t, id.Code.String())
		assert.Empty(t, id.Email)
		assert.Equal(t, "https://inbox.example.com", id.URL)
	})
	t.Run("with peppol participant code", func(t *testing.T) {
		id := &org.Inbox{
			Key:  org.InboxKeyPeppol,
			Code: "0004:1234567890",
		}
		id.Normalize()
		assert.Equal(t, "1234567890", id.Code.String())
		assert.Equal(t, "0004", id.Scheme.String())
		assert.Equal(t, org.InboxKeyPeppol, id.Key)
	})
	t.Run("with peppol participant code without key", func(t *testing.T) {
		id := &org.Inbox{
			Code: "0004:1234567890",
		}
		id.Normalize()
		assert.Equal(t, "0004:1234567890", id.Code.String())
	})
}

func TestInboxValidate(t *testing.T) {
	t.Run("with basics", func(t *testing.T) {
		id := &org.Inbox{
			Scheme: "0004",
			Code:   "BAR",
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
		assert.ErrorContains(t, err, "code: cannot be blank without url or email")
	})
	t.Run("with URL", func(t *testing.T) {
		id := &org.Inbox{
			URL: "https://inbox.example.com",
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("with invalid URL", func(t *testing.T) {
		id := &org.Inbox{
			URL: "https:/inbox",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "url: must be a valid URL")
	})
	t.Run("with code and URL", func(t *testing.T) {
		id := &org.Inbox{
			Code: "FOOO",
			URL:  "https://inbox.example.com",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "url: must be blank with code or email")
	})
	t.Run("with code and email", func(t *testing.T) {
		id := &org.Inbox{
			Code:  "FOOO",
			Email: "dev@invopop.com",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "email: must be blank with code or url")
	})
	t.Run("with email and url", func(t *testing.T) {
		id := &org.Inbox{
			Email: "dev@invopop.com",
			URL:   "https://inbox.example.com",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "email: must be blank with code or url; url: must be blank with code or email")
	})
	t.Run("with email", func(t *testing.T) {
		id := &org.Inbox{
			Email: "dev@invopop.com",
		}
		err := id.Validate()
		assert.NoError(t, err)
	})
	t.Run("with invalid email", func(t *testing.T) {
		id := &org.Inbox{
			Email: "dev@invopop",
		}
		err := id.Validate()
		assert.ErrorContains(t, err, "email: must be a valid email address")
	})

}
