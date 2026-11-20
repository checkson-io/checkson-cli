package cmd

import (
	"github.com/checkson-io/checkson-cli/operations"
	"github.com/checkson-io/checkson-cli/output"
	"github.com/spf13/cobra"
)

const listChannelsDesc = `
List checks
`

func newListChannelsCmd() *cobra.Command {

	var flags operations.ListChannelsOperationFlags

	var cmdList = &cobra.Command{
		Use:    "list",
		Short:  "List all notification channels",
		Long:   listChannelsDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ListChannelsOperation{}).ListChannelsOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
