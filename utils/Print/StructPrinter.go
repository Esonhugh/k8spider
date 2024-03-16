package Print

import (
	"fmt"
	"reflect"
)

// PrintStruct prints a table of the given struct
// even print the tags when your struct has. like json or others.
func PrintStruct(in interface{}, tags ...string) {
	if in == nil {
		return
	}
	table := &Table{
		Header: append([]string{"Key", "Value"}, tags...),
		Body:   getBody(in, tags...),
	}
	table.Print("")
}

// getBody returns a slice of strings representing the struct fields
// building the Table body.
func getBody(v interface{}, tags ...string) [][]string {
	t := reflect.TypeOf(v).Elem()
	r := reflect.ValueOf(v)
	rows := make([][]string, 0)
	for i := 0; i < t.NumField(); i++ {
		row := make([]string, 0)
		field := t.Field(i)
		// unexported key is hide.
		if !field.IsExported() {
			continue
		}
		row = append(row,
			field.Name,
			fmt.Sprint(reflect.Indirect(r).FieldByName(field.Name)),
		)
		// for each tag get the field value.
		for _, v := range tags {
			column := field.Tag.Get(v)
			row = append(row, column)
		}
		rows = append(rows, row)
	}
	return rows
}
