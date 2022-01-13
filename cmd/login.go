package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
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

			if err := (&operations.LoginOperation{}).LoginOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	cmdLogin.Flags().StringVarP(&flags.PersonalAccessToken, "personal-access-token", "t", "", "The personal access token. They can be created on checkson.io.")

	return cmdLogin
}
