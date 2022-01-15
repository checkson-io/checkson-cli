package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const runDesc = `
List runs
`

func newListRunsCmd() *cobra.Command {

	var flags operations.ListRunsOperationFlags

	var cmdRun = &cobra.Command{
		Use:    "runs",
		Short:  "list all runs",
		Long:   runDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ListRunsOperation{}).ListRunsOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdRun
}
