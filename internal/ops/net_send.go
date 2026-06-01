package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/net"
)

const netSendTimeout = 10 * time.Second

// NetSendOptions configures the gobl net send command.
type NetSendOptions struct {
	Input    io.Reader
	To       net.Address
	Insecure bool        // when true: use http:// and accept host:port form
	Client   *http.Client // optional; defaults to a 10s-timeout client
}

// NetSend reads a GOBL envelope from opts.Input and POSTs it to the
// destination address's inbox endpoint. Returns ErrInboxRejected if
// the inbox does not respond with 202.
func NetSend(ctx context.Context, opts *NetSendOptions) error {
	body, err := io.ReadAll(cancelableReader(ctx, opts.Input))
	if err != nil {
		return gobl.ErrInput.WithCause(err)
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(body, env); err != nil {
		return gobl.ErrInput.WithCause(err)
	}
	if err := env.Validate(); err != nil {
		return gobl.ErrValidation.WithCause(err)
	}

	url, err := inboxURL(opts.To, opts.Insecure)
	if err != nil {
		return err
	}

	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: netSendTimeout}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("net send: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("net send: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("%w: HTTP %d: %s", net.ErrInboxRejected, resp.StatusCode, bytes.TrimSpace(respBody))
}

func inboxURL(addr net.Address, insecure bool) (string, error) {
	if insecure {
		if addr == "" {
			return "", net.ErrAddressEmpty
		}
		return "http://" + string(addr) + net.InboxPath, nil
	}
	if err := addr.Validate(); err != nil {
		return "", err
	}
	return addr.InboxURL(), nil
}
