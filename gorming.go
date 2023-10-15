package gorming

import (
	"embed"
	"path"

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

		writeTemplate("common/db", "db/db.go", data, types.DB)
		writeTemplate("common/migrate", "db/migrate.go", data, types.Migrate)
		writeTemplate("common/utils", "utils/utils.go", data, types.Utils)
		writeTemplate("common/model", "models/model.go", data, types.Model)

		writeTemplate("common/query", "handlers/query.go", data, types.Query)
		writeTemplate("server/fiber/error", "handlers/error.go", data, types.Error)
		writeTemplate("server/fiber/handler", "handlers/handler.go", data, types.Handler)
		writeTemplate("server/fiber/images", "handlers/images.go", data, types.Images)
		writeTemplate("server/fiber/response", "handlers/response.go", data, types.Response)
		writeTemplate("server/fiber/ws", "handlers/ws.go", data, types.Ws)
		writeTemplate("server/fiber/routes", "routes/routes.go", data, types.Routes)

		writeTemplate("client/typescript/api", path.Join(data.Config.Paths.TypescriptClient, "api.ts"), data, types.TsApi)
		writeTemplate("client/typescript/types", path.Join(data.Config.Paths.TypescriptClient, "types.ts"), data, types.TsTypes)

		writeTemplate("client/dart/api", path.Join(data.Config.Paths.DartClient, "api.dart"), data, types.DartApi)
		writeTemplate("client/dart/types", path.Join(data.Config.Paths.DartClient, "types.dart"), data, types.DartTypes)
	}
}
