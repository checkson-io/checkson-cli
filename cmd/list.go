package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
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
			if err := (&operations.ListOperation{}).ListOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
