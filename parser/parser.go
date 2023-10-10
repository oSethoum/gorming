package parser

import (
	"reflect"
	"strings"

	"github.com/oSethoum/gorming/types"
)

type TypeMap map[string]reflect.Type

func Tables(tablesMap, typesMap *TypeMap) []types.Table {
	tables := []types.Table{}

	for name, table := range *tablesMap {
		tables = append(tables, types.Table{
			Name:    name,
			Table:   table.Name(),
			Columns: Columns(tablesMap, typesMap, table),
		})
	}

	return tables
}

func Columns(tablesMap, typesMap *TypeMap, table reflect.Type) []types.Column {
	columns := []types.Column{}

	for i := 0; i < table.NumField(); i++ {
		f := table.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			columns = append(columns, Columns(tablesMap, typesMap, f.Type)...)
			continue
		}

		slices := strings.Split(f.Type.String(), ".")
		rawType := slices[len(slices)-1]
		rawType = cleanString(rawType, "[]", "*")

		column := types.Column{
			Name:    f.Name,
			Type:    f.Type.String(),
			RawType: rawType,
		}

		if table, ok := (*tablesMap)[rawType]; ok {
			column.Edge = &types.Edge{
				Table: table.Name(),
			}
		}
		columns = append(columns, column)
	}

	return columns
}

func Parse(_tables []any, _types ...any) *types.Schema {

	tablesMap := TypeMap{}
	typesMap := TypeMap{}

	for _, v := range _tables {
		t := reflect.TypeOf(v)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		tablesMap[t.Name()] = t
	}

	for _, v := range _types {
		t := reflect.TypeOf(v)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		typesMap[t.Name()] = t
	}

	return &types.Schema{
		Tables: Tables(&tablesMap, &typesMap),
	}
}
