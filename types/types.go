package types

import "reflect"

type File uint
type DBKind string
type Case float64
type Server int
type FilesAction bool
type IgnoreHandler byte

const (
	FileDB = iota + 10
	FileRoutes
	FileHandler
	FileHooks
	FileMigration
	FileQuery
	FileError
	FileResponse
	FileMiddlewares
	FileWs
	FilePrivacy
	FileUtils
	FileApi
	FileResource
	FileSchema
	FileTsApi
	FileTsTypes
)

const (
	SQLite   DBKind = "sqlite"
	MySQL    DBKind = "mysql"
	Postgres DBKind = "postgres"
)

const (
	Camel Case = iota + 1
	Pascal
	Snake
)

const (
	Fiber Server = iota + 40
	Wails
)

const (
	Generate      FilesAction = true
	DoNotGenerate FilesAction = false
)

type Engine = func(tables []any, types ...any)
type TypeMap map[string]reflect.Value
type FieldMap map[string]reflect.StructField

type Paths struct {
	BasePath         string   `json:"base_path,omitempty"`
	BackendPath      string   `json:"backend_path,omitempty"`
	TypescriptClient []string `json:"typescript_client,omitempty"`
	ApiPath          string   `json:"api_path,omitempty"`
}

type Config struct {
	DBKind         DBKind            `json:"db_kind,omitempty"`
	Case           Case              `json:"case,omitempty"`
	Server         Server            `json:"server,omitempty"`
	FilesAction    FilesAction       `json:"files_action,omitempty"`
	Paths          Paths             `json:"paths,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	Files          []File            `json:"files,omitempty"`
	Package        string            `json:"package,omitempty"`
	ApiPackage     string            `json:"api_package,omitempty"`
	BackendPackage string            `json:"backend_package,omitempty"`
	SkipRoutes     map[string]string `json:"skip_routes,omitempty"`
}

type Schema struct {
	Tables []Table `json:"tables,omitempty"`
	Types  []Table `json:"types,omitempty"`
}

type Table struct {
	Name         string   `json:"name,omitempty"`
	Table        string   `json:"table,omitempty"`
	HasTableFunc bool     `json:"has_table_func,omitempty"`
	Columns      []Column `json:"columns,omitempty"`
	Skip         []string `json:"skip,omitempty"`
}

type ValidatorTag struct {
	Rule      string `json:"rule,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type Tags struct {
	Gorm       GormTag        `json:"gorm,omitempty"`
	Gorming    GormingTag     `json:"gorming_tag,omitempty"`
	Json       JsonTag        `json:"json,omitempty"`
	Typescript TypescriptTag  `json:"typescript,omitempty"`
	Validator  []ValidatorTag `json:"validator,omitempty"`
	Swagger    SwaggerTag     `json:"swagger,omitempty"`
	Dart       DartTag        `json:"dart,omitempty"`
}

type TemplateData struct {
	Schema         *Schema `json:"schema,omitempty"`
	ValidationList string  `json:"validation_list,omitempty"`
	Config         Config  `json:"config,omitempty"`
}

type GormTag struct {
	Column     string `json:"column,omitempty"`
	Default    string `json:"default,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	ForeignKey string `json:"foreign_key,omitempty"`
	References string `json:"reference,omitempty"`
	Ignore     bool   `json:"ignore,omitempty"`
	Many2Many  string `json:"many2many,omitempty"`
}

type DartTag struct {
	Type string `json:"type,omitempty"`
}

type SwaggerTag struct {
	Type    string `json:"type,omitempty"`
	Example string `json:"example,omitempty"`
}

type TypescriptTag struct {
	Type string   `json:"type,omitempty"`
	Enum []string `json:"enum,omitempty"`
}

type GormingTag struct {
	TsType      string   `json:"ts_type,omitempty"`
	DartType    string   `json:"dart_type,omitempty"`
	SwaggerType string   `json:"swagger_type,omitempty"`
	Skip        []string `json:"skip,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type JsonTag struct {
	Name      string `json:"name,omitempty"`
	Ignore    bool   `json:"ignore,omitempty"`
	OmitEmpty bool   `json:"omit_empty,omitempty"`
}

type Edge struct {
	Table     string `json:"table,omitempty"`
	Unique    bool   `json:"unique,omitempty"`
	LocalKey  string `json:"local_key,omitempty"`
	TableKey  string `json:"table_key,omitempty"`
	Many2Many string `json:"many2many,omitempty"`
}

type Column struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	RawType string `json:"raw_type,omitempty"`
	Edge    *Edge  `json:"edge,omitempty"`
	Slice   bool   `json:"slice,omitempty"`
	Tags    Tags   `json:"tags,omitempty"`
}
