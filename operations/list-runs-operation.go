package operations

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
	"os"
	"strconv"
)

type ListRunsOperation struct {
}

type ListRunsOperationFlags struct {
	DevMode bool
}

func (operation *ListRunsOperation) ListRunsOperation(flags ListRunsOperationFlags) error {

	// TODO: Call the to-be-created runs API

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	runs, err1 := services.ListRuns(authToken, flags.DevMode)
	if err1 != nil {
		return err1
	}

	var data [][]string
	for _, run := range runs {
		successString := "Success"
		if !run.Success {
			successString = "Failure"
		}
		data = append(data, []string{run.CheckName, run.Id, run.StartTime, run.EndTime, strconv.Itoa(run.DurationInSeconds), successString})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Check", "Id", "Start Time", "End Time", "Duration (s)", "Success"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.Render()

	return nil
}
