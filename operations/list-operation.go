package operations

import (
	"errors"
	"github.com/olekukonko/tablewriter"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
	"os"
	"strconv"
)

type ListOperation struct {
}

type ListOperationFlags struct {
	DevMode bool
}

func (operation *ListOperation) ListOperation(flags ListOperationFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	checks, err1 := services.ListChecks(authToken, flags.DevMode)
	if err1 != nil {
		return err1
	}

	var data [][]string
	for _, check := range checks {
		data = append(data, []string{check.Name, check.Status, check.LastStatusChange, check.DockerImage, strconv.Itoa(int(check.CheckIntervalInMinutes))})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Since", "Docker Image", "Interval"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.Render()

	return nil
}
