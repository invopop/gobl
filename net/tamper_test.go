package net

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/require"
)

// The tampering suite checks a receiver (VerifySender) against an
// attacker who can copy, edit, reorder and re-address any envelope and
// signature it sees, and sign with its own keys. Each case moves one
// fact away from what its signature binds it to (README.md §12.1), and
// expects a rejection, or the same result where the spec says the change
// does not matter.

const (
	twSupplier   = Address("supplier.example.com")
	twShadow     = Address("supplier-shadow.example.com") // the supplier's, unverified
	twMallory    = Address("mallory.example.com")
	twOther      = Address("other.example.com") // another receiver
	twAuthority  = Address("lookup.example.org")
	twVerifier   = Address("kyb.example.net")
	twSbAuth     = Address("lookup.sandbox.example.org")
	twSbVerifier = Address("email.sandbox.example.net")

	twSupplierName = "Supplier Ltd"
	twMalloryName  = "Mallory Ltd"
	twRejected     = ""
)

// twReceiver is the client running the checks.
var twReceiver = testRequester

var (
	twSupplierKey  = dsig.NewES256Key()
	twSupplierKey2 = dsig.NewES256Key() // published by the supplier, not bound
	twMalloryKey   = dsig.NewES256Key()
	twOtherKey     = dsig.NewES256Key()
	twAuthorityKey = dsig.NewES256Key()
	twVerifierKey  = dsig.NewES256Key()
	twSbAuthKey    = dsig.NewES256Key()
	twSbVerKey     = dsig.NewES256Key()
)

// twPublished lists the keys each address publishes.
var twPublished = map[Address][]*dsig.PrivateKey{
	twSupplier:   {twSupplierKey, twSupplierKey2},
	twShadow:     {twSupplierKey},
	twMallory:    {twMalloryKey},
	twOther:      {twOtherKey},
	twAuthority:  {twAuthorityKey},
	twVerifier:   {twVerifierKey},
	twSbAuth:     {twSbAuthKey},
	twSbVerifier: {twSbVerKey},
}

// exchange is what a receiver sees: a delivery, and the party envelope
// each address serves at /who.
type exchange struct {
	t        *testing.T
	who      map[Address]*gobl.Envelope
	delivery *gobl.Envelope
}

// newExchange returns a valid exchange: an invoice from the supplier,
// and verified, registered parties for the supplier and Mallory.
func newExchange(t *testing.T) *exchange {
	ex := &exchange{t: t, who: map[Address]*gobl.Envelope{}}
	ex.who[twSupplier] = ex.verifiedParty(twSupplierName, twSupplier, twSupplierKey)
	ex.who[twMallory] = ex.verifiedParty(twMalloryName, twMallory, twMalloryKey)
	ex.delivery = ex.invoice(twSupplier, twReceiver, twSupplierKey)
	return ex
}

func (ex *exchange) verifiedParty(name string, subject Address, key *dsig.PrivateKey) *gobl.Envelope {
	party := &org.Party{Name: name}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(ex.t, err)
	require.NoError(ex.t, Sign(env, key, subject, twReceiver))
	ex.add(env, twVerifierKey, verification(twVerifier, subject, key))
	ex.add(env, twAuthorityKey, registration(twAuthority, subject, key, twVerifier))
	return env
}

func (ex *exchange) invoice(from, to Address, keys ...*dsig.PrivateKey) *gobl.Envelope {
	msg := &note.Message{Content: "Invoice 1"}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(ex.t, err)
	for _, k := range keys {
		require.NoError(ex.t, Sign(env, k, from, to))
	}
	return env
}

// verification returns the claims of a verification countersignature.
func verification(iss, subject Address, bound *dsig.PrivateKey) head.SigningPayload {
	return head.SigningPayload{
		Iss:          iss.String(),
		Aud:          subject.String(),
		IssuedAt:     time.Now().Unix(),
		Confirmation: &head.Confirmation{JKT: bound.JKT()},
	}
}

// registration returns the claims of a registration countersignature.
func registration(iss, subject Address, bound *dsig.PrivateKey, verifier Address) head.SigningPayload {
	p := verification(iss, subject, bound)
	p.ExpiresAt = time.Now().Add(30 * 24 * time.Hour).Unix()
	p.Verifier = verifier.String()
	return p
}

