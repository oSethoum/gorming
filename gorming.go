package gorming

import (
	"bytes"
	"embed"
	"encoding/json"

	"github.com/oSethoum/gorming/parser"
	"github.com/oSethoum/gorming/types"
)

//go:embed templates
var templates embed.FS

type Engine = func(tables []any, types ...any)

func New(config types.Config) Engine {

	return func(tables []any, Types ...any) {
		schema := parser.Parse(tables, Types...)

		data := types.TemplateData{
			Schema: schema,
			Config: config,
		}

		writeTemplateData("common/query", "query.go", data, types.Query)

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.Encode(schema)

		writeFile("schema.json", buffer.Bytes())
	}
}
