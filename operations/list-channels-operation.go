package operations

import (
	"errors"
	"github.com/olekukonko/tablewriter"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
	"os"
)

type ListChannelsOperation struct {
}

type ListChannelsOperationFlags struct {
	DevMode bool
}

func (operation *ListChannelsOperation) ListChannelsOperation(flags ListChannelsOperationFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	channels, err1 := services.ListChannels(authToken, flags.DevMode)
	if err1 != nil {
		return err1
	}

	var data [][]string
	for _, channel := range channels {
		data = append(data, []string{channel.Name, channel.Type})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type"})

	for _, v := range data {
		table.Append(v)
	}
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.Render()

	return nil
}
