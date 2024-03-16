package Print

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultShrinkTableFunc is the default function to shrink the table.
// It will shrink the table to 7 columns and width of each unit is less than 50 characters.
func DefaultShrinkTableFunc(table *Table) *Table {
	if len(table.Header) > 7 {
		newTable := NewTable()

		tableWidth := len(table.Header)
		var leftTable, rightTable []string
		if tableWidth%2 == 0 {
			leftTable = append([]string{"UID"}, table.Header[:tableWidth/2]...)
			rightTable = append([]string{"UID"}, table.Header[tableWidth/2:]...)
		} else {
			half := tableWidth/2 + 1
			leftTable = table.Header[:half]
			rightTable = append([]string{table.Header[0]}, table.Header[half:]...)
		}
		for k := 0; k < len(leftTable); k++ {
			newTable.Header = append(newTable.Header, fmt.Sprintf("%v/%v", leftTable[k], rightTable[k]))
		}

		for rowId, row := range table.Body {
			if tableWidth%2 == 0 {
				row1 := []string{
					fmt.Sprintf("%v", rowId),
				}
				row1 = append(row1, row[:tableWidth/2]...)
				row2 := []string{
					fmt.Sprintf("%v", rowId),
				}
				row2 = append(row2, row[tableWidth/2:]...)
				newTable.Body = append(newTable.Body, row1, row2)
			} else {
				half := tableWidth/2 + 1
				row1 := []string{}
				row1 = append(row1, row[:half]...)
				row2 := []string{
					row[0],
				}
				row2 = append(row2, row[half:]...)
				newTable.Body = append(newTable.Body, row1, row2)
			}
		}
		table = newTable
		table.SetAutoMergeCells(true)
		table.SetAutoMergeCellsByColumnIndex([]int{0})
	}

	for rowID, row := range table.Body {
		for unitID, unit := range row {
			lines := strings.Split(unit, "\n")
			for i, line := range lines {
				if len(line) > 50 {
					lines[i] = line[:47] + "..."
				}
			}
			table.Body[rowID][unitID] = strings.Join(lines, "\n")
		}
	}
	return table
}

func (t *Table) OrderByColumnIndex(index int) {
	tempBody := t.Body
	sort.Slice(tempBody, func(i, j int) bool {
		return t.Body[i][index] < t.Body[j][index]
	})
}
