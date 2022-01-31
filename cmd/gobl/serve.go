package main

import "github.com/spf13/cobra"

type serveOpts struct {
}

func serve() *serveOpts {
	return &serveOpts{}
}

func (s *serveOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "serve",
		RunE: s.runE,
	}
	return cmd
}

func (s *serveOpts) runE(*cobra.Command, []string) error {
	return nil
}
