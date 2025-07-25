package gorming

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
	"github.com/samber/lo"
)

func templateFunctions(data *types.TemplateData) template.FuncMap {

	normalizeParamFunc := func(column *types.Column, param string) string {
		if column.RawType == "string" {
			return `"` + param + `"`
		}
		return param
	}

	tableHasValidationFunc := func(table types.Table) bool {
		for _, c := range table.Columns {
			if len(c.Tags.Validator) > 0 {
				return true
			}
		}
		return false
	}

	allowValidationFunc := func(s string, rawType ...string) bool {

		for _, t := range data.Schema.Tables {
			for _, c := range t.Columns {
				for _, v := range c.Tags.Validator {
					if v.Rule == s && len(rawType) == 0 {
						return true
					}

					if v.Rule == s {
						for _, rt := range rawType {
							if rt == c.RawType {
								return true
							}
						}
					}

				}
			}
		}
		return false
	}

	hasRegexValidationFunc := func() bool {
		for _, t := range data.Schema.Tables {
			for _, c := range t.Columns {
				for _, v := range c.Tags.Validator {
					if utils.In(v.Rule, "url", "email", "numeric", "alpha", "alphanumeric", "number", "alphaSpace") {
						return true
					}
				}
			}
		}
		return false
	}

	ignoreRouteFunc := func(resource string, method string) bool {

		if data.Config.SkipRoutes != nil {
			if skips, ok := data.Config.SkipRoutes[resource]; ok {
				return strings.Contains(skips, method)
			}
		}

		return false
	}

	ignoreAllRouteFunc := func(resource string) bool {
		if data.Config.SkipRoutes != nil {
			if skips, ok := data.Config.SkipRoutes[resource]; ok {
				return strings.Contains(skips, "all") ||
					(strings.Contains(skips, "query") &&
						strings.Contains(skips, "create") &&
						strings.Contains(skips, "update") &&
						strings.Contains(skips, "delete"))
			}
		}
		return false
	}

	tsNameFunc := func(column types.Column) string {

		if column.Tags.Json.Name != "" {
			return column.Tags.Json.Name
		}

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

	tablePascalFunc := func(table types.Table) string {
		return inflection.Plural(table.Name)
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
		if column.Tags.Typescript.Type != "" {
			return column.Tags.Typescript.Type
		}

		t := column.RawType
		found := false
		for _, v := range data.Schema.Tables {
			if v.Name == t {
				found = true
			}
		}

		for _, v := range data.Schema.Types {
			if v.Name == t {
				found = true
			}
		}

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

		if len(column.Tags.Typescript.Enum) > 0 {
			return `"` + strings.Join(column.Tags.Typescript.Enum, `" | "`) + `"`
		}

		if strings.Contains(column.Type, "[]") {
			t += "[]"
		}

		if t == column.RawType && t != "string" && !found {
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
			strings.HasPrefix(column.Type, "*") || strings.HasSuffix(column.Name, "ID") || column.Edge != nil || column.Tags.Typescript.Optional ||
			len(column.Tags.Gorm.Default) > 0 {
			return "?"
		}
		return ""
	}

	tsNullableCreateFunc := func(column types.Column) string {
		if utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") || strings.HasSuffix(column.Name, "ID") || column.Edge != nil ||
			len(column.Tags.Gorm.Default) > 0 {
			return " | null"
		}
		return ""
	}

	tsOptionalFunc := func(column types.Column) string {
		if utils.In(column.Name, "DeletedAt") || strings.HasPrefix(column.Type, "*") || column.Slice || column.Tags.Typescript.Optional ||
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

	cleanRequiredEdgesFunc := func(columns []types.Column) []types.Column {
		return columns
	}

	_getRequiredEdges := func(columns []types.Column) [][2]types.Column {
		out := [][2]types.Column{}
		for _, v := range columns {
			if v.Edge != nil && v.Edge.Unique && v.Edge.LocalKey != "ID" {
				localKeyColumn, ok := lo.Find(columns, func(c types.Column) bool {
					return c.Edge == nil && c.Name == v.Edge.LocalKey && !strings.HasPrefix(c.Type, "*")
				})
				if ok {
					out = append(out, [2]types.Column{v, localKeyColumn})
				}
			}
		}

		return out
	}

	getTableFKConstraintsFunc := func(table types.Table) string {
		ss := ""
		edges := _getRequiredEdges(table.Columns)

		for _, v := range edges {
			u := ""

			if v[0].Tags.Gorm.OnUpdate != "" {
				u += fmt.Sprintf(`ON UPDATE %s `, v[0].Tags.Gorm.OnUpdate)
			} else if v[1].Tags.Gorm.OnUpdate != "" {
				u += fmt.Sprintf(`ON UPDATE %s `, v[1].Tags.Gorm.OnUpdate)
			}

			if v[0].Tags.Gorm.OnDelete != "" {
				u += fmt.Sprintf(`ON DELETE %s`, v[0].Tags.Gorm.OnDelete)
			} else if v[1].Tags.Gorm.OnDelete != "" {
				u += fmt.Sprintf(`ON DELETE %s`, v[1].Tags.Gorm.OnDelete)
			}

			s := fmt.Sprintf(`DB.Exec("ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(%s) %s")`,
				tableNameStringFunc(table.Name),
				tableNameStringFunc(table.Name),
				tableNameStringFunc(v[0].Edge.Table),
				tsNameStringFunc(v[0].Edge.LocalKey),
				tableNameStringFunc(v[0].Edge.Table),
				tsNameStringFunc(v[0].Edge.TableKey),
				u,
			)

			ss += fmt.Sprintf("%s\n", s)

		}

		return ss
	}

	getTableFKMigratorFunc := func(table types.Table) string {
		ss := ""
		edges := _getRequiredEdges(table.Columns)

		for _, v := range edges {
			u := ""

			if v[0].Tags.Gorm.OnUpdate != "" {
				u += fmt.Sprintf(`ON UPDATE %s `, v[0].Tags.Gorm.OnUpdate)
			} else if v[1].Tags.Gorm.OnUpdate != "" {
				u += fmt.Sprintf(`ON UPDATE %s `, v[1].Tags.Gorm.OnUpdate)
			}

			if v[0].Tags.Gorm.OnDelete != "" {
				u += fmt.Sprintf(`ON DELETE %s`, v[0].Tags.Gorm.OnDelete)
			} else if v[1].Tags.Gorm.OnDelete != "" {
				u += fmt.Sprintf(`ON DELETE %s`, v[1].Tags.Gorm.OnDelete)
			}
			_table := tableNameStringFunc(table.Name)
			_targetTable := tableNameStringFunc(v[0].Edge.Table)
			_constraint := fmt.Sprintf("fk_%s_%s", _table, _targetTable)

			s := fmt.Sprintf(`if DB.Migrator().HasConstraint("%s", "%s") { 
				DB.Migrator().DropConstraint("%s", "%s") 
			}
			`, _table, _constraint, _table, _constraint)

			s += fmt.Sprintf(`DB.Exec("ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(%s) %s")`,
				tableNameStringFunc(table.Name),
				tableNameStringFunc(table.Name),
				tableNameStringFunc(v[0].Edge.Table),
				tsNameStringFunc(v[0].Edge.LocalKey),
				tableNameStringFunc(v[0].Edge.Table),
				tsNameStringFunc(v[0].Edge.TableKey),
				u,
			)

			ss += fmt.Sprintf("%s\n", s)

		}

		return ss
	}

	requiredEdge := func(column types.Column) bool {
		return !(utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") || column.Edge != nil ||
			len(column.Tags.Gorm.Default) > 0) && strings.HasSuffix(column.Name, "ID")
	}

	hasRequiredEdge := func(table types.Table, column types.Column) bool {
		v := false
		column.Name += "ID"
		for _, c := range table.Columns {
			if c.Name == column.Name && requiredEdge(c) {
				v = true
			}
		}
		return v
	}

	tsCreateIgnoreFunc := func(table types.Table, column types.Column) bool {
		return requiredEdge(column) || hasRequiredEdge(table, column)
	}

	tsCreateUnionFunc := func(table types.Table) string {
		edges := _getRequiredEdges(table.Columns)
		buffer := []string{}

		for _, e := range edges {
			buffer = append(buffer, fmt.Sprintf("({ %s: %sCreateInput } | { %s: %s })", tsNameFunc(e[0]), tsTypeFunc(e[0]), tsNameFunc(e[1]), tsTypeFunc(e[1])))
		}
		if len(buffer) > 0 {
			return " & " + strings.Join(buffer, " & ")
		}
		return ""
	}

	tsOptionalKeyFunc := func(column types.Column) string {
		if !(utils.In(column.Name, "ID", "CreatedAt", "UpdatedAt", "DeletedAt") ||
			strings.HasPrefix(column.Type, "*") || column.Edge != nil ||
			len(column.Tags.Gorm.Default) > 0) && strings.HasSuffix(column.Name, "ID") {
			return ` /** required edge */`
		}
		return ""
	}

	setNullFieldTypeFunc := func(table types.Table) string {
		nullableFields := []string{}
		for _, v := range table.Columns {
			if (v.Slice || strings.HasPrefix(v.Type, "*") || v.RawType == "DeletedAt") && v.Name != "SetNULL" {
				nullableFields = append(nullableFields, tsNameFunc(v))
			}
		}

		if len(nullableFields) > 0 {
			return fmt.Sprintf(`Array<"%s">`, strings.Join(nullableFields, `" | "`))
		}

		return "null"
	}

	return template.FuncMap{
		"plural":                inflection.Plural,
		"models":                modelsFunc,
		"tsName":                tsNameFunc,
		"tsNameString":          tsNameStringFunc,
		"tableName":             tableNameFunc,
		"tableNameString":       tableNameStringFunc,
		"tsType":                tsTypeFunc,
		"columnOptional":        columnOptionalFunc,
		"columnOptionalCreate":  columnOptionalCreateFunc,
		"uniqueRelations":       uniqueRelationsFunc,
		"tsOptionalCreate":      tsOptionalCreateFunc,
		"tsNullableCreate":      tsNullableCreateFunc,
		"tsOptionalKey":         tsOptionalKeyFunc,
		"tsOptional":            tsOptionalFunc,
		"tsCreateOmit":          tsCreateOmitFunc,
		"ignoreRoute":           ignoreRouteFunc,
		"ignoreAllRoute":        ignoreAllRouteFunc,
		"allowValidation":       allowValidationFunc,
		"hasRegexValidation":    hasRegexValidationFunc,
		"tableHasValidation":    tableHasValidationFunc,
		"normalizeParam":        normalizeParamFunc,
		"tsCreateIgnore":        tsCreateIgnoreFunc,
		"cleanRequiredEdges":    cleanRequiredEdgesFunc,
		"tsCreateUnion":         tsCreateUnionFunc,
		"tablePascal":           tablePascalFunc,
		"setNullFieldType":      setNullFieldTypeFunc,
		"getTableFKConstraints": getTableFKConstraintsFunc,
		"getTableFKMigrator":    getTableFKMigratorFunc,
	}
}
