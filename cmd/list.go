package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const listDesc = `
List checks
`

func newListChecksCmd() *cobra.Command {

	var flags operations.ListOperationFlags

	var cmdList = &cobra.Command{
		Use:    "list",
		Short:  "List all checks",
		Long:   listDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ListOperation{}).ListOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
