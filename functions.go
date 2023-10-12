package gorming

import (
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
	"github.com/oSethoum/gorming/types"
)

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"plural": inflection.Plural,
		"models": func(tables []types.Table) string {
			tablesString := []string{}
			for _, t := range tables {
				tablesString = append(tablesString, t.Name)
			}
			return strings.Join(tablesString, `" | "`)
		},
		"tsName": func(column types.Column) string {
			return ""
		},
		"tsType": func(column types.Column) string {
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

			if strings.Contains(column.Type, "[]") {
				t += "[]"
			}

			return t
		},
	}
}
