package main

import (
	"{{ cookiecutter.project_slug }}/cmd/{{ cookiecutter.project_slug }}/scripts"
	"{{ cookiecutter.project_slug }}/cmd/{{ cookiecutter.project_slug }}/server"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "{{ cookiecutter.project_slug }}-cli",
	Short: "The `{{ cookiecutter.project_slug }}` CLI",
}

func main() {
	rootCmd.AddCommand(server.CmdRun)
	rootCmd.AddCommand(scripts.CmdScripts)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute root command")
	}
}
