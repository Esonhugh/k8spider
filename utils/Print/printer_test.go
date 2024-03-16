package Print

import "testing"

func TestTablePrintWithNewlines(t *testing.T) {
	tab := Table{
		Header: []string{"Header1", "Header2"},
		Body: [][]string{
			{"Row1", "Row2"},
		},
	}
	tab.Print("testTablePrintWithNewlines")
	tab.Body = [][]string{
		{"Row1\nRow1", "Row2"},
	}
	tab.Print("testTablePrintWithNewlines")
}

var BigTable = Table{
	Header: []string{
		"Header1", "Header2", "Header3", "Header4", "Header5", "Header6", "Header7", "Header8", "Header9", "Header10",
	},
	Body: [][]string{
		{"Unit1", "Unit2", "Unit3", "Unit4", "Unit5", "Unit6", "Unit7", "Unit8", "Unit9", "Unit10"},
		{"Unit11", "Unit12", "Unit13", "Unit14", "Unit15", "Unit16", "Unit17", "Unit18", "Unit19", "Unit20"},
		{"Unit21", "Unit22", "Unit23", "Unit24", "Unit25", "Unit26", "Unit27", "Unit28", "Unit29", "Unit30"},
		{"Unit31", "Unit32", "Unit33", "Unit34", "Unit35", "Unit36", "Unit37", "Unit38", "Unit39", "Unit40"},
		{"Unit41", "Unit42", "Unit43", "Unit44", "Unit45", "Unit46", "Unit47", "Unit48", "Unit49", "Unit50"},
		{"Unit51", "Unit52", "Unit53", "Unit54", "Unit55", "Unit56", "Unit57", "Unit58", "Unit59", "Unit60"},
		{"Unit61", "Unit62", "Unit63", "Unit64", "Unit65", "Unit66", "Unit67", "Unit68", "Unit69", "Unit70"},
		{"Unit71", "Unit72", "Unit73", "Unit74", "Unit75", "Unit76", "Unit77", "Unit78", "Unit79", "Unit80"},
	},
	ShrinkBigTableInPrint: true,
	ShrinkBigTableFunc:    DefaultShrinkTableFunc,
}

var BigTable2 = Table{
	Header: []string{
		"Header1", "Header2", "Header3", "Header4", "Header5", "Header6", "Header7", "Header8", "Header9",
	},
	Body: [][]string{
		{"Unit1", "Unit2", "Unit3", "Unit4", "Unit5", "Unit6", "Unit7", "Unit8", "Unit9"},
		{"Unit11", "Unit12", "Unit13", "Unit14", "Unit15", "Unit16", "Unit17", "Unit18", "Unit19"},
		{"Unit21", "Unit22", "Unit23", "Unit24", "Unit25", "Unit26", "Unit27", "Unit28", "Unit29"},
		{"Unit31", "Unit32", "Unit33", "Unit34", "Unit35", "Unit36", "Unit37", "Unit38", "Unit39"},
		{"Unit41", "Unit42", "Unit43", "Unit44", "Unit45", "Unit46", "Unit47", "Unit48", "Unit49"},
		{"Unit51", "Unit52", "Unit53", "Unit54", "Unit55", "Unit56", "Unit57", "Unit58", "Unit59"},
		{"Unit61", "Unit62", "Unit63", "Unit64", "Unit65", "Unit66", "Unit67", "Unit68", "Unit69"},
		{"Unit71", "Unit72", "Unit73", "Unit74", "Unit75", "Unit76", "Unit77", "Unit78", "Unit79"},
	},
	ShrinkBigTableInPrint: true,
	ShrinkBigTableFunc:    DefaultShrinkTableFunc,
}

func TestDefaultShrinkTableFunc(t *testing.T) {
	BigTable.Print("TestDefaultShrinkTableFunc")
	BigTable2.Print("")
	BigTable2.Print("")
}

var RandomOrderTable = Table{
	Header: []string{
		"Header1", "Header2", "Header3", "Header4", "Header5", "Header6", "Header7", "Header8", "Header9", "Header10",
	},
	Body: [][]string{
		{"ORDER1", "UNIT 1-2", "UNIT 1-3", "UNIT 1-4", "UNIT 1-5", "UNIT 1-6", "UNIT 1-7", "UNIT 1-8", "UNIT 1-9", "UNIT 1-10"},
		{"ORDER4", "UNIT 7-2", "UNIT 7-3", "UNIT 7-4", "UNIT 7-5", "UNIT 7-6", "UNIT 7-7", "UNIT 7-8", "UNIT 7-9", "UNIT 7-10"},
		{"ORDER2", "UNIT 3-2", "UNIT 3-3", "UNIT 3-4", "UNIT 3-5", "UNIT 3-6", "UNIT 3-7", "UNIT 3-8", "UNIT 3-9", "UNIT 3-10"},
		{"ORDER3", "UNIT 6-2", "UNIT 6-3", "UNIT 6-4", "UNIT 6-5", "UNIT 6-6", "UNIT 6-7", "UNIT 6-8", "UNIT 6-9", "UNIT 6-10"},
		{"ORDER2", "UNIT 4-2", "UNIT 4-3", "UNIT 4-4", "UNIT 4-5", "UNIT 4-6", "UNIT 4-7", "UNIT 4-8", "UNIT 4-9", "UNIT 4-10"},
		{"ORDER1", "UNIT 2-2", "UNIT 2-3", "UNIT 2-4", "UNIT 2-5", "UNIT 2-6", "UNIT 2-7", "UNIT 2-8", "UNIT 2-9", "UNIT 2-10"},
		{"ORDER3", "UNIT 5-2", "UNIT 5-3", "UNIT 5-4", "UNIT 5-5", "UNIT 5-6", "UNIT 5-7", "UNIT 5-8", "UNIT 5-9", "UNIT 5-10"},
		{"ORDER4", "UNIT 8-2", "UNIT 8-3", "UNIT 8-4", "UNIT 8-5", "UNIT 8-6", "UNIT 8-7", "UNIT 8-8", "UNIT 8-9", "UNIT 8-10"},
	},
}

func TestTablePrint(t *testing.T) {
	t.Log("TestTablePrint")
	RandomOrderTable.Print("TestTablePrint After")
	RandomOrderTable.OrderByColumnIndex(0)
	RandomOrderTable.Print("TestTablePrint After")
}
