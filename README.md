I apologize for the inconvenience. Here is the raw text without markdown rendering:

# Gorming

Gorming is a powerful code generation tool designed to streamline the process of generating query and TypeScript clients, Go Fiber handlers, and Swagger documentation. It simplifies the development workflow by automating repetitive tasks and ensuring consistency across your project.

## Installation

To install Gorming, use the following:

```bash
go get -u github.com/oSethoum/gorming
```

## Usage

To use Gorming, create a configuration struct and customize it according to your project requirements. Here's an example configuration:

```go
//go:build ignore
// +build ignore

package main

import (
	"app/db"
	"os/exec"

	"github.com/oSethoum/gorming"
	"github.com/oSethoum/gorming/types"
)

func main() {
	engine := gorming.New(types.Config{
		Package:      "app",
		DBKind:       types.SQLite,
		Server:       types.Fiber,
		WithSwagger:  true,
		WithSecurity: true,
		SwaggerConfig: types.SwaggerConfig{
			Output:          "../docs.json",
			BasePath:        "http://localhost:5000/api",
			Title:           "App API",
			Version:         "1.0.0",
			PreservePaths:   []string{"/auth"},
			PreserveSchemas: []string{"AuthBody"},
		},
		FilesAction: types.DoNotGenerate,
		Paths: types.Paths{
			TypescriptClient: "../../frontend/src/api",
		},
	})

    // Provide extra types
    type Address struct {
        Street string `json:"street,omitempty"`
        Zip string `json:"zip,omitempty"`
        ...
    }

	engine([]interface{}{
		db.User{},
		db.Role{},
	}, Address{})

	// Format the output
	exec.Command("gofmt", "-w", "-l", "../").Run()
}
```

This example showcases a basic setup. Customize the configuration according to your project's needs.

## Configuration Options

### `DBKind`

Specify the type of database: `SQLite`, `MySQL`, or `Postgres`.

### `Server`

Choose between `Fiber` or `Wails` as your server.

### `Case`

Define the naming convention for your code: `Camel`, `Pascal`, or `Snake`.

### `FilesAction`

Control whether to generate or skip files.

### `Paths`

Configure various paths for generated files.

### `File`

Choose specific files to generate based on your project requirements.

## Struct Tags

## License

Gorming is licensed under the [License](LICENSE).
