package parser

import (
	"log"
	"reflect"
	"strings"

	"github.com/oSethoum/gorming/types"
)

func Tables(tablesMap, typesMap *types.TypeMap) []types.Table {
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

func Columns(tablesMap, typesMap *types.TypeMap, table reflect.Type) []types.Column {

	fieldsMap := &types.FieldMap{}
	fields(fieldsMap, table)
	columns := []types.Column{}

	for name, f := range *fieldsMap {
		slices := strings.Split(f.Type.String(), ".")
		rawType := slices[len(slices)-1]
		rawType = cleanString(rawType, "[]", "*")

		column := types.Column{
			Name:    name,
			Type:    f.Type.String(),
			RawType: rawType,
			Tags:    tags(f),
			Slice:   strings.Contains(f.Type.String(), "[]"),
		}

		if edgeTable, ok := (*tablesMap)[column.RawType]; ok {

			edgeTableFieldsMap := &types.FieldMap{}
			fields(edgeTableFieldsMap, edgeTable)

			edge := &types.Edge{
				Table:  column.RawType,
				Unique: !strings.Contains(column.Type, "[]"),
			}

			var keyFound, referenceFound bool

			if edge.Unique {
				key := choice(column.Tags.Gorm.ForeignKey, table.Name()+"ID")

				if _, keyFound = (*edgeTableFieldsMap)[key]; keyFound {
					edge.TableKey = key
				}

				if !keyFound {
					key = choice(column.Tags.Gorm.ForeignKey, column.Name+"ID")
				}

				var keyLocal bool
				if !keyFound {
					if _, keyFound = (*fieldsMap)[key]; keyFound {
						edge.LocalKey = key
						keyLocal = true
					}
				}

				reference := choice(column.Tags.Gorm.References, "ID")
				if keyLocal {
					if _, referenceFound = (*edgeTableFieldsMap)[reference]; referenceFound {
						edge.TableKey = reference
					}

				} else {
					if _, referenceFound = (*fieldsMap)[reference]; referenceFound {
						edge.LocalKey = reference
					}
				}

			} else {
				key := choice(column.Tags.Gorm.ForeignKey, table.Name()+"ID")
				reference := choice(column.Tags.Gorm.References, "ID")

				if _, keyFound = (*edgeTableFieldsMap)[key]; keyFound {
					edge.TableKey = key
				}

				if _, referenceFound = (*fieldsMap)[reference]; referenceFound {
					edge.LocalKey = key
				}

			}

			if !keyFound {
				log.Fatalf("gorming: cannot find foreignKey for %s.%s", table.Name(), column.Name)
			}

			if !referenceFound {
				log.Fatalf("gorming: cannot find reference for %s.%s", table.Name(), column.Name)
			}

			column.Edge = edge
		}
		columns = append(columns, column)
	}
	return columns
}

func Parse(tablesArray []any, typesArray ...any) *types.Schema {

	tablesMap := types.TypeMap{}
	typesMap := types.TypeMap{}

	for _, v := range tablesArray {
		t := reflect.TypeOf(v)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		tablesMap[t.Name()] = t
	}

	for _, v := range tablesArray {
		t := reflect.TypeOf(v)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		typesMap[t.Name()] = t
	}

	tables := Tables(&tablesMap, &typesMap)
	return &types.Schema{
		Tables: tables,
	}
}
