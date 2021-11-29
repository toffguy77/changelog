package templates

import (
	"fmt"
	"strconv"

	"github.com/alexeyco/simpletable"
)

func FillTable(diff [][]string) {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Ticker ID"},
			{Align: simpletable.AlignCenter, Text: "Commit"},
		},
	}

	for i, row := range diff {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(i)},
			{Align: simpletable.AlignRight, Text: string(row[0])},
			{Align: simpletable.AlignRight, Text: string(row[1])},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleMarkdown)
	fmt.Println(table.String())
}
