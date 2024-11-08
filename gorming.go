package gorming

import (
	"embed"
	"path"
	"path/filepath"

	"github.com/oSethoum/gorming/parser"
	"github.com/oSethoum/gorming/types"
)

//go:embed templates
var templates embed.FS

func New(config types.Config) types.Engine {
	config = defaultConfig(config)
	return func(tables []any, Types ...any) {
		schema := parser.Parse(tables, Types...)

		if config.Debug {
			writeJSON("schema.json", schema)
		}

		data := types.TemplateData{
			Schema: schema,
			Config: config,
		}

		writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.FileDB)
		writeTemplate("common/migration", filepath.Join(config.Paths.BackendPath, "db/migration.go"), data, types.FileMigration)
		writeTemplate("common/query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.FileQuery)
		writeTemplate("common/schema", filepath.Join(config.Paths.BackendPath, "db/schema.go"), data, types.FileSchema)
		writeTemplate("common/privacy", filepath.Join(config.Paths.BackendPath, "db/privacy.go"), data, types.FileSchema)
		writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.FileUtils)
		writeTemplate("common/error", filepath.Join(config.Paths.BackendPath, "handlers/error.go"), data, types.FileError)
		writeTemplate("server/handler", filepath.Join(config.Paths.BackendPath, "handlers/handler.go"), data, types.FileHandler)
		writeTemplate("server/response", filepath.Join(config.Paths.BackendPath, "handlers/response.go"), data, types.FileResponse)
		writeTemplate("server/ws", filepath.Join(config.Paths.BackendPath, "handlers/ws.go"), data, types.FileWs)
		writeTemplate("server/routes", filepath.Join(config.Paths.BackendPath, "routes/routes.go"), data, types.FileRoutes)
		writeTemplate("client/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.FileTsApi)
		writeTemplate("client/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.FileTsTypes)

	}
}
