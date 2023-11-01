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

		if config.Server == types.Fiber {
			writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.DB)
			writeTemplate("common/migrate", filepath.Join(config.Paths.BackendPath, "db/migrate.go"), data, types.Migrate)
			writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.Utils)
			writeTemplate("common/model", filepath.Join(config.Paths.BackendPath, "models/model.go"), data, types.Model)
			writeTemplate("common/query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.Query)
			writeTemplate("server/fiber/error", filepath.Join(config.Paths.BackendPath, "handlers/error.go"), data, types.Error)
			writeTemplate("server/fiber/handler", filepath.Join(config.Paths.BackendPath, "handlers/handler.go"), data, types.Handler)
			writeTemplate("server/fiber/images", filepath.Join(config.Paths.BackendPath, "handlers/images.go"), data, types.Images)
			writeTemplate("server/fiber/response", filepath.Join(config.Paths.BackendPath, "handlers/response.go"), data, types.Response)
			writeTemplate("server/fiber/ws", filepath.Join(config.Paths.BackendPath, "handlers/ws.go"), data, types.Ws)
			writeTemplate("server/fiber/routes", filepath.Join(config.Paths.BackendPath, "routes/routes.go"), data, types.Routes)
			writeTemplate("client/typescript/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.TsApi)
			writeTemplate("client/typescript/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.TsTypes)
			writeTemplate("client/dart/api", path.Join(data.Config.Paths.DartClient, "api.dart"), data, types.DartApi)
			writeTemplate("client/dart/types", path.Join(data.Config.Paths.DartClient, "types.dart"), data, types.DartTypes)
			if config.WithPrivacy {
				writeTemplate("server/fiber/token", filepath.Join(config.Paths.BackendPath, "handlers/middleware.go"), data, types.TokenMiddleware)
			}
		}

		if config.Server == types.Echo {
			writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.DB)
			writeTemplate("common/migrate", filepath.Join(config.Paths.BackendPath, "db/migrate.go"), data, types.Migrate)
			writeTemplate("common/model", filepath.Join(config.Paths.BackendPath, "models/model.go"), data, types.Model)
			writeTemplate("common/query", filepath.Join(config.Paths.BackendPath, "db/query.go"), data, types.Query)
			writeTemplate("server/echo/error", filepath.Join(config.Paths.BackendPath, "handlers/error.go"), data, types.Error)
			writeTemplate("server/echo/handler", filepath.Join(config.Paths.BackendPath, "handlers/handler.go"), data, types.Handler)
			writeTemplate("server/echo/images", filepath.Join(config.Paths.BackendPath, "handlers/images.go"), data, types.Images)
			writeTemplate("server/echo/response", filepath.Join(config.Paths.BackendPath, "handlers/response.go"), data, types.Response)
			writeTemplate("server/echo/ws", filepath.Join(config.Paths.BackendPath, "handlers/ws.go"), data, types.Ws)
			writeTemplate("server/echo/routes", filepath.Join(config.Paths.BackendPath, "routes/routes.go"), data, types.Routes)
			writeTemplate("client/typescript/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.TsApi)
			writeTemplate("client/typescript/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.TsTypes)
			writeTemplate("client/dart/api", path.Join(data.Config.Paths.DartClient, "api.dart"), data, types.DartApi)
			writeTemplate("client/dart/types", path.Join(data.Config.Paths.DartClient, "types.dart"), data, types.DartTypes)
		}

		if config.Server == types.Wails {
			writeTemplate("common/db", filepath.Join(config.Paths.BackendPath, "db/db.go"), data, types.DB)
			writeTemplate("server/wails/migrate", filepath.Join(config.Paths.BackendPath, "db/migrate.go"), data, types.Migrate)
			writeTemplate("common/utils", filepath.Join(config.Paths.BackendPath, "utils/utils.go"), data, types.Utils)
			writeTemplate("common/model", filepath.Join(config.Paths.BackendPath, "models/model.go"), data, types.Model)

			writeTemplate("server/wails/api", "api.go", data, types.Api)
			writeTemplate("server/wails/error", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "error.go"), data, types.Error)
			writeTemplate("server/wails/query", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "query.go"), data, types.Query)
			writeTemplate("server/wails/resource", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "resource.go"), data, types.Resource)
			writeTemplate("server/wails/response", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "response.go"), data, types.Response)
			writeTemplate("server/wails/subscription", path.Join(config.Paths.BackendPath, data.Config.Paths.ApiPath, "subscription.go"), data, types.Response)

			writeTemplate("client/wails/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.TsApi)
			writeTemplate("client/wails/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.TsTypes)
		}
	}
}
