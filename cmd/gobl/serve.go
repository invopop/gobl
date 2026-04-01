package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/invopop/gobl"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl/internal/api"
)

const defaultHTTPPort = 80

type serveOpts struct {
	httpPort int
}

func serve() *serveOpts {
	return &serveOpts{}
}

func (s *serveOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Launch an HTTP server",
		RunE:  s.runE,
	}
	f := cmd.Flags()
	f.IntVarP(&s.httpPort, "port", "p", defaultHTTPPort, "HTTP port to listen on")
	return cmd
}

func (s *serveOpts) runE(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(s.httpPort),
		Handler: api.NewHandler(),
	}

	addr := srv.Addr
	if addr == "" {
		addr = ":80"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "GOBL %s\n", gobl.VERSION)
	fmt.Fprintf(cmd.OutOrStdout(), "Listening on %s\n", addr)

	var startErr error
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			startErr = err
		}
		cancel()
	}()

	<-ctx.Done()
	fmt.Fprintln(cmd.OutOrStdout(), "Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return startErr
}
