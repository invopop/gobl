package dsig_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	Foo string `json:"foo"`
	Bar int64  `json:"bar"`
}

type structWithSig struct {
	Name string          `json:"name"`
	Sig  *dsig.Signature `json:"sig"`
}

func TestNewSignature(t *testing.T) {
	data := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4","d":"oJM3Ogl9uYUpSbc4oHV25DpFs_gOGP5nHJcLAtQxL6U"}`)
	wantID := "3500bbee-966c-4b7a-8fbc-c763ae2aec62"
	k := new(dsig.PrivateKey)
	if err := json.Unmarshal(data, k); err != nil {
		t.Errorf("failed to parse test key: %v", err.Error())
		return // abort
	}

	// jo, _ := json.Marshal(k.Public())
	// t.Logf("PUB KEY: %v", string(jo))

	p := new(payload)
	p.Foo = "foo" // nolint:goconst
	p.Bar = 1234
	s, err := dsig.NewSignature(k, p)
	if err != nil {
		t.Errorf("failed to create signature: %v", err.Error())
		return
	}

	t.Logf("signature: %v", s.String())
	if s.KeyID() != wantID {
		t.Errorf("execpted key IDs to be the same, got: %v", s.KeyID())
	}

	pub := k.Public()
	p2 := new(payload)
	if err := s.VerifyPayload(pub, p2); err != nil {
		t.Errorf("unexpected error verifying payload: %v", err.Error())
	}
	if p2.Foo != p.Foo || p2.Bar != p.Bar {
		t.Errorf("expected payloads to be the same, got: %+v", p2)
	}

	p3 := new(payload)
	if err := s.UnsafePayload(p3); err != nil {
		t.Errorf("did not expect unsafe payload to fail, got: %v", err.Error())
	}
	if p3.Foo != p.Foo || p3.Bar != p.Bar {
		t.Errorf("expected payloads to be the same, got: %+v", p3)
	}
}

func TestParseSignature(t *testing.T) {
	data := "eyJhbGciOiJFUzI1NiIsImtpZCI6IjM1MDBiYmVlLTk2NmMtNGI3YS04ZmJjLWM3NjNhZTJhZWM2MiJ9.eyJmb28iOiJmb28iLCJiYXIiOjEyMzR9.96eDPg1RJ4EXMXnCYYTbC3mIGU_DaKnULUdx6LxDeLh6kp-7G8V1CEr1Lwc-tqZ29iq6fwi0Pte-bnkBO0xh9w"
	pubData := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4"}`)

	pub := new(dsig.PublicKey)
	if err := json.Unmarshal(pubData, pub); err != nil {
		t.Errorf("expected to unmarshal pub key: %v", err.Error())
		return
	}

	s, err := dsig.ParseSignature(data)
	if err != nil {
		t.Errorf("expected signature to be parsed, got: %v", err.Error())
	}

	if s.KeyID() == "" || s.KeyID() != pub.ID() {
		t.Errorf("expected key ID to match, got: %v", s.KeyID())
	}

	p1 := new(payload)
	if err := s.UnsafePayload(p1); err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if p1.Foo != "foo" || p1.Bar != 1234 {
		t.Errorf("payload does not match: %+v", p1)
	}

	p2 := new(payload)
	if err := s.VerifyPayload(pub, p2); err != nil {
		t.Errorf("did not expect verify payload to fail: %v", err.Error())
	}
	if p2.Foo != "foo" || p2.Bar != 1234 {
		t.Errorf("failed to extract payload, got: %+v", p2)
	}
}

func TestSignaturesWithJKU(t *testing.T) {
	kdata := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4","d":"oJM3Ogl9uYUpSbc4oHV25DpFs_gOGP5nHJcLAtQxL6U"}`)
	k := new(dsig.PrivateKey)
	require.NoError(t, json.Unmarshal(kdata, k))

	p := new(payload)
	p.Foo = "foo"
	p.Bar = 1234
	jku := "https://ks.invopop.dev/NKFs8SbnEeykMgJCrBUACQ"
	s, err := dsig.NewSignature(k, p, dsig.WithJKU(jku))
	require.NoError(t, err)

	out := s.String()

	sig, err := dsig.ParseSignature(out)
	require.NoError(t, err)

	assert.Equal(t, jku, sig.JKU(), "should be included in signature output")
}

