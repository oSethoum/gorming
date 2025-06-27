package gorming

import (
	"embed"
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
		if config.DBKind == types.MySQL {
			writeTemplate("common/m_query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.FileQuery)
		} else {
			writeTemplate("common/query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.FileQuery)
		}
		writeTemplate("common/schema", filepath.Join(config.Paths.BackendPath, "db/schema.go"), data, types.FileSchema)
		writeTemplate("common/hooks", filepath.Join(config.Paths.BackendPath, "db/hooks.go"), data, types.FileHooks)
		writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.FileUtils)
		writeTemplate("common/error", filepath.Join(config.Paths.BackendPath, "handlers/error.go"), data, types.FileError)
		writeTemplate("server/handler", filepath.Join(config.Paths.BackendPath, "handlers/handler.go"), data, types.FileHandler)
		writeTemplate("server/response", filepath.Join(config.Paths.BackendPath, "handlers/response.go"), data, types.FileResponse)
		writeTemplate("server/ws", filepath.Join(config.Paths.BackendPath, "handlers/ws.go"), data, types.FileWs)
		writeTemplate("server/routes", filepath.Join(config.Paths.BackendPath, "routes/routes.go"), data, types.FileRoutes)

		for _, v := range data.Config.Paths.TypescriptClient {
			writeTemplate("client/api", filepath.Join(v, "api.ts"), data, types.FileTsApi)
			writeTemplate("client/types", filepath.Join(v, "types.ts"), data, types.FileTsTypes)
			writeTemplate("client/event", filepath.Join(v, "event.ts"), data, types.FileTsEvent)
			writeTemplate("client/request", filepath.Join(v, "request.ts"), data, types.FileRequest)
		}
	}
}
