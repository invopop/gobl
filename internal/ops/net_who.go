package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
)

// schemeRewriteFetcher rewrites well-known https://<addr>/... URLs to
// the given http:// base so --insecure mode reuses net.Client logic
// (which always builds https URLs) over plain HTTP.
type schemeRewriteFetcher struct {
	base  string // e.g. http://acme.example
	inner net.Fetcher
}

func (s *schemeRewriteFetcher) Fetch(ctx context.Context, raw string) ([]byte, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return s.inner.Fetch(ctx, raw)
	}
	bu, err := url.Parse(s.base)
	if err == nil {
		u.Scheme = bu.Scheme
		u.Host = bu.Host
		raw = u.String()
	}
	return s.inner.Fetch(ctx, raw)
}

const netWhoTimeout = 10 * time.Second

// NetWhoOptions configures NetWho.
type NetWhoOptions struct {
	Target    net.Address       // domain being queried
	From      net.Address       // caller's GOBL Net address (signs the request)
	FromKey   *dsig.PrivateKey  // caller's signing key
	FromParty *org.Party        // caller's party, sent as the request document
	Insecure  bool              // query over http:// and permit host:port
	Fetcher   net.Fetcher       // optional (for /keys); defaults to net.NewHTTPFetcher()
	Client    *http.Client      // optional (for POST /who); defaults to 10s timeout
}

// NetWho performs an authenticated GOBL Net party exchange: it POSTs a
// signed request envelope (the caller's party, iss=gobl:from,
// aud=gobl:target) to the target's /who endpoint, verifies the response
// is signed by the target (iss=gobl:target) and bound to the caller
// (aud=gobl:from), and returns the target's verified org.Party.
func NetWho(ctx context.Context, opts *NetWhoOptions) (*org.Party, error) {
	if opts.Target == "" {
		return nil, gobl.ErrInput.WithReason("target address is required")
	}
	if opts.From == "" || opts.FromKey == nil || opts.FromParty == nil {
		return nil, gobl.ErrInput.WithReason("a --from identity (key + party) is required to authenticate the request")
	}

	scheme := "https"
	if opts.Insecure {
		scheme = "http"
	}
	base := scheme + "://" + string(opts.Target)

	// Build and sign the request envelope: iss=from, aud=target.
	reqEnv, err := gobl.Envelop(opts.FromParty)
	if err != nil {
		return nil, fmt.Errorf("net who: build request: %w", err)
	}
	if err := reqEnv.Sign(opts.FromKey, opts.From.URI(), opts.Target.URI()); err != nil {
		return nil, fmt.Errorf("net who: sign request: %w", err)
	}
	reqBody, err := json.Marshal(reqEnv)
	if err != nil {
		return nil, fmt.Errorf("net who: encode request: %w", err)
	}

	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: netWhoTimeout}
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, base+net.WhoPath, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("net who: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("net who: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, netInboxMaxBody))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("net who: %w: HTTP %d: %s", net.ErrFetchFailed, resp.StatusCode, bytes.TrimSpace(respBody))
	}

	respEnv := new(gobl.Envelope)
	if err := json.Unmarshal(respBody, respEnv); err != nil {
		return nil, fmt.Errorf("net who: invalid /who response: %w", err)
	}
	if !respEnv.Signed() {
		return nil, fmt.Errorf("net who: /who response is not signed")
	}

	// Verify the response is signed by the target, bound to us. The
	// fetcher is wrapped so /key/<kid> URLs honour --insecure by being
	// rewritten to the http:// base.
	fetcher := opts.Fetcher
	if fetcher == nil {
		fetcher = net.NewHTTPFetcher()
	}
	if opts.Insecure {
		fetcher = &schemeRewriteFetcher{base: base, inner: fetcher}
	}
	verifyClient := net.NewClient(net.WithFetcher(fetcher))

	wantIss := opts.Target.URI()
	wantAud := opts.From.URI()
	verified := false
	for _, sig := range respEnv.Signatures {
		p, perr := head.SignedPayload(sig)
		if perr != nil || p.Iss != wantIss {
			continue
		}
		pubKey, kerr := verifyClient.FetchKey(ctx, opts.Target, sig.KeyID())
		if kerr != nil {
			continue
		}
		// VerifySignature enforces the key's validity window via
		// head.Header.Verify, so no extra Allows call is needed here.
		if respEnv.VerifySignature(sig, pubKey) != nil {
			continue
		}
		if p.Aud != "" && p.Aud != wantAud {
			return nil, fmt.Errorf("net who: response audience mismatch (got %q, want %q)", p.Aud, wantAud)
		}
		verified = true
		break
	}
	if !verified {
		return nil, fmt.Errorf("net who: response not signed by %s with a published key", wantIss)
	}

	party, ok := respEnv.Extract().(*org.Party)
	if !ok || party == nil {
		return nil, fmt.Errorf("net who: /who response document is not an org.Party")
	}
	return party, nil
}
