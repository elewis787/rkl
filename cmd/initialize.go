package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elewis787/rcl/internal/tui"
	"github.com/spf13/cobra"
)

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialze",
		Short:   "init the rcl cfg.",
		Long:    "init provision the rcl configuration file.",
		Example: "rkl init",
		Aliases: []string{"i", "init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := tea.NewProgram(tui.NewInitPrompt()).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	return init
}
