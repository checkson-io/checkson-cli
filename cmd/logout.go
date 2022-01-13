package cmd

import (
	"github.com/huditech/checkson/operations"
	"github.com/huditech/checkson/output"
	"github.com/spf13/cobra"
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
