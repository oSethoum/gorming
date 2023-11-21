package gorming

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func templateFunctions(data *types.TemplateData) template.FuncMap {
	tsNameFunc := func(column types.Column) string {
		if data.Config.Case == types.Camel {
			return utils.Camel(column.Name)
		}

		if data.Config.Case == types.Snake {
			return utils.Snake(column.Name)
		}

		if len(column.Tags.Json.Name) > 0 {
			return column.Tags.Json.Name
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

	dartCreateOptionalFunc := func(column types.Column) bool {
		return utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") || column.Edge != nil || strings.HasSuffix(column.Name, "ID") ||
			len(column.Tags.Gorm.Default) > 0

	}

	dartCreateOptionalStringFunc := func(column types.Column) string {
		if dartCreateOptionalFunc(column) {
			return "?"
		}
		return ""
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
		if column.Tags.Gorming.TsType != "" {
			return column.Tags.Gorming.TsType
		}

		t := column.RawType
		typesMap := map[string]string{
			"Time":   "string",
			"bool":   "boolean",
			"int":    "number",
			"uint":   "number",
			"float":  "number",
			"string": "string",
		}

		for k, v := range typesMap {
			if strings.HasPrefix(t, k) {
				t = v
			}
		}

		if strings.Contains(column.Type, "gorm.DeletedAt") {
			t = "string"
		}

		if len(column.Tags.Gorming.Enum) > 0 {
			return `"` + strings.Join(column.Tags.Gorming.Enum, `" | "`) + `"`
		}

		if strings.Contains(column.Type, "[]") {
			t += "[]"
		}

		if t == column.RawType {
			t = "any"
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
			strings.HasPrefix(column.Type, "*") || strings.HasSuffix(column.Name, "ID") || column.Edge != nil ||
			len(column.Tags.Gorm.Default) > 0 {
			return "?"
		}
		return ""
	}

	tsOptionalFunc := func(column types.Column) string {
		if utils.In(column.Name, "DeletedAt") || strings.HasPrefix(column.Type, "*") || column.Slice ||
			len(column.Tags.Gorm.Default) > 0 {
			return "?"
		}
		return ""
	}

	columnOptionalFunc := func(column types.Column) bool {
		return utils.In(column.Name, "DeletedAt") || strings.HasPrefix(column.Type, "*") ||
			len(column.Tags.Gorm.Default) > 0 || column.Slice || column.Edge != nil
	}

	columnOptionalCreateFunc := func(column types.Column) bool {
		return utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") ||
			len(column.Tags.Gorm.Default) > 0 || column.Edge != nil
	}

	dartTypeFunc := func(column types.Column, mode ...string) string {

		t := column.RawType
		typesMap := map[string]string{
			"Time":   "DateTime",
			"bool":   "bool",
			"int":    "int",
			"uint":   "int",
			"float":  "double",
			"string": "String",
		}

		for k, v := range typesMap {
			if strings.HasPrefix(t, k) {
				t = v
			}
		}

		if strings.Contains(column.Type, "gorm.DeletedAt") {
			t = "String"
		}

		if column.Edge != nil && len(mode) > 0 && mode[0] == "create" {
			t += "CreateInput"
		}

		if column.Edge != nil && len(mode) > 0 && mode[0] == "update" {
			t += "UpdateInput"
		}

		if strings.Contains(column.Type, "[]") {
			t = fmt.Sprintf("List<%s>", t)
		}

		return t

	}

	dartNameFunc := func(column types.Column) string {
		return utils.Camel(column.Name)
	}

	tsOptionalKeyFunc := func(column types.Column) string {
		if !(utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") || column.Edge != nil ||
			len(column.Tags.Gorm.Default) > 0) && strings.HasSuffix(column.Name, "ID") {
			return ` /** required edge */`
		}
		return ""
	}

	return template.FuncMap{
		"plural":                   inflection.Plural,
		"models":                   modelsFunc,
		"tsName":                   tsNameFunc,
		"tsNameString":             tsNameStringFunc,
		"tableName":                tableNameFunc,
		"tableNameString":          tableNameStringFunc,
		"tsType":                   tsTypeFunc,
		"columnOptional":           columnOptionalFunc,
		"columnOptionalCreate":     columnOptionalCreateFunc,
		"uniqueRelations":          uniqueRelationsFunc,
		"tsOptionalCreate":         tsOptionalCreateFunc,
		"tsOptionalKey":            tsOptionalKeyFunc,
		"tsOptional":               tsOptionalFunc,
		"tsCreateOmit":             tsCreateOmitFunc,
		"dartType":                 dartTypeFunc,
		"dartName":                 dartNameFunc,
		"dartOptionalCreate":       dartCreateOptionalFunc,
		"dartOptionalCreateString": dartCreateOptionalStringFunc,
	}
}
