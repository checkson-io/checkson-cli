package cmd

import (
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/config"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

const loginDebugDesc = `
Login to checkson (debug mode - must be removed before release)
`

type LoginDebugFlags struct {
	OrgId string
}

func newLoginDebugCmd() *cobra.Command {

	//var flags LoginDebugFlags

	var cmdLoginDebug = &cobra.Command{
		Use:    "logindebug",
		Short:  "Login to checkson",
		Long:   loginDesc,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

			orgId := args[0]
			resp, err := http.Post("http://127.0.0.1:8080/token/"+orgId, "", nil)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			output.PrintStrings("Response status:", resp.Status)

			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				panic(readErr)
			}

			authToken := string(body[:])

			fmt.Println("Received auth token:", authToken)

			configPath := config.GetConfigPath()

			err = os.MkdirAll(configPath, os.ModePerm)
			if err != nil {
				panic(err)
			}

			writeErr := os.WriteFile(auth.GetAuthFile(), []byte(authToken), os.ModePerm)
			if writeErr != nil {
				panic(writeErr)
			}
		},
	}

	//cmdLoginDebug.Flags().StringVarP(&flags.OrgId, "org", "o", "", "The id of the org to login to")

	return cmdLoginDebug
}
