package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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
			}
		},
	}

	return cmdDeleteCheck
}
