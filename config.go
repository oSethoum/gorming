package gorming

import (
	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func defaultConfig(config types.Config) types.Config {
	basePath, pkg := utils.CurrentGoMod()
	config.Package = utils.Choice(config.Package, pkg)
	config.ApiPackage = utils.Choice(config.ApiPackage, "main")
	config.BackendPackage = utils.Choice(config.BackendPackage, "backend")
	config.Case = utils.Choice(config.Case, types.Snake)
	config.DBKind = utils.Choice(config.DBKind, types.SQLite)
	config.Paths.TypescriptClient = utils.Choice(config.Paths.TypescriptClient, "client/typescript/gorming")
	config.Paths.DartClient = utils.Choice(config.Paths.DartClient, "client/dart/gorming")
	config.Paths.BasePath = utils.Choice(config.Paths.BasePath, basePath)
	config.SwaggerConfig.FileName = utils.Choice(config.SwaggerConfig.FileName, "swagger")
	if !utils.In(config.Server, types.Fiber, types.Wails) {
		config.Server = types.Fiber
	}
	return config
}
