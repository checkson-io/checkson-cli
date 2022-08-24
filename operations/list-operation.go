package operations

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"io/ioutil"
	"log"
	"net/http"
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
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	client := &http.Client{}
	req, err1 := http.NewRequest("GET", getApiUrl(flags.DevMode, "api/checks"), nil)

	if err1 != nil {
		return fmt.Errorf("problem preparing request: %w", err1)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()
	output.Debugf("Response status: %s", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var checks []Check
	jsonErr := json.Unmarshal(body, &checks)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	var data [][]string
	for _, check := range checks {
		data = append(data, []string{check.Name, check.LastCheckOutcome, check.LastOutcomeChangeDuration, check.DockerImage, strconv.Itoa(int(check.CheckIntervalInMinutes))})
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