// add signs env's document with key and exactly the claims in p.
func (ex *exchange) add(env *gobl.Envelope, key *dsig.PrivateKey, p head.SigningPayload) {
	p.UUID, p.Digest = env.Head.UUID, env.Head.Digest
	sig, err := dsig.NewSignature(key, &p)
	require.NoError(ex.t, err)
	env.Signatures = append(env.Signatures, sig)
}

// by returns env's signatures whose (unverified) issuer is iss.
func by(env *gobl.Envelope, iss Address) []*dsig.Signature {
	var out []*dsig.Signature
	for _, s := range env.Signatures {
		if p, err := head.SignedPayload(s); err == nil && p.Iss == iss.String() {
			out = append(out, s)
		}
	}
	return out
}

// drop removes env's signatures by iss.
func drop(env *gobl.Envelope, iss Address) {
	env.Signatures = slices.DeleteFunc(env.Signatures, func(s *dsig.Signature) bool {
		p, err := head.SignedPayload(s)
		return err == nil && p.Iss == iss.String()
	})
}

// replace swaps env's signatures by iss for sigs.
func replace(env *gobl.Envelope, iss Address, sigs ...*dsig.Signature) {
	drop(env, iss)
	env.Signatures = append(env.Signatures, sigs...)
}

// clone returns a deep copy of env.
func (ex *exchange) clone(env *gobl.Envelope) *gobl.Envelope {
	data, err := json.Marshal(env)
	require.NoError(ex.t, err)
	out := new(gobl.Envelope)
	require.NoError(ex.t, json.Unmarshal(data, out))
	return out
}

// editClaims changes a signature's payload without signing it again.
func (ex *exchange) editClaims(sig *dsig.Signature, edit func(map[string]any)) *dsig.Signature {
	parts := strings.Split(sig.String(), ".")
	require.Len(ex.t, parts, 3)
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(ex.t, err)
	claims := map[string]any{}
	require.NoError(ex.t, json.Unmarshal(raw, &claims))
	edit(claims)
	raw, err = json.Marshal(claims)
	require.NoError(ex.t, err)
	parts[1] = base64.RawURLEncoding.EncodeToString(raw)
	out, err := dsig.ParseSignature(strings.Join(parts, "."))
	require.NoError(ex.t, err)
	return out
}

// editAll applies editClaims to env's signatures by iss.
func (ex *exchange) editAll(env *gobl.Envelope, iss Address, edit func(map[string]any)) {
	sigs := by(env, iss)
	require.NotEmpty(ex.t, sigs)
	edited := make([]*dsig.Signature, len(sigs))
	for i, s := range sigs {
		edited[i] = ex.editClaims(s, edit)
	}
	replace(env, iss, edited...)
}

// rehash replaces env's document with one where was is replaced by is,
// with a digest recomputed to match it, keeping the signatures.
func (ex *exchange) rehash(env *gobl.Envelope, was, is string) *gobl.Envelope {
	out := tamper(ex.t, env, was, is)
	d, err := out.Digest()
	require.NoError(ex.t, err)
	out.Head.Digest = d
	return out
}

// verifySender runs the receiver's checks with the given trust options,
// returning the sender's party name, or twRejected and the error.
func (ex *exchange) verifySender(opts ...ClientOption) (string, error) {
	data := map[string][]byte{}
	for addr, keys := range twPublished {
		for _, k := range keys {
			data[addr.KeyURL(k.ID())] = jwkFromKey(ex.t, k)
		}
	}
	for addr, env := range ex.who {
		raw, err := json.Marshal(env)
		require.NoError(ex.t, err)
		data[addr.WhoURL()] = raw
	}
	opts = append(opts, WithIdentity(twReceiver, testRequesterKey), WithFetcher(&mapFetcher{data: data}))
	party, err := NewClient(opts...).VerifySender(context.Background(), ex.delivery)
	if err != nil {
		return twRejected, err
	}
	return party.Name, nil
}

type tamperCase struct {
	name string
	edit func(ex *exchange)
	// authority and direct are the party names expected by a client
	// trusting the authority, and one trusting the verifier directly;
	// twRejected expects a rejection.
	authority, direct string
}

