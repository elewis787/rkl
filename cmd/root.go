package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/elewis787/boa"
	"github.com/elewis787/rkl/cmd/history"
	"github.com/elewis787/rkl/cmd/note"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Execute is the command line applications entry function
func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "rkl",
		Long:    "Rekall (rkl) is a CLI that helps you remember things. Easily manage past commands, todos and notes all from your command line.",
		Example: "rkl",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cfg := homeDir + cfgDir + cfgFile
			if _, err := os.Stat(cfg); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error() + ": please run init to configure rkl\n")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rkl")
			return nil
		},
	}

	rootCmd.SetUsageFunc(boa.UsageFunc)
	rootCmd.SetHelpFunc(boa.HelpFunc)

	// Add sub commands
	rootCmd.AddCommand(initialize())
	rootCmd.AddCommand(history.HistoryCmd())
	rootCmd.AddCommand(note.NoteCmd())

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rootCmd.PersistentFlags().String(cfgPath, dir+cfgDir+cfgFile, "location of the mazzaroth config file")
	// Applicaiton execution
	ctx, cancel := context.WithCancel(context.Background())
	errGrp, errctx := errgroup.WithContext(ctx)
	errGrp.Go(func() error {
		defer cancel()
		if err := rootCmd.ExecuteContext(errctx); err != nil {
			return err
		}
		return nil
	})

	return errGrp.Wait()
}
