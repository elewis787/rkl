package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elewis787/rkl/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialize",
		Short:   "init the rcl cfg.",
		Long:    "init provision the rcl configuration file.",
		Example: "rkl init",
		Aliases: []string{"i", "init"},
		// used to overwrite/skip the parent commands persistentPreRunE func
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind Cobra flags with viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			// Environment variables are expected to be ALL CAPS
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rkl")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			if err := tea.NewProgram(tui.NewInitPrompt(viper.GetString(cfgPath), homeDir)).Start(); err != nil {
				return err
			}
			return nil
		},
	}
	return init
}
