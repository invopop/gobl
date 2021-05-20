package dsig_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/dsig"
)

func TestNewKey(t *testing.T) {
	k := dsig.NewES256Key()

	if k.IsPublic() {
		t.Errorf("expected private key")
	}
	if k.ID() == "" {
		t.Errorf("expected an ID")
	}
	if !k.Valid() {
		t.Errorf("expected a valid new key")
	}
	if k.Thumbprint() == "" {
		t.Errorf("expected a thumbprint")
	}
	data, err := json.Marshal(k)
	if err != nil {
		t.Errorf("failed to generate JSON data: %v", err.Error())
	}
	t.Logf("json output of new key: %v", string(data))
}

func TestKeyParsing(t *testing.T) {
	data := []byte(`{"use":"sig","kty":"EC","kid":"3500bbee-966c-4b7a-8fbc-c763ae2aec62","crv":"P-256","x":"Fd4a9pj2gtDLnW3GX30S06qXHrkBrAsmg3aHb4kOCL4","y":"_I4ZuddZtZ86kDBvGKcsOPbU0gWh13Kt6R2m6bfWAK4","d":"oJM3Ogl9uYUpSbc4oHV25DpFs_gOGP5nHJcLAtQxL6U"}`)
	kid := "3500bbee-966c-4b7a-8fbc-c763ae2aec62"
	thumb := "eb9793f9b542e25e13959591a41391329f2eb05a0d8f9b12eb071b4d90d91c51"

	k := new(dsig.Key)
	if err := json.Unmarshal(data, k); err != nil {
		t.Errorf("failed to parse key: %v", err.Error())
		return
	}

	if k.IsPublic() {
		t.Errorf("did not expect public key")
	}
	if !k.Valid() {
		t.Errorf("did no expect invalid key")
	}
	if k.Thumbprint() != thumb {
		t.Errorf("got unexpected thumbprint: %v", k.Thumbprint())
	}
	if k.ID() != kid {
		t.Errorf("unexpected key id, got: %v", k.ID())
	}

	pk := k.Public()
	if !pk.Valid() {
		t.Errorf("expected valid public key")
	}
	if !pk.IsPublic() {
		t.Errorf("expected public key")
	}
	if pk.Thumbprint() != thumb {
		t.Errorf("expected same thumbprint, got: %v", pk.Thumbprint())
	}
	if pk.ID() != kid {
		t.Errorf("unexpected public key id, got: %v", k.ID())
	}
}
