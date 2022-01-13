package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
)

const showDesc = `
show a check
`

func newShowCheckCmd() *cobra.Command {

	var flags operations.ShowOperationFlags

	var cmdList = &cobra.Command{
		Use:    "show CHECK",
		Short:  "show a check",
		Long:   showDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := (&operations.ShowOperation{}).ShowOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
