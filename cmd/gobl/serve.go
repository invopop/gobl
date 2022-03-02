package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cmd/gobl/internal"
	"github.com/invopop/gobl/dsig"
)

const (
	defaultHTTPPort = 80

	// If you customize this server, you should change this.
	vendorName = "Invopop Ltd."
)

type serveOpts struct {
	httpPort int
}

func serve() *serveOpts {
	return &serveOpts{}
}

func (s *serveOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "serve",
		RunE: s.runE,
	}
	f := cmd.Flags()

	f.IntVarP(&s.httpPort, "port", "p", defaultHTTPPort, "HTTP port to listen on")

	return cmd
}

func (s *serveOpts) runE(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	e := echo.New()

	e.GET("/", s.version())
	e.POST("/build", s.build())
	e.POST("/verify", s.verify())
	e.GET("/keygen", s.keygen())

	var startErr error
	go func() {
		err := e.Start(":" + strconv.Itoa(s.httpPort))
		if !errors.Is(err, http.ErrServerClosed) {
			startErr = err
		}
		cancel() // Ensure we stop the context if we have a startup error
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return startErr
}

func (s *serveOpts) version() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"gobl":    "Welcome",
			"version": gobl.VERSION,
			"vendor": map[string]interface{}{
				"name": vendorName,
			},
		})
	}
}

type buildRequest struct {
	Template json.RawMessage `json:"template"`
	Data     json.RawMessage `json:"data"`
}

func (s *serveOpts) build() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ct, _, _ := mime.ParseMediaType(c.Request().Header.Get("Content-Type"))
		if ct != "application/json" {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType)
		}
		req := new(buildRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		if len(req.Data) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "no payload")
		}
		opts := internal.BuildOptions{
			Data: bytes.NewReader(req.Data),
		}
		if len(req.Template) != 0 {
			opts.Template = bytes.NewReader(req.Template)
		}
		env, err := internal.Build(ctx, opts)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, env)
	}
}

type verifyRequest struct {
	Data json.RawMessage `json:"data"`
}

type verifyResponse struct {
	OK bool `json:"ok"`
}

func (s *serveOpts) verify() echo.HandlerFunc {
	return func(c echo.Context) error {
		ct, _, _ := mime.ParseMediaType(c.Request().Header.Get("Content-Type"))
		if ct != "application/json" {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType)
		}
		req := new(verifyRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		if err := internal.Verify(c.Request().Context(), bytes.NewReader(req.Data)); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &verifyResponse{OK: true})
	}
}

type keygenResponse struct {
	Private *dsig.PrivateKey `json:"private"`
	Public  *dsig.PublicKey  `json:"public"`
}

func (s *serveOpts) keygen() echo.HandlerFunc {
	return func(c echo.Context) error {
		key := dsig.NewES256Key()

		return c.JSON(http.StatusOK, keygenResponse{
			Private: key,
			Public:  key.Public(),
		})
	}
}
