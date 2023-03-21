package statics

import (
	"{{ cookiecutter.project_slug }}"
	_ "embed"
	"strings"
)

//go:embed description.md
var description string

//go:embed debugger.html
var DebuggerHTML string

func GenDescription() string {
	replacer := strings.NewReplacer("{build_time}", {{ cookiecutter.project_slug }}.BuildTime, "{go_version}", {{ cookiecutter.project_slug }}.GoVersion)
	return replacer.Replace(description)
}
