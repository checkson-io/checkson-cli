package cmd

import (
	"github.com/checkson-io/checkson-cli/operations"
	"github.com/checkson-io/checkson-cli/output"
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

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.DeleteCheckOperation{}).DeleteCheckOperation(args[0], flags); err != nil {
				output.Fail(err)
			} else {
				output.PrintStrings("Check deleted")
			}
		},
	}

	return cmdDeleteCheck
}
