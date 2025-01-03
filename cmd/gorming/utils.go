package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func parseTemplate(templateName string) *bytes.Buffer {

	file, err := templates.ReadFile("templates/" + templateName + ".tmpl")
	if err != nil {
		log.Fatalf("gorming: error parsing templates %s, %s \n", templateName, err.Error())
	}
	buffer := new(bytes.Buffer)
	engine, err := template.New(templateName).Parse(string(file))
	if err != nil {
		log.Fatalf("gorming: error executing template %s, %s \n", templateName, err.Error())
	}
	err = engine.Execute(buffer, nil)
	if err != nil {
		log.Fatalf("gorming: error executing template %s, %s \n", templateName, err.Error())
	}
	return buffer
}

func writeTemplate(templateName string, filepath string) {
	wd, _ := os.Getwd()
	buffer := parseTemplate(templateName)
	writeFile(path.Join(wd, filepath), buffer.Bytes())
}

func writeFile(outPath string, data []byte) {
	err := os.MkdirAll(filepath.Dir(outPath), 0777)

	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}
	_, err = os.Stat(outPath)
	if err == nil {
		// File exist already
		return
	}
	err = os.WriteFile(outPath, data, 07777)
	if err != nil {
		log.Fatalf("gorming: %s \n", err.Error())
	}
}
