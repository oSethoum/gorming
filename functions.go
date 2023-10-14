package gorming

import (
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func templateFunctions(data *types.TemplateData) template.FuncMap {
	tsNameFunc := func(column types.Column) string {

		if len(column.Tags.Json.Name) > 0 {
			return column.Tags.Json.Name
		}

		if data.Config.Case == types.Camel {
			return utils.Camel(column.Name)
		}

		if data.Config.Case == types.Snake {
			return utils.Snake(column.Name)
		}

		return column.Name
	}

	tsNameStringFunc := func(name string) string {

		if data.Config.Case == types.Camel {
			return utils.Camel(name)
		}

		if data.Config.Case == types.Snake {
			return utils.Snake(name)
		}

		return name
	}

	tableByName := func(name string) types.Table {
		for _, t := range data.Schema.Tables {
			if name == t.Name {
				return t
			}
		}
		panic("Cannot find table with name " + name)
	}

	modelsFunc := func(tables []types.Table) string {
		tablesString := []string{}
		for _, t := range tables {
			tablesString = append(tablesString, t.Name)
		}
		return strings.Join(tablesString, ` | `)
	}

	tsCreateOmitFunc := func(column types.Column) struct {
		Field  string
		Should bool
	} {
		omit := ""
		if column.Edge != nil && column.Edge.TableKey != "ID" {
			omit = tsNameStringFunc(column.Edge.TableKey)
		}
		return struct {
			Field  string
			Should bool
		}{
			Field:  omit,
			Should: omit != "",
		}
	}

	tableNameFunc := func(table types.Table) string {

		if len(table.Table) > 0 {
			return table.Table
		}

		if data.Config.Case == types.Camel {
			return utils.Camel(inflection.Plural(table.Name))
		}

		if data.Config.Case == types.Snake {
			return utils.Snake(inflection.Plural(table.Name))
		}

		return table.Name
	}

	tableNameStringFunc := func(name string) string {
		return tableNameFunc(tableByName(name))
	}

	tsTypeFunc := func(column types.Column) string {
		t := column.RawType
		typesMap := map[string]string{
			"Time":  "string",
			"bool":  "boolean",
			"int":   "number",
			"uint":  "number",
			"float": "number",
		}

		for k, v := range typesMap {
			if strings.HasPrefix(k, column.RawType) {
				t = v
			}
		}

		if strings.Contains(column.Type, "gorm.DeletedAt") {
			t = "string"
		}

		if len(column.Tags.Gorming.Enum) > 0 {
			return strings.Join(column.Tags.Gorming.Enum, `" | "`)
		}

		if strings.Contains(column.Type, "[]") {
			t += "[]"
		}

		return t
	}

	uniqueRelationsFunc := func(table types.Table) string {
		relations := []string{}

		for _, column := range table.Columns {
			if column.Edge != nil && column.Edge.Unique {
				relations = append(relations, tsNameFunc(column))
			}
		}

		return strings.Join(relations, `" | "`)
	}

	tsOptionalCreateFunc := func(column types.Column) string {
		if utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") ||
			len(column.Tags.Gorm.Default) > 0 {
			return "?"
		}
		return ""
	}

	tsOptionalFunc := func(column types.Column) string {
		if utils.In(column.Name, "DeletedAt") || strings.HasPrefix(column.Type, "*") ||
			len(column.Tags.Gorm.Default) > 0 {
			return "?"
		}
		return ""
	}

	return template.FuncMap{
		"plural":           inflection.Plural,
		"models":           modelsFunc,
		"tsName":           tsNameFunc,
		"tsNameString":     tsNameStringFunc,
		"tableName":        tableNameFunc,
		"tableNameString":  tableNameStringFunc,
		"tsType":           tsTypeFunc,
		"uniqueRelations":  uniqueRelationsFunc,
		"tsOptionalCreate": tsOptionalCreateFunc,
		"tsOptional":       tsOptionalFunc,
		"tsCreateOmit":     tsCreateOmitFunc,
	}
}
