package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
)

const logsDesc = `
Show log output. By default the most recent run of a job is shown.
`

func newLogsCmd() *cobra.Command {

	var flags operations.LogsOperationFlags

	var cmdLogs = &cobra.Command{
		Use:    "logs",
		Short:  "Show log output of runs",
		Long:   logsDesc,
		Hidden: false,
		Args:   cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
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
