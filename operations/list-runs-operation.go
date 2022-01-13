package operations

import (
	"encoding/json"
	"fmt"
	"github.com/huditech/checkson/operations/auth"
	"github.com/huditech/checkson/output"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ListRunsOperation struct {
}

type ListRunsOperationFlags struct {
}

type Run struct {
	Id                string `json:"id"`
	CheckName         string `json:"checkName"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	Success           bool   `json:"success"`
	DurationInSeconds int    `json:"durationInSeconds"`
}

func (operation *ListRunsOperation) ListRunsOperation(flags ListRunsOperationFlags) error {

	// TODO: Call the to-be-created runs API

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/api/runs", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.Debugf("Response status: %s", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	// TODO: Make this reusable and use it for other operations as well
	if resp.StatusCode == 401 {
		return errors.New("Not logged in. Please log in using 'checkson login'")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response from API. Status code: %d, Status: %s", resp.StatusCode, resp.Status)
	}

	var runs []Run
	jsonErr := json.Unmarshal(body, &runs)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
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
