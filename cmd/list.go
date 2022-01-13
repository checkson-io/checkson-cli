package cmd

import (
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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