func TestTampering(t *testing.T) {
	now := time.Now()
	sup := func(ex *exchange) *gobl.Envelope { return ex.who[twSupplier] }

	cases := []tamperCase{
		{name: "unchanged", edit: func(*exchange) {},
			authority: twSupplierName, direct: twSupplierName},

		// Document: every signature covers the document's UUID and
		// digest, and the digest must match the document.
		{name: "invoice body changed", edit: func(ex *exchange) {
			ex.delivery = tamper(t, ex.delivery, "Invoice 1", "Invoice 2")
		}},
		{name: "invoice body changed with its digest recomputed", edit: func(ex *exchange) {
			ex.delivery = ex.rehash(ex.delivery, "Invoice 1", "Invoice 2")
		}},
		{name: "invoice UUID changed", edit: func(ex *exchange) {
			ex.delivery.Head.UUID = uuid.V7()
		}},
		{name: "confirming signature moved from another invoice", edit: func(ex *exchange) {
			ex.delivery.Signatures = ex.invoice(twSupplier, twReceiver, twSupplierKey).Signatures
		}},
		{name: "party name changed", edit: func(ex *exchange) {
			ex.who[twSupplier] = tamper(t, sup(ex), twSupplierName, "Someone Else")
		}},
		{name: "party name changed with its digest recomputed", edit: func(ex *exchange) {
			ex.who[twSupplier] = ex.rehash(sup(ex), twSupplierName, "Someone Else")
		}},
		{name: "party UUID changed", edit: func(ex *exchange) {
			sup(ex).Head.UUID = uuid.V7()
		}},
		{name: "verification moved from mallory's party", edit: func(ex *exchange) {
			replace(sup(ex), twVerifier, by(ex.who[twMallory], twVerifier)...)
		}},
		{name: "registration moved from mallory's party", edit: func(ex *exchange) {
			replace(sup(ex), twAuthority, by(ex.who[twMallory], twAuthority)...)
		}, direct: twSupplierName},

		// Address: countersignatures name their subject as audience,
		// and the confirming signature names sender and receiver.
		{name: "verified party served at another address with the same key", edit: func(ex *exchange) {
			env := ex.clone(sup(ex))
			drop(env, twSupplier)
			require.NoError(t, Sign(env, twSupplierKey, twShadow, twReceiver))
			ex.who[twShadow] = env
			ex.delivery = ex.invoice(twShadow, twReceiver, twSupplierKey)
		}},
		{name: "verified party served by mallory as its own", edit: func(ex *exchange) {
			env := ex.clone(sup(ex))
			drop(env, twSupplier)
			require.NoError(t, Sign(env, twMalloryKey, twMallory, twReceiver))
			ex.who[twMallory] = env
			ex.delivery = ex.invoice(twMallory, twReceiver, twMalloryKey)
		}},
		{name: "verification for another address of the subject", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
			ex.add(sup(ex), twVerifierKey, verification(twVerifier, twShadow, twSupplierKey))
		}},
		{name: "invoice sent to another receiver, re-addressed", edit: func(ex *exchange) {
			ex.delivery = ex.invoice(twSupplier, twOther, twSupplierKey)
			ex.delivery.Head.To = twReceiver.URI()
		}},
		{name: "invoice from changed to mallory", edit: func(ex *exchange) {
			ex.delivery.Head.From = twMallory.URI()
		}},
		{name: "invoice confirmed without an audience", edit: func(ex *exchange) {
			ex.delivery.Signatures = nil
			ex.add(ex.delivery, twSupplierKey, head.SigningPayload{Iss: twSupplier.String(), IssuedAt: now.Unix()})
		}},
		{name: "who response made for another requester", edit: func(ex *exchange) {
			drop(sup(ex), twSupplier)
			require.NoError(t, Sign(sup(ex), twSupplierKey, twSupplier, twOther))
		}},
		{name: "who response from mallory served at the supplier", edit: func(ex *exchange) {
			ex.who[twSupplier] = ex.who[twMallory]
		}},
		// Passing an envelope on as one's own is allowed (README.md
		// §5): the sender is then mallory, not the supplier.
		{name: "invoice forwarded by mallory as its own", edit: func(ex *exchange) {
			require.NoError(t, Sign(ex.delivery, twMalloryKey, twMallory, twReceiver))
		}, authority: twMalloryName, direct: twMalloryName},
		{name: "registration audience in another form", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			ex.add(sup(ex), twAuthorityKey,
				registration(twAuthority, "SUPPLIER.Example.COM.", twSupplierKey, twVerifier))
		}, authority: twSupplierName, direct: twSupplierName},

		// Key: countersignatures bind the subject's key, which must
		// confirm both the party envelope and the delivery.
		{name: "invoice confirmed with an unbound key of the supplier", edit: func(ex *exchange) {
			ex.delivery = ex.invoice(twSupplier, twReceiver, twSupplierKey2)
		}},
		{name: "who confirmed with an unbound key of the supplier", edit: func(ex *exchange) {
			drop(sup(ex), twSupplier)
			require.NoError(t, Sign(sup(ex), twSupplierKey2, twSupplier, twReceiver))
		}},
		{name: "key rotation signs with both keys", edit: func(ex *exchange) {
			require.NoError(t, Sign(sup(ex), twSupplierKey2, twSupplier, twReceiver))
			ex.delivery = ex.invoice(twSupplier, twReceiver, twSupplierKey, twSupplierKey2)
		}, authority: twSupplierName, direct: twSupplierName},
		{name: "registration binds a key that confirmed nothing", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			ex.add(sup(ex), twAuthorityKey, registration(twAuthority, twSupplier, twSupplierKey2, twVerifier))
		}, direct: twSupplierName},
		{name: "verification binds a key that confirmed nothing", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
			ex.add(sup(ex), twVerifierKey, verification(twVerifier, twSupplier, twSupplierKey2))
		}, authority: twSupplierName},
		{name: "registration without a bound key", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			p := registration(twAuthority, twSupplier, twSupplierKey, twVerifier)
			p.Confirmation = nil
			ex.add(sup(ex), twAuthorityKey, p)
		}, direct: twSupplierName},
		{name: "verification without a bound key", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
			p := verification(twVerifier, twSupplier, twSupplierKey)
			p.Confirmation = nil
			ex.add(sup(ex), twVerifierKey, p)
		}},

		// Trust: only countersignatures by trusted participants count.
		{name: "registration by an untrusted authority", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			ex.add(sup(ex), twMalloryKey, registration(twMallory, twSupplier, twSupplierKey, twVerifier))
		}, direct: twSupplierName},
		{name: "sandbox verification and registration", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			drop(sup(ex), twVerifier)
			ex.add(sup(ex), twSbVerKey, verification(twSbVerifier, twSupplier, twSupplierKey))
			ex.add(sup(ex), twSbAuthKey, registration(twSbAuth, twSupplier, twSupplierKey, twSbVerifier))
		}},
		{name: "registration naming the authority as verifier", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			ex.add(sup(ex), twAuthorityKey, registration(twAuthority, twSupplier, twSupplierKey, twAuthority))
		}, direct: twSupplierName},
		{name: "registration naming a verifier with no countersignature", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			ex.add(sup(ex), twAuthorityKey, registration(twAuthority, twSupplier, twSupplierKey, "kyb2.example.net"))
		}, direct: twSupplierName},
		{name: "verification removed", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
		}},
		{name: "registration removed", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
		}, direct: twSupplierName},

		// Time: signing times are not in the future, expiries have not
		// passed, and registrations last at most 90 days.
		{name: "registration expired", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			p := registration(twAuthority, twSupplier, twSupplierKey, twVerifier)
			p.IssuedAt, p.ExpiresAt = now.Add(-2*time.Hour).Unix(), now.Add(-time.Hour).Unix()
			ex.add(sup(ex), twAuthorityKey, p)
		}, direct: twSupplierName},
		{name: "registration lasting more than 90 days", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			p := registration(twAuthority, twSupplier, twSupplierKey, twVerifier)
			p.ExpiresAt = now.Add(91 * 24 * time.Hour).Unix()
			ex.add(sup(ex), twAuthorityKey, p)
		}, direct: twSupplierName},
		{name: "registration without an expiry", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			p := registration(twAuthority, twSupplier, twSupplierKey, twVerifier)
			p.ExpiresAt = 0
			ex.add(sup(ex), twAuthorityKey, p)
		}, direct: twSupplierName},
		{name: "registration signed in the future", edit: func(ex *exchange) {
			drop(sup(ex), twAuthority)
			p := registration(twAuthority, twSupplier, twSupplierKey, twVerifier)
			p.IssuedAt = now.Add(time.Hour).Unix()
			ex.add(sup(ex), twAuthorityKey, p)
		}, direct: twSupplierName},
		{name: "verification expired", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
			p := verification(twVerifier, twSupplier, twSupplierKey)
			p.IssuedAt, p.ExpiresAt = now.Add(-2*time.Hour).Unix(), now.Add(-time.Hour).Unix()
			ex.add(sup(ex), twVerifierKey, p)
		}},
		{name: "verification signed in the future", edit: func(ex *exchange) {
			drop(sup(ex), twVerifier)
			p := verification(twVerifier, twSupplier, twSupplierKey)
			p.IssuedAt = now.Add(time.Hour).Unix()
			ex.add(sup(ex), twVerifierKey, p)
		}},
		{name: "invoice confirmed without a signing time", edit: func(ex *exchange) {
			ex.delivery.Signatures = nil
			ex.add(ex.delivery, twSupplierKey, head.SigningPayload{
				Iss: twSupplier.String(), Aud: twReceiver.String()})
		}},
		{name: "invoice confirmed in the future", edit: func(ex *exchange) {
			ex.delivery.Signatures = nil
			ex.add(ex.delivery, twSupplierKey, head.SigningPayload{
				Iss: twSupplier.String(), Aud: twReceiver.String(), IssuedAt: now.Add(time.Hour).Unix()})
		}},

		// Signatures: order, duplicates and unsigned header fields
		// change nothing; too many signatures are refused.
		{name: "party signatures reordered", edit: func(ex *exchange) {
			slices.Reverse(sup(ex).Signatures)
		}, authority: twSupplierName, direct: twSupplierName},
		{name: "party signatures duplicated", edit: func(ex *exchange) {
			sup(ex).Signatures = append(sup(ex).Signatures, sup(ex).Signatures...)
		}, authority: twSupplierName, direct: twSupplierName},
		{name: "party with too many signatures", edit: func(ex *exchange) {
			for len(sup(ex).Signatures) <= maxEnvelopeSignatures {
				sup(ex).Signatures = append(sup(ex).Signatures, sup(ex).Signatures[0])
			}
		}},
		{name: "invoice header stamp and notes added", edit: func(ex *exchange) {
			ex.delivery.Head.Notes = "unsigned"
			ex.delivery.Head.Stamps = append(ex.delivery.Head.Stamps, &head.Stamp{Provider: "x", Value: "y"})
		}, authority: twSupplierName, direct: twSupplierName},
	}

	// Editing any claim of a signature without signing again.
	edits := map[string]any{
		"iss":      twMallory.String(),
		"aud":      twMallory.String(),
		"iat":      now.Add(-time.Minute).Unix(),
		"exp":      now.Add(10 * 24 * time.Hour).Unix(),
		"cnf":      map[string]any{"jkt": twMalloryKey.JKT()},
		"verifier": twMallory.String(),
		"uuid":     uuid.V7().String(),
	}
	roles := []struct {
		name   string
		env    func(ex *exchange) *gobl.Envelope
		iss    Address
		direct string
	}{
		{"invoice confirmation", func(ex *exchange) *gobl.Envelope { return ex.delivery }, twSupplier, twRejected},
		{"who confirmation", sup, twSupplier, twRejected},
		{"verification", sup, twVerifier, twRejected},
		{"registration", sup, twAuthority, twSupplierName},
	}
	for _, r := range roles {
		for claim, value := range edits {
			cases = append(cases, tamperCase{
				name: r.name + " claim " + claim + " edited",
				edit: func(ex *exchange) {
					ex.editAll(r.env(ex), r.iss, func(c map[string]any) { c[claim] = value })
				},
				direct: r.direct,
			})
		}
	}

	clients := []struct {
		name string
		opts []ClientOption
		want func(tamperCase) string
	}{
		{"authority", []ClientOption{WithAuthorities(twAuthority)}, func(c tamperCase) string { return c.authority }},
		{"direct", []ClientOption{WithAuthorities(), WithVerifiers(twVerifier)}, func(c tamperCase) string { return c.direct }},
	}
	for _, tc := range cases {
		for _, cl := range clients {
			t.Run(tc.name+"/"+cl.name, func(t *testing.T) {
				ex := newExchange(t)
				tc.edit(ex)
				got, err := ex.verifySender(cl.opts...)
				want := cl.want(tc)
				if want == twRejected {
					require.Error(t, err, "accepted as %q", got)
					// Any other error means the case failed to set up,
					// not that the receiver rejected the change.
					require.True(t, errors.Is(err, ErrVerifyFailed) || errors.Is(err, ErrNotEndorsed) ||
						errors.Is(err, ErrSignatureExpired), "rejected for an unrelated reason: %v", err)
					return
				}
				require.NoError(t, err)
				require.Equal(t, want, got)
			})
		}
	}
}
