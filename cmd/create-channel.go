package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const createChannelDesc = `
Create check
`

func newCreateChannelCmd() *cobra.Command {

	var flags operations.CreateChannelFlags

	var cmdCreateChannel = &cobra.Command{
		Use:    "create CHANNEL",
		Short:  "create a notification channel",
		Long:   createChannelDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode
			if err := (&operations.CreateChannelOperation{}).CreateChannelOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	cmdCreateChannel.Flags().StringVarP(&flags.WebHookUrl, "webhook-url", "w", "", "The webhook to call when the status of the check changes")
	cmdCreateChannel.Flags().StringVarP(&flags.WebHookUrl, "email", "e", "", "The email address to inform when the status of the check changes")

	return cmdCreateChannel
}
