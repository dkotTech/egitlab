package comands

import (
	"fmt"
	"sort"

	"github.com/dkotTech/egitlab/internal"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/urfave/cli/v2"
)

func PrintStatusStyles() {
	rows := make([][]string, 0, len(internal.StatusStyles))
	for statusName, style := range internal.StatusStyles {
		rows = append(rows, []string{statusName, style.Emoji, style.Style.Render(statusName)})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(internal.TableBorderStyle).
		BorderRow(true).
		BorderColumn(true).
		Rows(rows...)

	fmt.Println(t.Render())
}

func NewStylesTestCommand() *cli.Command {
	return &cli.Command{
		Name:  "styles-test",
		Usage: "print all styles",
		Action: func(cCtx *cli.Context) error {
			PrintStatusStyles()

			return nil
		},
	}
}
