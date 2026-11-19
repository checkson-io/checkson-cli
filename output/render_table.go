package output

import (
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// RenderTable writes the CLI's borderless tables with a separator below the header.
func RenderTable(w io.Writer, headers []string, rows [][]string) error {
	table := tablewriter.NewTable(w, tablewriter.WithHeaderAutoFormat(tw.Off), tablewriter.WithRendition(tw.Rendition{
		Borders:  tw.BorderNone,
		Symbols:  tw.NewSymbols(tw.StyleASCII),
		Settings: tw.Settings{Lines: tw.Lines{ShowHeaderLine: tw.On}},
	}))
	table.Header(headers)
	for _, row := range rows {
		if err := table.Append(row); err != nil {
			return err
		}
	}
	return table.Render()
}
