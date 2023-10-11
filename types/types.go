package types

import "reflect"

type File uint

const (
	DB = iota
	Roues
	Handler
	Query
	Error
	Response
	Subscription
	Privacy
)

type TypeMap map[string]reflect.Type
type FieldMap map[string]reflect.StructField

type Schema struct {
	Tables []Table `json:"tables,omitempty"`
	Types  []Table `json:"types,omitempty"`
}

type Table struct {
	Name    string   `json:"name,omitempty"`
	Table   string   `json:"table,omitempty"`
	Columns []Column `json:"columns,omitempty"`
}

type Tags struct {
	Gorm    GormTag    `json:"gorm,omitempty"`
	Gorming GormingTag `json:"gorming,omitempty"`
	Json    JsonTag    `json:"json,omitempty"`
}

type GormTag struct {
	Column     string `json:"column,omitempty"`
	Default    string `json:"default,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	ForeignKey string `json:"foreign_key,omitempty"`
	References string `json:"reference,omitempty"`
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

type Config struct {
	Files []File `json:"files,omitempty"`
}
