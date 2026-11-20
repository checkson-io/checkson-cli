package operations

import (
	"fmt"
	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/output"
	"github.com/checkson-io/checkson-cli/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strconv"
	"strings"
)

type ListRunsOperation struct {
}

type ListRunsOperationFlags struct {
	DevMode bool
}

func (operation *ListRunsOperation) ListRunsOperation(checkName string, flags ListRunsOperationFlags) error {

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	runs, err1 := services.ListRuns(authToken, checkName, flags.DevMode)
	if err1 != nil {
		return err1
	}

	var data [][]string
	for _, run := range runs {
		caser := cases.Title(language.AmericanEnglish)
		successString := caser.String(strings.ToLower(run.Outcome))
		data = append(data, []string{run.CheckName, run.Id, run.StartTime, run.EndTime, strconv.Itoa(run.DurationInSeconds), successString})
	}

	return output.RenderTable(os.Stdout, []string{"CHECK", "ID", "START TIME", "END TIME", "DURATION (s)", "SUCCESS"}, data)
}
