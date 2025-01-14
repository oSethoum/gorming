package parser

import (
	"log"
	"reflect"
	"strings"

	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func Tables(tablesMap *types.TypeMap, typesMode bool) []types.Table {
	tables := []types.Table{}

	for name, table := range *tablesMap {
		_, ok := table.Type().MethodByName("Table")
		method := table.MethodByName("Table")

		newTable := types.Table{
			Name:    name,
			Columns: Columns(tablesMap, table.Type(), typesMode),
		}
		if ok {
			newTable.Table = method.Call(nil)[0].String()
		}

		tables = append(tables, newTable)
	}

	return tables
}

func Columns(tablesMap *types.TypeMap, table reflect.Type, typesMode bool) []types.Column {

	fieldsMap := &types.FieldMap{}
	fields(fieldsMap, table)
	columns := []types.Column{}

	for name, f := range *fieldsMap {
		slices := strings.Split(f.Type.String(), ".")
		rawType := slices[len(slices)-1]
		rawType = utils.CleanString(rawType, "[]", "*")

		column := types.Column{
			Name:    name,
			Type:    f.Type.String(),
			RawType: rawType,
			Tags:    tags(f),
			Slice:   strings.Contains(f.Type.String(), "[]"),
		}

		if edgeTable, ok := (*tablesMap)[column.RawType]; ok && !typesMode && !column.Tags.Typescript.SkipEdge {
			edge := &types.Edge{
				Table:  column.RawType,
				Unique: !strings.Contains(column.Type, "[]"),
			}

			if column.Tags.Gorm.Many2Many != "" {
				edge.Many2Many = column.Tags.Gorm.Many2Many
				edge.LocalKey = "ID"
				edge.TableKey = "ID"
			} else {
				edgeTableFieldsMap := &types.FieldMap{}
				fields(edgeTableFieldsMap, edgeTable.Type())

				var keyFound, referenceFound bool

				if edge.Unique {
					key := utils.Choice(column.Tags.Gorm.ForeignKey, table.Name()+"ID")

					if _, keyFound = (*edgeTableFieldsMap)[key]; keyFound {
						edge.TableKey = key
					}

					if !keyFound {
						key = utils.Choice(column.Tags.Gorm.ForeignKey, column.Name+"ID")
					}

					var keyLocal bool
					if !keyFound {
						if _, keyFound = (*fieldsMap)[key]; keyFound {
							edge.LocalKey = key
							keyLocal = true
						}
					}

					reference := utils.Choice(column.Tags.Gorm.References, "ID")
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
					key := utils.Choice(column.Tags.Gorm.ForeignKey, table.Name()+"ID")
					reference := utils.Choice(column.Tags.Gorm.References, "ID")

					if _, keyFound = (*edgeTableFieldsMap)[key]; keyFound {
						edge.TableKey = key
					}

					if _, referenceFound = (*fieldsMap)[reference]; referenceFound {
						edge.LocalKey = reference
					}

				}

				if !keyFound {
					log.Fatalf("gorming: cannot find foreignKey for %s.%s", table.Name(), column.Name)
				}

				if !referenceFound {
					log.Fatalf("gorming: cannot find reference for %s.%s", table.Name(), column.Name)
				}
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
		t := reflect.ValueOf(v)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		tablesMap[t.Type().Name()] = t
	}

	for _, v := range typesArray {
		if v == nil {
			continue
		}
		t := reflect.ValueOf(v)

		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		typesMap[t.Type().Name()] = t
	}

	return &types.Schema{
		Tables: Tables(&tablesMap, false),
		Types:  Tables(&typesMap, true),
	}
}
