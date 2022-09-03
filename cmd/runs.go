package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const runDesc = `
List runs. If not check if given, list runs of all checks.
`

func newListRunsCmd() *cobra.Command {

	var flags operations.ListRunsOperationFlags

	var cmdRun = &cobra.Command{
		Use:    "runs [CHECK]",
		Short:  "list runs",
		Long:   runDesc,
		Hidden: false,
		Args:   cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			var checkName string
			if len(args) == 1 {
				checkName = args[0]
			} else {
				checkName = ""
			}

			if err := (&operations.ListRunsOperation{}).ListRunsOperation(checkName, flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdRun
}
