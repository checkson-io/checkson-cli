package cmd

import (
	"github.com/checkson-io/checkson-cli/operations"
	"github.com/checkson-io/checkson-cli/output"
	"github.com/spf13/cobra"
)

const statusDesc = `
List checks, alias for list
`

func newStatusChecksCmd() *cobra.Command {

	var flags operations.ListOperationFlags

	var cmdStatus = &cobra.Command{
		Use:    "status",
		Short:  "Status all checks, alias for list",
		Long:   statusDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ListOperation{}).ListOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdStatus
}
