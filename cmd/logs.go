package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const logsDesc = `
Show log output. By default the most recent run of a job is shown.
`

func newLogsCmd() *cobra.Command {

	var flags operations.LogsOperationFlags

	var cmdLogs = &cobra.Command{
		Use:    "logs CHECK [RUN_ID]",
		Short:  "Show log output of runs",
		Long:   logsDesc,
		Hidden: false,
		Args:   cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			var checkName string
			if len(args) == 1 {
				checkName = "_last"
			} else {
				checkName = args[1]
			}
			if err := (&operations.LogsOperation{}).LogsOperation(args[0], checkName, flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdLogs
}
