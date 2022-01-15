package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ShowOperation{}).ShowOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
