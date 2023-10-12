package gorming

import (
	"bytes"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/oSethoum/gorming/types"
)

func parseTemplate(templateName string, data any) *bytes.Buffer {

	file, err := templates.ReadFile("templates/" + templateName + ".tmpl")
	if err != nil {
		log.Fatalf("gorming: error parsing templates %s, %s \n", templateName, err.Error())
	}
	buffer := new(bytes.Buffer)
	engine, _ := template.New(templateName).Funcs(templateFunctions()).Parse(string(file))
	err = engine.Execute(buffer, data)
	if err != nil {
		log.Fatalf("gorming: error executing template %s, %s \n", templateName, err.Error())
	}
	return buffer
}

func writeTemplateData(templateName string, filepath string, data types.TemplateData, file types.File) {
	writeTemplate(templateName, filepath, data, in(file, data.Config.Fiels...) == bool(data.Config.FilesAction))
}

func writeTemplate(templateName string, filepath string, data any, write bool) {
	if !write {
		return
	}
	buffer := parseTemplate(templateName, data)
	writeFile(filepath, buffer.Bytes())
}

func writeFile(filepath string, data []byte) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}

	err = os.MkdirAll(path.Dir(path.Join(cwd, filepath)), 0777)

	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}

	err = os.WriteFile(path.Join(cwd, filepath), data, 07777)
	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}
}

func in[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
