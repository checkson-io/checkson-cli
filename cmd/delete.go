package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
)

const deleteDesc = `
Delete checks
`

func newDeleteCheckCmd() *cobra.Command {

	var flags operations.DeleteCheckFlags

	var cmdDeleteCheck = &cobra.Command{
		Use:    "delete CHECK",
		Short:  "delete a check",
		Long:   deleteDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if err := (&operations.DeleteCheckOperation{}).DeleteCheckOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdDeleteCheck
}
