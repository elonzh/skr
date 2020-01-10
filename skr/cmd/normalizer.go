package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

const Indentation = `  `

// LongDesc normalizes a command's long description to follow the conventions.
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().string
}

// Examples normalizes a command's examples to follow the conventions.
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().indent().string
}

// Normalize perform all required normalizations on a given command.
func Normalize(cmd *cobra.Command) *cobra.Command {
	if len(cmd.Long) > 0 {
		cmd.Long = LongDesc(cmd.Long)
	}
	if len(cmd.Example) > 0 {
		cmd.Example = Examples(cmd.Example)
	}
	return cmd
}

// NormalizeAll perform all required normalizations in the entire command tree.
func NormalizeAll(cmd *cobra.Command) *cobra.Command {
	if cmd.HasSubCommands() {
		for _, subCmd := range cmd.Commands() {
			NormalizeAll(subCmd)
		}
	}
	Normalize(cmd)
	return cmd
}

type normalizer struct {
	string
}

func (s normalizer) trim() normalizer {
	s.string = strings.TrimSpace(s.string)
	return s
}

func (s normalizer) indent() normalizer {
	var indentedLines []string
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}
