package gorming

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func parseTemplate(templateName string, data types.TemplateData) *bytes.Buffer {

	file, err := templates.ReadFile("templates/" + templateName + ".tmpl")
	if err != nil {
		log.Fatalf("gorming: error parsing templates %s, %s \n", templateName, err.Error())
	}
	buffer := new(bytes.Buffer)
	engine, _ := template.New(templateName).Funcs(templateFunctions(&data)).Parse(string(file))
	err = engine.Execute(buffer, data)
	if err != nil {
		log.Fatalf("gorming: error executing template %s, %s \n", templateName, err.Error())
	}
	return buffer
}

func writeTemplate(templateName string, filepath string, data types.TemplateData, file types.File) {
	if !utils.In(file, data.Config.Fiels...) == bool(data.Config.FilesAction) {
		return
	}
	buffer := parseTemplate(templateName, data)
	writeFile(path.Join(data.Config.Paths.BasePath, filepath), buffer.Bytes())
}

func writeFile(filepath string, data []byte) {
	err := os.MkdirAll(path.Dir(filepath), 0777)

	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}

	err = os.WriteFile(filepath, data, 07777)
	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}
}

func defaultConfig(config types.Config) types.Config {
	basePath, pkg := utils.CurrentGoMod()
	config.Package = utils.Choice(config.Package, pkg)
	config.Case = utils.Choice(config.Case, types.Snake)
	config.DBKind = utils.Choice(config.DBKind, types.SQLite)
	config.Paths.TypescriptClient = utils.Choice(config.Paths.TypescriptClient, "client/gorming/")
	config.Paths.DartClient = utils.Choice(config.Paths.DartClient, "client/dart/")
	config.Paths.BasePath = utils.Choice(config.Paths.BasePath, basePath)
	return config
}

func writeJSON(filename string, data any) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
	writeFile(filename, buffer.Bytes())
}
