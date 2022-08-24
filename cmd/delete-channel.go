package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const deleteChannelDesc = `
Delete channel
`

func newDeleteChannelCmd() *cobra.Command {

	var flags operations.DeleteChannelFlags

	var cmdDeleteChannel = &cobra.Command{
		Use:    "delete CHANNEL",
		Short:  "delete a notification channel",
		Long:   deleteChannelDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode
			if err := (&operations.DeleteChannelOperation{}).DeleteChannelOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdDeleteChannel
}
