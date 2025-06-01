# Gorming

Gorming is a powerful code generation tool designed to streamline the process of generating query and TypeScript clients, Go Fiber handlers, and Swagger documentation. It simplifies the development workflow by automating repetitive tasks and ensuring consistency across your project.

## Getting Started

To install Gorming, use the following:

```bash
go install github.com/oSethoum/gorming/cmd/gorming@latest
```

## Init in new project

```bash
gorming
```

this will result to this folder structure

```
app
│
└───db
│   │   db.go
│   │   models.go
│
│
└───generate
    │   generate.go
    │   main.go
```

where each file will contain:

-  db.go

where the connection to the database happen

```go
package db

import (
	"fmt"
	"github.com/oSethoum/sqlite"

	"gorm.io/gorm"
)

var Client *gorm.DB

func Init() error {

	dsn := fmt.Sprintf("file:%s?_fk=1&_pragma_key=%s", "db.sqlite", "")
	dialect := sqlite.Open(dsn)

	client, err := gorm.Open(dialect, &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return err
	}

	Client = client
	return nil
}

func Close() error {
	if Client != nil {
		conn, err := Client.DB()
		if err != nil {
			return err
		}
		return conn.Close()
	}
	return nil
}
```

-  db/models.go

where you define your models in gorm syntax

```go
package db

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint         `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

```

-  generate.go

```go
package generate

import (
	_ "github.com/oSethoum/gorming"
)

//go:generate go run main.go

```

-  main.go

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
		DBKind:       types.SQLite,
		Case:         types.Snake,
		Server:       types.Fiber,
		FilesAction:  types.Generate,
		WithSecurity: true,
		Package:      "app",
		Debug:        true,
		Files: []types.File{
			types.FileDB,
			types.FileApi,
			types.FileError,
			types.FileHandler,
			types.FileHooks,
			types.FileMigration,
			types.FilePrivacy,
			types.FileQuery,
			types.FileResource,
			types.FileResponse,
			types.FileRoutes,
			types.FileSchema,
			types.FileTsApi,
			types.FileTsTypes,
			types.FileUtils,
			types.FileValidation,
			types.FileWs,
		},
	})

	engine([]any{})
	exec.Command("gofmt", "-w", "-s", "-l", "../..").Run()
}

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

`typescript=`: this tag will help you override th default type that gorming generate, gorming default to any when the type isn't defined or primitive. the tag list would be:
| Tga | Used for | Example |
| :---: | :---: | :---: |
| enum | define enum type instead of using string or any | `typescript:"enum=individual, company" |
| type | override the type gorming outputs for the field | `typescript:"type={name:string, value:number}" |

`gorm=`: gorming is aware of gorm tags so if you wanna change the name of the foreign key, gorming will use the name provider in the gorm tag. the tags we support are: **foreignKey**, **references**, **column**, **default**, **many2many**.

``

## License

Gorming is licensed under the [License](LICENSE).
