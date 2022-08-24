package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/operations"
	"github.com/stefan-hudelmaier/checkson-cli/output"
)

const logoutDesc = `
Logout from checkson
`

func newLogoutCmd() *cobra.Command {

	var flags operations.LogoutOperationFlags

	var cmdLogout = &cobra.Command{
		Use:    "logout",
		Short:  "Logout to checkson",
		Long:   logoutDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			if err := (&operations.LogoutOperation{}).LogoutOperation(flags); err != nil {
				output.Fail(err)
			}
		},
	}

	return cmdLogout
}
