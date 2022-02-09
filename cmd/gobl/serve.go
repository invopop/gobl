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

	"github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cmd/gobl/internal"
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
	Data json.RawMessage `json:"data"`
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
		env, err := internal.Build(ctx, internal.BuildOptions{
			Data: bytes.NewReader(req.Data),
		})
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
		env := new(gobl.Envelope)
		if err := yaml.Unmarshal(req.Data, env); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := env.Validate(); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		return c.JSON(http.StatusOK, &verifyResponse{OK: true})
	}
}
