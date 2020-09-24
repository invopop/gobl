package gobl_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl"
)

func TestNewECDSAKey(t *testing.T) {
	// With ID
	k, err := gobl.NewECDSAKey("key1")
	if err != nil {
		t.Errorf("did not expect error: %v", err)
	}
	if k.ID() != "key1" {
		t.Errorf("failed to assign key ID")
	}

	// Without ID
	k, err = gobl.NewECDSAKey("")
	if err != nil {
		t.Errorf("did not expect error")
	}
	if k.ID() == "" {
		t.Errorf("did not expect empty key ID")
	}

	if k.IsPublic() {
		t.Errorf("expected private key")
	}
	if !k.Valid() {
		t.Errorf("did not expect key to be invalid")
	}

}

func TestParseKey(t *testing.T) {
	data := []byte(`{"use":"sig","kty":"EC","kid":"5186c2af-606f-46d9-85c5-16ea0d5e23b3","crv":"P-256","x":"k7C64o1DiRZlF_e-8XGUxAWYmhThwKeM8yoU_YipfKQ","y":"GBcnnygic7LVfu0-2LjCVG7ieuoH-6EVMaEcpuTVW04","d":"46iMLaPQn9msYKjuCvlu1sS5p_qLC1DxZw-7jWo79r0"}`)
	thumb := "47dce63f7d988e81015563269ac6652df65295eaca07608349b4e2860a232637"

	k := new(gobl.Key)
	if err := json.Unmarshal(data, k); err != nil {
		t.Errorf("Failed to parse key: %v", err.Error())
		return
	}

	th, err := k.Thumbprint()
	if err != nil {
		t.Errorf("did not expect error from thumbprint, got: %v", err)
		return
	}
	if th != thumb {
		t.Errorf("invalid thumbprint, got: %v", th)
	}

}
