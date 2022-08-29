package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const showChannelDesc = `
show a channel
`

func newShowChannelCmd() *cobra.Command {

	var flags operations.ShowChannelFlags

	var cmdList = &cobra.Command{
		Use:    "show CHANNEL",
		Short:  "show a channel",
		Long:   showChannelDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.ShowChannelOperation{}).ShowChannelOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdList
}
