package cmd

import (
	"github.com/spf13/cobra"
)

func newChannelCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "channels",
		Short: "manage notification channels",
	}

	cmd.AddCommand(newCreateChannelCmd())
	cmd.AddCommand(newDeleteChannelCmd())

	return cmd
}
