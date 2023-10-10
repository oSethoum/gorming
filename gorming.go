package gorming

import (
	"bytes"
	"encoding/json"
	"os"
	"path"

	"github.com/oSethoum/gorming/parser"
	"github.com/oSethoum/gorming/types"
)

type Engine = func(tables []any, types ...any)

func New(config types.Config) Engine {
	return func(tables []any, types ...any) {
		schema := parser.Parse(tables, types...)

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.Encode(schema)

		cwd, _ := os.Getwd()
		os.WriteFile(path.Join(cwd, "schema.json"), buffer.Bytes(), 07777)
	}
}
