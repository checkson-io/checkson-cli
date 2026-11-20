package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
)

func newEditCheckCmd() *cobra.Command {
	var filename string
	command := &cobra.Command{
		Use:   "edit CHECK --file FILE",
		Short: "replace a check configuration (all fields required)",
		Long: `Replace an existing check's configuration. Requires a complete JSON configuration;
no merge is performed. All fields must be explicit, including empty values.

Required fields: dockerImage, enabled, dockerCredentials, checkIntervalInMinutes,
environment, description, failureThreshold, notificationChannels.
Use null for unused dockerCredentials, {} for empty environment, and [] for no channels.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			devMode, _ := cmd.Flags().GetBool("dev-mode")
			if err := operations.EditCheck(args[0], filename, devMode); err != nil {
				return err
			}
			cmd.Println("Check updated")
			return nil
		},
	}
	command.Flags().StringVarP(&filename, "file", "f", "", "Complete check configuration (JSON)")
	_ = command.MarkFlagRequired("file")
	return command
}
