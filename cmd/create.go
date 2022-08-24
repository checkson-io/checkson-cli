package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const createDesc = `
Create check
`

func newCreateCheckCmd() *cobra.Command {

	var flags operations.CreateCheckFlags

	var cmdCreateCheck = &cobra.Command{
		Use:    "create CHECK",
		Short:  "create a check",
		Long:   createDesc,
		Hidden: false,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode
			if err := (&operations.CreateCheckOperation{}).CreateCheckOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	cmdCreateCheck.Flags().StringVarP(&flags.WebHookUrl, "webhook-url", "w", "", "The webhook to call when the status of the check changes")
	cmdCreateCheck.Flags().StringVarP(&flags.WebHookUrl, "email", "m", "", "The email address to inform when the status of the check changes")
	cmdCreateCheck.Flags().StringVarP(&flags.DockerImage, "docker-image", "i", "", "The Docker image to execute")
	cmdCreateCheck.Flags().Int16VarP(&flags.CheckIntervalInMinutes, "check-interval", "d", 5, "The check interval in minutes")
	cmdCreateCheck.Flags().StringToStringVarP(&flags.Environment, "env", "e", map[string]string{}, "Environment variables")

	return cmdCreateCheck
}
