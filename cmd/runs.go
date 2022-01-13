package cmd

import (
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/spf13/cobra"
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
			if err := (&operations.ListRunsOperation{}).ListRunsOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdRun
}