func TestJSONSignatures(t *testing.T) {
	pubData := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4"}`)
	pub := new(dsig.PublicKey)
	if err := json.Unmarshal(pubData, pub); err != nil {
		t.Errorf("expected to unmarshal pub key: %v", err.Error())
		return
	}

	// Test empty signature
	d1 := []byte(`{"name":"foo"}`)
	s1 := new(structWithSig)
	if err := json.Unmarshal(d1, s1); err != nil {
		t.Errorf("unexpected error parsing struct: %v", err.Error())
		return
	}
	if s1.Sig != nil {
		t.Errorf("expected nil signature")
	}

	// Test with signature and payload
	sigData := "eyJhbGciOiJFUzI1NiIsImtpZCI6IjM1MDBiYmVlLTk2NmMtNGI3YS04ZmJjLWM3NjNhZTJhZWM2MiJ9.eyJmb28iOiJmb28iLCJiYXIiOjEyMzR9.Cr-Kg0rHiPKmlTldQ5mgMIX4WRpKgbPA55TBf-avuPsfnxNpUgiwRV6lbbwcHIMsTl956oL62pJBnM5zSeITfw"
	d2 := []byte(`{"name":"foo","sig":"` + sigData + `"}`)
	s2 := new(structWithSig)
	if err := json.Unmarshal(d2, s2); err != nil {
		t.Errorf("unexpected unmarshal error: %v", err.Error())
	}
	if s2.Sig.KeyID() != pub.ID() {
		t.Errorf("expected key IDs to match, got: %v", s2.Sig.KeyID())
	}

	p2 := new(payload)
	if err := s2.Sig.VerifyPayload(pub, p2); err != nil {
		t.Errorf("unexpected verify error: %v", err.Error())
	}
	if p2.Foo != s2.Name {
		t.Errorf("expected names to be the same")
	}

	d3, err := json.Marshal(s2)
	if err != nil {
		t.Errorf("failed to marshal sig struct: %v", err.Error())
	}
	if !strings.Contains(string(d3), sigData) {
		t.Errorf("expected marshaled struct to include signature")
	}
}

func TestNewSignatureInvalidKey(t *testing.T) {
	// Invalid private key (no JWK) -> NewSignature short-circuits with
	// ErrKeyInvalid.
	_, err := dsig.NewSignature(&dsig.PrivateKey{}, struct{}{})
	require.Error(t, err)
	assert.ErrorIs(t, err, dsig.ErrKeyInvalid)
}

func TestNewSignatureUnmarshalablePayload(t *testing.T) {
	// A function value cannot be JSON-marshaled, exercising the
	// json.Marshal error branch in NewSignature.
	priv := dsig.NewES256Key()
	_, err := dsig.NewSignature(priv, func() {})
	require.Error(t, err)
}

func TestParseSignatureInvalid(t *testing.T) {
	_, err := dsig.ParseSignature("not.a.jws")
	require.Error(t, err)
}

func TestSignatureAccessorsEmpty(t *testing.T) {
	s := new(dsig.Signature)
	assert.Equal(t, "", s.KeyID())
	assert.Equal(t, "", s.JKU())
	assert.Equal(t, "", s.String())
	assert.Nil(t, s.JSONWebSignature())
}

func TestSignatureUnmarshalEmptyString(t *testing.T) {
	s := new(dsig.Signature)
	require.NoError(t, json.Unmarshal([]byte(`""`), s))
	assert.Nil(t, s.JSONWebSignature())
}

func TestSignatureUnmarshalInvalidJSON(t *testing.T) {
	s := new(dsig.Signature)
	err := json.Unmarshal([]byte(`not-a-string`), s)
	require.Error(t, err)
}

func TestSignatureUnmarshalBadCompact(t *testing.T) {
	s := new(dsig.Signature)
	err := json.Unmarshal([]byte(`"not.a.jws"`), s)
	require.Error(t, err)
}

func TestSignatureJSONSchema(t *testing.T) {
	js := dsig.Signature{}.JSONSchema()
	require.NotNil(t, js)
	assert.Equal(t, "string", js.Type)
}

func TestSignatureVerifyPayloadUnmarshalFail(t *testing.T) {
	// Sign a string payload, then try to verify into a struct that
	// can't accept it (target is an int).
	priv := dsig.NewES256Key()
	pub := new(dsig.PublicKey)
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, pub))
	sig, err := dsig.NewSignature(priv, "hello")
	require.NoError(t, err)
	var target int
	err = sig.VerifyPayload(pub, &target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dsig verify")
}

func TestSignatureUnsafePayloadUnmarshalFail(t *testing.T) {
	priv := dsig.NewES256Key()
	sig, err := dsig.NewSignature(priv, "hello")
	require.NoError(t, err)
	var target int
	err = sig.UnsafePayload(&target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsafe payload")
}

func TestSignatureVerifyWrongKey(t *testing.T) {
	// Verify against a key that didn't sign the payload — exercises
	// the ErrKeyMismatch return in Verify.
	signer := dsig.NewES256Key()
	other := dsig.NewES256Key()

	sig, err := dsig.NewSignature(signer, "hello")
	require.NoError(t, err)

	_, err = sig.Verify(other.Public())
	require.Error(t, err)
	assert.ErrorIs(t, err, dsig.ErrKeyMismatch)
}

func TestSignatureVerifyPayloadWrongKey(t *testing.T) {
	signer := dsig.NewES256Key()
	other := dsig.NewES256Key()

	sig, err := dsig.NewSignature(signer, "hello")
	require.NoError(t, err)

	var target string
	err = sig.VerifyPayload(other.Public(), &target)
	require.Error(t, err)
	assert.ErrorIs(t, err, dsig.ErrKeyMismatch)
}

func TestSignatureKeyID(t *testing.T) {
	priv := dsig.NewES256Key()
	sig, err := dsig.NewSignature(priv, "hello")
	require.NoError(t, err)
	assert.Equal(t, priv.ID(), sig.KeyID())
}
