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

const (
	gold     = `#E3BD2D`
	darkGrey = `#353C3B`
	teal     = `#01A299`
	white    = `#e5e5e5`
	red      = `#F31849`

	width = 50
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.ThickBorder(), true, true, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: gold, Dark: gold}).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
			MarginLeft(1).
			Padding(0, 1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: teal, Dark: teal}).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: gold, Dark: gold}).
			BorderBottom(true).
			Padding(1, 1, 0, 1).Align(lipgloss.Center)

	textStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	subTextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 2).
			Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white})
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:     "rcl",
		Short:   "rcl root command",
		Example: "rcl [sub command] [flags]",
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

	rootCmd.PersistentFlags().String("yolo", "", "defaults to active channel address in the cfg")
	//rootCmd.SetUsageFunc(styleUsageFunc)
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

	return nil
}

func buildUsage(c *cobra.Command) string {
	cmdTitle := titleStyle.Render(c.Root().Name())

	usageOutput := sectionStyle.Render("Cmd Description:")
	short := textStyle.Render(c.Short)

	usage := sectionStyle.Render("Usage:")
	usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, short, usage)

	if c.Runnable() {
		useLine := textStyle.Render(c.UseLine())
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, useLine)
	}

	if c.HasAvailableSubCommands() {
		commandPath := textStyle.Render(c.CommandPath() + " [command]")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, commandPath)
	}

	if len(c.Aliases) > 0 {
		aliases := sectionStyle.Render("Aliases:")
		nameAndAlias := textStyle.Render(c.NameAndAliases())
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, aliases, nameAndAlias)

	}

	if c.HasAvailableSubCommands() {
		commands := sectionStyle.Render("Available Commands:")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, commands)

		for _, subcmd := range c.Commands() {
			if subcmd.Name() == "help" || subcmd.IsAvailableCommand() {
				cmd := textStyle.Render(subcmd.Name()) + lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: white, Dark: white}).Bold(true).
					PaddingLeft(subcmd.NamePadding()-len(subcmd.Name())+1).Render(subcmd.Short)
				usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, cmd)
			}
		}
		subCmdHelp := subTextStyle.Render("\nUse \"" + c.CommandPath() + " [command] --help\" for more information about a command.")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, subCmdHelp)
	}

	if c.HasAvailableLocalFlags() {
		localFlags := sectionStyle.Render("Flags:")
		flagUsage := textStyle.Render(c.LocalFlags().FlagUsages())
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, localFlags, flagUsage)
	}
	/*
		if c.HasAvailableInheritedFlags() {
			globalFlags = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
				PaddingRight(5).Render("Global Flags:\n")
			globalFlags += lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353C3B", Dark: "#FFFFFF"}).
				PaddingRight(5).Render(c.InheritedFlags().FlagUsages())
		}

		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usage, subCommandPath, aliases, examples, commands, localFlags, globalFlags)
	*/

	if c.HasExample() {
		examples := sectionStyle.Render("Examples:")
		example := textStyle.Render(c.Example)
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, examples, example)
	}

	usageOutput = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).
		Border(lipgloss.ThickBorder()).Render(usageOutput)

	return lipgloss.JoinVertical(0, cmdTitle, usageOutput)
}

func styleHelpFunc(c *cobra.Command, s []string) {
	fmt.Println(buildUsage(c))
}
