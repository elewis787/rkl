package cmd

import "github.com/spf13/cobra"

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "initialze",
		Short:   "init the rcl cfg.",
		Long:    "init provision the rcl configuration file.",
		Example: "rcl init",
		Aliases: []string{"i", "init"},
	}
	return init
}
