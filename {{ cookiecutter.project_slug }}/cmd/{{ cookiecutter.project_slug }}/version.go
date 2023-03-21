package main

import (
	"fmt"

	"{{ cookiecutter.project_slug }}"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "print the version info",
	Run: func(*cobra.Command, []string) {
		fmt.Println("version:", {{ cookiecutter.project_slug }}.Version)
		fmt.Println("build time:", {{ cookiecutter.project_slug }}.BuildTime)
		fmt.Println("build go version:", {{ cookiecutter.project_slug }}.GoVersion)
	},
}

func init() { // nolint
	rootCmd.AddCommand(cmdVersion)
}
