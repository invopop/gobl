package dsig_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/invopop/gobl/dsig"
)

func TestNewKeyPair(t *testing.T) {
	k := dsig.NewES256Key()

	if k.ID() == "" {
		t.Errorf("expected an ID")
	}
	if err := k.Validate(); err != nil {
		t.Errorf("expected a valid new key")
	}
	if k.Thumbprint() == "" {
		t.Errorf("expected a thumbprint")
	}
	data, err := json.Marshal(k)
	if err != nil {
		t.Errorf("failed to generate JSON data: %v", err.Error())
	}
	t.Logf("json output of private key: %v", string(data))

	pub := k.Public()
	if err := k.Validate(); err != nil {
		t.Errorf("expected a valid new key")
	}
	if k.ID() != pub.ID() {
		t.Errorf("expected matching keys: %v", pub.ID())
	}
	if k.Thumbprint() != pub.Thumbprint() {
		t.Errorf("expected matching thumbprints, got: %v != %v", k.Thumbprint(), pub.Thumbprint())
	}
	pubdata, err := json.Marshal(pub)
	if err != nil {
		t.Errorf("failed to generate JSON data: %v", err.Error())
	}
	t.Logf("json output of public key: %v", string(pubdata))
}

func TestKeyParsing(t *testing.T) {
	data := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4","d":"oJM3Ogl9uYUpSbc4oHV25DpFs_gOGP5nHJcLAtQxL6U"}`)
	kid := "3500bbee-966c-4b7a-8fbc-c763ae2aec62"
	thumb := "eb9793f9b542e25e13959591a41391329f2eb05a0d8f9b12eb071b4d90d91c51"

	k := new(dsig.PrivateKey)
	if err := json.Unmarshal(data, k); err != nil {
		t.Errorf("failed to parse key: %v", err.Error())
		return
	}

	if err := k.Validate(); err != nil {
		t.Errorf("invalid key: %v", err.Error())
	}
	if k.Thumbprint() != thumb {
		t.Errorf("got unexpected thumbprint: %v", k.Thumbprint())
	}
	if k.ID() != kid {
		t.Errorf("unexpected key id, got: %v", k.ID())
	}

	pk := k.Public()
	if ty := reflect.TypeOf(pk); ty != reflect.TypeOf(&dsig.PublicKey{}) {
		t.Errorf("expected a public key")
	}
	if err := pk.Validate(); err != nil {
		t.Errorf("expected valid public key: %v", err.Error())
	}
	if pk.Thumbprint() != thumb {
		t.Errorf("expected same thumbprint, got: %v", pk.Thumbprint())
	}
	if pk.ID() != kid {
		t.Errorf("unexpected public key id, got: %v", k.ID())
	}
}
