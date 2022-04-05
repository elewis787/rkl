package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/elewis787/rcl/cmd/history"
	"github.com/elewis787/rcl/cmd/note"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "rcl",
		Short: "recall",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			viper.AutomaticEnv()
			viper.SetEnvPrefix("rcl")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cfgPath := homeDir + rclDir + cfgFile
			if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
				return errors.New(err.Error() + ", please run init to configure rcl\n")
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().String("", "", "defaults to active channel address in the cfg")
	rootCmd.SetUsageFunc(styleUsageFunc)
	rootCmd.SetHelpFunc(styleHelpFunc)
	rootCmd.AddCommand(history.HistoryCmd())
	rootCmd.AddCommand(note.NoteCmd())

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

func styleUsageFunc(c *cobra.Command) error {
	useTitle := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
		Render("Usage:")
	fmt.Println(useTitle)
	return nil
}

func buildUsage(c *cobra.Command) string {
	var useTitle, useBody, subCommandPath,
		aliases, examples, commands string

	useTitle = lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
		Render("Usage:")

	if c.Runnable() {
		useBody = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render(c.Use)
	}

	if c.HasAvailableSubCommands() {
		subCommandPath = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render(c.CommandPath() + "[command]")
	}

	if len(c.Aliases) > 0 {
		aliases = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render("Aliases:\n")
		aliases += lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render(c.NameAndAliases())
	}

	if c.HasExample() {
		examples = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render("Examples:\n")
		examples += lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render(c.Example)
	}

	if c.HasSubCommands() {
		commands = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
			PaddingRight(5).Render("Available Commands:\n")
		for _, subcmd := range c.Commands() {
			if subcmd.Name() == "help" || subcmd.IsAvailableCommand() {
				commands += lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
					PaddingRight(5).Render(c.Name() + " " + c.Short)
			}
		}
	}

	fmt.Println(useTitle, useBody, subCommandPath, aliases, examples, commands)

	return ""
}

func styleHelpFunc(c *cobra.Command, s []string) {

	if c.HasAvailableSubCommands() {

	}
	for _, r := range c.Commands() {

		fmt.Println(r.Use)
	}

	fs := c.Flags().FlagUsages()

	fmt.Println(fs)
	fmt.Println(s)

}
