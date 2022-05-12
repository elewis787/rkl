package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/elewis787/boa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute is the command line applications entry function
func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "rkl",
		Long:    "Rekall (rkl) is a CLI that helps you remember things. Easily manage past commands, todos and notes all from your command line.",
		Example: "rkl",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rkl")

			if _, err := os.Stat(viper.GetString(cfgPath)); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error() + ": please run init to configure rkl\n")
			}
			return nil
		},
	}

	rootCmd.SetHelpFunc(boa.HelpFunc)
	rootCmd.SetUsageFunc(boa.UsageFunc)

	// Add sub commands
	rootCmd.AddCommand(initialize())

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	//Define root flags
	rootCmd.PersistentFlags().String(cfgPath, dir+cfgDir+cfgFile, "location of the rkl config file")

	return rootCmd.ExecuteContext(context.Background())
}
