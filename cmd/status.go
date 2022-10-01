package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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
