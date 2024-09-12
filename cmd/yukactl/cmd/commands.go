package cmd

import (
	"yuka/cmd/yukactl/cmd/client"

	"github.com/spf13/cobra"
)

var subcommands = []*cobra.Command{
	client.SubCommand(),
}

func init() {
	for _, subcommand := range subcommands {
		rootCmd.AddCommand(subcommand)
	}
}
