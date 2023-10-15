package types

import "reflect"

type File uint
type DBKind string
type Case float64
type Server int
type FilesAction bool
type IgnoreHandler byte

const (
	DB = iota + 10
	Routes
	Handler
	Model
	Migrate
	Query
	Images
	Error
	Response
	Ws
	Privacy
	Utils
	DartApi
	DartTypes
	TsApi
	TsTypes
)

const (
	SQLite   = "sqlite"
	MySQL    = "mysql"
	Postgres = "postgres"
)

const (
	Camel = iota + 1
	Pascal
	Snake
)

const (
	Generate      = true
	DoNotGenerate = false
)

type Engine = func(tables []any, types ...any)
type TypeMap map[string]reflect.Value
type FieldMap map[string]reflect.StructField

type Paths struct {
	BasePath         string `json:"base_path,omitempty"`
	TypescriptClient string `json:"typescript_client,omitempty"`
	DartClient       string `json:"dart_client,omitempty"`
}

type Config struct {
	DBKind         DBKind                     `json:"db_kind,omitempty"`
	Case           Case                       `json:"case,omitempty"`
	FilesAction    FilesAction                `json:"files_action,omitempty"`
	Paths          Paths                      `json:"paths,omitempty"`
	Debug          bool                       `json:"debug,omitempty"`
	Fiels          []File                     `json:"fiels,omitempty"`
	IgnoreHandlers map[string][]IgnoreHandler `json:"ignore_handlers,omitempty"`
	Package        string                     `json:"package,omitempty"`
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
}

type Tags struct {
	Gorm    GormTag    `json:"gorm,omitempty"`
	Gorming GormingTag `json:"gorming,omitempty"`
	Json    JsonTag    `json:"json,omitempty"`
}

type TemplateData struct {
	Schema *Schema `json:"schema,omitempty"`
	Config Config  `json:"config,omitempty"`
}

type GormTag struct {
	Column     string `json:"column,omitempty"`
	Default    string `json:"default,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	ForeignKey string `json:"foreign_key,omitempty"`
	References string `json:"reference,omitempty"`
	Ignore     bool   `json:"ignore,omitempty"`
}

type GormingTag struct {
	Enum []string `json:"enum,omitempty"`
}

type JsonTag struct {
	Name      string `json:"name,omitempty"`
	Ignore    bool   `json:"ignore,omitempty"`
	OmitEmpty bool   `json:"omit_empty,omitempty"`
}

type Edge struct {
	Table    string `json:"table,omitempty"`
	Unique   bool   `json:"unique,omitempty"`
	LocalKey string `json:"local_key,omitempty"`
	TableKey string `json:"table_key,omitempty"`
}

type Column struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	RawType string `json:"raw_type,omitempty"`
	Edge    *Edge  `json:"edge,omitempty"`
	Slice   bool   `json:"slice,omitempty"`
	Tags    Tags   `json:"tags,omitempty"`
}
