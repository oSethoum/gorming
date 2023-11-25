package gorming

import (
	"embed"
	"path"
	"path/filepath"

	"github.com/oSethoum/gorming/parser"
	"github.com/oSethoum/gorming/swagger"
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

		if config.Server == types.Fiber {
			writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.FileDB)
			writeTemplate("common/migration", filepath.Join(config.Paths.BackendPath, "db/migration.go"), data, types.FileMigration)
			writeTemplate("common/query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.FileQuery)
			writeTemplate("common/schema", filepath.Join(config.Paths.BackendPath, "db/schema.go"), data, types.FileSchema)
			writeTemplate("common/privacy", filepath.Join(config.Paths.BackendPath, "db/privacy.go"), data, types.FileSchema)
			writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.FileUtils)
			writeTemplate("common/error", filepath.Join(config.Paths.BackendPath, "handlers/error.go"), data, types.FileError)

			writeTemplate("server/fiber/handler", filepath.Join(config.Paths.BackendPath, "handlers/handler.go"), data, types.FileHandler)
			writeTemplate("server/fiber/response", filepath.Join(config.Paths.BackendPath, "handlers/response.go"), data, types.FileResponse)
			writeTemplate("server/fiber/ws", filepath.Join(config.Paths.BackendPath, "handlers/ws.go"), data, types.FileWs)
			writeTemplate("server/fiber/routes", filepath.Join(config.Paths.BackendPath, "routes/routes.go"), data, types.FileRoutes)

			writeTemplate("client/typescript/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.FileTsApi)
			writeTemplate("client/typescript/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.FileTsTypes)

			if config.WithDartClient {
				writeTemplate("client/dart/api", path.Join(data.Config.Paths.DartClient, "api.dart"), data, types.FileDartApi)
				writeTemplate("client/dart/types", path.Join(data.Config.Paths.DartClient, "types.dart"), data, types.FileDartTypes)
			}

			if config.WithSwagger {
				writeJSON(path.Join(config.SwaggerConfig.Output, config.SwaggerConfig.FileName+".json"), swagger.Generate(&config, schema.Tables, schema.Types...))
			}

		}

		if config.Server == types.Wails {
			writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.FileDB)
			writeTemplate("server/wails/migrate", filepath.Join(config.Paths.BackendPath, "db/migrate.go"), data, types.FileMigration)
			writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.FileUtils)

			writeTemplate("server/wails/api", "api.go", data, types.FileApi)
			writeTemplate("server/wails/error", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "error.go"), data, types.FileError)
			writeTemplate("server/wails/query", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "query.go"), data, types.FileQuery)
			writeTemplate("server/wails/resource", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "resource.go"), data, types.FileResource)
			writeTemplate("server/wails/response", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "response.go"), data, types.FileResponse)
			writeTemplate("server/wails/subscription", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "subscription.go"), data, types.FileResponse)

			writeTemplate("client/wails/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.FileTsApi)
			writeTemplate("client/wails/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.FileTsTypes)
		}
	}
}
