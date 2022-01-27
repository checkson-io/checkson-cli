package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const loginDesc = `
Login to checkson
`

func newLoginCmd() *cobra.Command {

	// TODO: Add flag: --personal-access-token
	// Login via curl -X POST -H Content-Type:application/json -d '{"personalAccessToken": "xxx"}' http://.../europe-west1/getCustomAuthTokenForPersonalAccessToken

	var flags operations.LoginOperationFlags

	var cmdLogin = &cobra.Command{
		Use:    "login",
		Short:  "Login to checkson",
		Long:   loginDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			devMode, _ := cmd.Flags().GetBool("dev-mode")
			flags.DevMode = devMode

			if err := (&operations.LoginOperation{}).LoginOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	cmdLogin.Flags().StringVarP(&flags.PersonalAccessToken, "personal-access-token", "t", "", "The personal access token. They can be created on checkson.io.")

	return cmdLogin
}
