package types

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
	Gorm    *GormTag    `json:"gorm,omitempty"`
	Gorming *GormingTag `json:"gorming,omitempty"`
	Json    *JsonTag    `json:"json,omitempty"`
}

type GormTag struct {
	Name       string `json:"name,omitempty"`
	Default    string `json:"default,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	ForeignKey string `json:"foreign_key,omitempty"`
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
	LocalKey string `json:"local_key,omitempty"`
	TableKey string `json:"table_key,omitempty"`
}

type Column struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	RawType string `json:"raw_type,omitempty"`
	Edge    *Edge  `json:"edge,omitempty"`
	Tags    Tags   `json:"tags,omitempty"`
}

type Config struct {
	Files []File `json:"files,omitempty"`
}
