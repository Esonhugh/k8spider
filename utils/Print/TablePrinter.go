package Print

import (
	"encoding/csv"
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// Table is Struct for printable table
type Table struct {
	// Table Header as []string{"id","name","value" .....}
	Header []string
	// Table Body is content of table.
	Body [][]string

	// They will false in default.

	ForceWriteLogFile   bool
	ForceNoWriteLogFile bool // ForceNoWriteLogFile will not write log file and ignore all other settings include ForceWriteLogFile.

	ShrinkBigTableInPrint bool                      // ShrinkBigTable will shrink the table if it's too large.
	ShrinkBigTableFunc    func(table *Table) *Table // ShrinkTableFunc is a function to shrink the table.

	*tablewriter.Table
}

func NewTable() (ret *Table) {
	ret = &Table{}
	ret.GenTableWriter()
	return
}

func (t *Table) GenTableWriter() {
	if t.Table != nil { // if t.Table is not nil, don't do it again.
		return
	}
	newTable := tablewriter.NewWriter(os.Stdout)
	newTable.SetAutoMergeCells(true)
	newTable.SetRowLine(true)
	newTable.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	newTable.SetAlignment(tablewriter.ALIGN_CENTER)
	t.Table = newTable
}

// Print function prints the table itself.
func (t *Table) Print(Caption string) {
	// If Table is too large. Hacker can use export CFP_TABLE_LOG=true to print it to file cfp_debug.log
	if t.ForceWriteLogFile || log.GetLevel() >= log.DebugLevel || os.Getenv("CFP_TABLE_LOG") == "true" {
		if !t.ForceNoWriteLogFile { // If true bypass no log
			t.writeToLogFile(Caption)
		}
	}
	if t.ShrinkBigTableInPrint {
		if t.ShrinkBigTableFunc != nil {
			t = t.ShrinkBigTableFunc(t)
		}
	}
	t.GenTableWriter() // if table is not init, init as default.
	table := t.Table
	table.SetHeader(t.Header)
	var TableHeaderColor = make([]tablewriter.Colors, len(t.Header))
	for i := range TableHeaderColor {
		TableHeaderColor[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
	}
	table.SetHeaderColor(TableHeaderColor...)
	if Caption != "" {
		table.SetCaption(true, Caption)
	}
	table.ClearRows()
	table.AppendBulk(t.Body)
	table.Render()
}

func (t *Table) writeToLogFile(cap string) {
	file, err := os.OpenFile("cfp_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString("\n# Begin of " + cap + " Table Print CSV Format (csv Comma == '$')\n")
	csvFile := csv.NewWriter(file)
	csvFile.Comma = rune('$')
	csvFile.Write(t.Header)
	csvFile.WriteAll(t.Body)
	csvFile.Flush()
	file.WriteString("\n# End of " + cap + " Table Print CSV Format\n\n")
	// Table Header
	/*
		file.WriteString("|")
		for _, v := range t.Header {
			file.WriteString(" " + v + " |")
		}
		// Header and Body Separator
		file.WriteString("\n|")
		for i := 0; i < len(t.Header); i++ {
			file.WriteString("---|")
		}
		// Body
		file.WriteString("\n")
		for _, v := range t.Body {
			file.WriteString("|")
			for _, vv := range v {
				linewrapper := strings.ReplaceAll(vv, "\n", "<br>")
				file.WriteString(" " + linewrapper + " |")
			}
			file.WriteString("\n")
		}
		file.WriteString("\n\n")
	*/
	return
}
