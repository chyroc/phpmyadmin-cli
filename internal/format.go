package internal

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func Format(name []string, values ...[]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(name)

	for _, v := range values {
		table.Append(v)
	}
	table.Render()
}

func FormatList(title string, values []string) {
	var vs [][]string
	for _, v := range values {
		vs = append(vs, []string{v})
	}
	Format([]string{title}, vs...)
}
