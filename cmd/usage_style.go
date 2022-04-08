package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	purple    = `#7e2fcc`
	darkGrey  = `#353C3B`
	lightTeal = `#03DAC5`
	darkTeal  = `#01A299`
	white     = `#e5e5e5`
	red       = `#F45353`
)

var (
	// TODO func to calc width
	physicalWidth, _, _ = term.GetSize(int(os.Stdout.Fd()))
	width               = physicalWidth / 3

	titleStyle = lipgloss.NewStyle().Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: purple, Dark: purple}).
			Padding(1, 1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
			Underline(true).
			BorderBottom(true).
			Padding(0, 1, 0, 1).Align(lipgloss.Center)

	textStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 5).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})

	subTextStyle = lipgloss.NewStyle().Bold(true).Padding(0, 0, 0, 2).
			Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white})
)

func styleUsageFunc(c *cobra.Command) error {

	return nil
}

func buildUsage(c *cobra.Command) string {
	// TODO switch to strings builder
	// output := strings.Builder{}
	fmt.Println(physicalWidth)
	cmdTitle := ""

	if !c.HasParent() {
		rootCmdName := sectionStyle.Render(c.Root().Name() + " " + c.Root().Version)
		rootCmdLong := lipgloss.NewStyle().Align(lipgloss.Center).Render(c.Root().Long)

		cmdTitle = titleStyle.Width(width).Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).Render(lipgloss.JoinVertical(lipgloss.Top, rootCmdName, rootCmdLong))
	}

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
				cmd := textStyle.Render(subcmd.Name()) +
					lipgloss.NewStyle().
						Foreground(lipgloss.AdaptiveColor{Light: darkGrey, Dark: white}).Bold(true).
						PaddingLeft(subcmd.NamePadding()-len(subcmd.Name())+1).Render(subcmd.Short)
				usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, cmd)
			}
		}
		subCmdHelp := subTextStyle.Render("\nUse \"" + c.CommandPath() + " [command] --help\" for more information about a command.")
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, subCmdHelp)
	}

	if c.HasAvailableLocalFlags() {
		localFlags := sectionStyle.Render("Flags:")
		flagUsage := textStyle.Render(strings.TrimFunc(c.LocalFlags().FlagUsages(), unicode.IsSpace))
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, localFlags, flagUsage)
	}

	if c.HasAvailableInheritedFlags() {
		globalFlags := sectionStyle.Render("Global Flags:")
		flagUsage := textStyle.Render(strings.TrimFunc(c.InheritedFlags().FlagUsages(), unicode.IsSpace))
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, globalFlags, flagUsage)
	}

	if c.HasExample() {
		examples := sectionStyle.Render("Examples:")
		example := textStyle.Render(c.Example)
		usageOutput = lipgloss.JoinVertical(lipgloss.Top, usageOutput, examples, example)
	}
	usageOutput = lipgloss.JoinVertical(lipgloss.Top, cmdTitle, usageOutput)
	usageOutput = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		BorderForeground(lipgloss.AdaptiveColor{Light: darkTeal, Dark: lightTeal}).
		Border(lipgloss.ThickBorder()).Render(usageOutput)

	return usageOutput
}

func styleHelpFunc(c *cobra.Command, s []string) {
	fmt.Println(buildUsage(c))
}
