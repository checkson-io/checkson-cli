package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
)

const createDesc = `
List checks
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

			if err := (&operations.CreateCheckOperation{}).CreateCheckOperation(args[0], flags); err != nil {
				output.Fail(err)
			}
		},
	}

	cmdCreateCheck.Flags().StringVarP(&flags.WebHookUrl, "webhook-url", "w", "", "The webhook to call when the status of the check changes")
	cmdCreateCheck.Flags().StringVarP(&flags.DockerImage, "docker-image", "i", "", "The Docker image to execute")
	cmdCreateCheck.Flags().Int16VarP(&flags.CheckIntervalInMinutes, "check-interval", "d", 5, "The check interval in minutes")

	return cmdCreateCheck
}
