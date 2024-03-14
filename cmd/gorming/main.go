package main

import "embed"

//go:embed templates
var templates embed.FS

func main() {
	writeTemplate("generate", "./generate/generate.go")
	writeTemplate("main", "./generate/main.go")
	writeTemplate("models", "./db/generate.go")
}
