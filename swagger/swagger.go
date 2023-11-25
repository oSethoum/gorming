package swagger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

type Info struct {
	Title   string `json:"title,omitempty"`
	Version string `json:"version"`
}

type Swagger struct {
	OpenAPI    string          `json:"openapi,omitempty"`
	Info       *Info           `json:"info,omitempty"`
	Paths      map[string]Path `json:"paths,omitempty"`
	Components Components      `json:"components,omitempty"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas,omitempty"`
}

type Path struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Parameter struct {
	In          string  `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type RequestBody struct {
	Description string  `json:"description,omitempty"`
	Content     Content `json:"content,omitempty"`
}

type Response struct {
	Description string  `json:"description,omitempty"`
	Content     Content `json:"content,omitempty"`
}

type Content map[string]struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Ref        string            `json:"$ref,omitempty"`
	Type       string            `json:"type,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Enum       []string          `json:"enum,omitempty"`
	Required   []string          `json:"required,omitempty"`
}

func skipRoute(config *types.Config, resource, method string) bool {
	if config.SkipRoutes != nil {
		if skips, ok := config.SkipRoutes[resource]; ok {
			return strings.Contains(skips, method)
		}
	}
	return false
}

func skipAllRoute(config *types.Config, resource string) bool {
	if config.SkipRoutes != nil {
		if skips, ok := config.SkipRoutes[resource]; ok {
			return strings.Contains(skips, "all") ||
				(strings.Contains(skips, "query") &&
					strings.Contains(skips, "create") &&
					strings.Contains(skips, "update") &&
					strings.Contains(skips, "delete"))
		}
	}
	return false
}

func columnName(config *types.Config, column types.Column) string {
	if config.Case == types.Camel {
		return utils.Camel(column.Name)
	}

	if config.Case == types.Snake {
		return utils.Snake(column.Name)
	}

	if len(column.Tags.Json.Name) > 0 {
		return column.Tags.Json.Name
	}

	return column.Name
}

func tableName(config *types.Config, name string) string {
	if config.Case == types.Camel {
		return utils.Camels(name)
	}

	if config.Case == types.Snake {
		return utils.Snakes(name)
	}

	return name
}

func columnType(t string) string {

	m := map[string]string{
		"string": "string",
		"uint":   "number",
		"int":    "number",
		"float":  "number",
		"bool":   "boolean",
	}

	for k, v := range m {
		if strings.HasPrefix(t, k) {
			return v
		}
	}

	if strings.Contains(t, "time.Time") {
		return "string"
	}

	if strings.Contains(t, "gorm.DeletedAt") {
		return "string"
	}

	return "object"
}

func _200(name string, operation string) Response {
	s := Schema{Type: "array",
		Items: &Schema{
			Ref: "#/components/schemas/" + name,
		},
	}

	if operation == "query" {
		s = Schema{
			Type: "object", Properties: map[string]Schema{
				"count":  {Type: "number"},
				"result": s,
			},
		}
	}

	return Response{
		Description: "Success",
		Content: Content{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]Schema{
						"status": {Type: "string", Enum: []string{"success"}},
						"data":   s,
					},
				},
			},
		},
	}
}

func Generate(config *types.Config, Tables []types.Table, Types ...types.Table) Swagger {
	docs := Swagger{
		OpenAPI: "3.0.2",
		Info: &Info{
			Title:   config.SwaggerConfig.Title,
			Version: config.SwaggerConfig.Version,
		},
	}
	components := Components{
		Schemas: map[string]Schema{},
	}
	paths := map[string]Path{}

	// standard responses

	_401 := Response{
		Description: "Unauthorized request",
		Content: Content{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]Schema{
						"status": {
							Type: "string",
							Enum: []string{"error"},
						},
						"error": {
							Type: "object",
							Properties: map[string]Schema{
								"type":    {Type: "string", Enum: []string{"authorization"}},
								"message": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}

	_400 := Response{
		Description: "Bad request",
		Content: Content{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]Schema{
						"status": {
							Type: "string",
							Enum: []string{"error"},
						},
						"error": {
							Type: "object",
							Properties: map[string]Schema{
								"type":    {Type: "string", Enum: []string{"authentication", "validation", ""}},
								"message": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}

	_500 := Response{
		Description: "Internal server error",
		Content: Content{
			"application/json": {
				Schema: &Schema{
					Type: "object",
					Properties: map[string]Schema{
						"status": {
							Type: "string",
							Enum: []string{"error"},
						},
						"error": {
							Type: "object",
							Properties: map[string]Schema{
								"type":    {Type: "string", Enum: []string{"authentication", "validation", ""}},
								"message": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}

	for _, t := range Tables {

		if skipAllRoute(config, tableName(config, t.Name)) {
			continue
		}

		p := Path{}

		if !skipRoute(config, tableName(config, t.Name), "query") {
			p.Get = &Operation{
				Tags:        []string{tableName(config, t.Name)},
				Description: fmt.Sprintf("Query %s", tableName(config, t.Name)),
				Parameters: []Parameter{
					{
						In:          "query",
						Name:        "query",
						Description: "JSON format of " + t.Name + "Query",
						Schema: &Schema{
							Type: "string",
						},
					},
				},
				Responses: map[string]Response{
					"200": _200(t.Name, "query"),
					"401": _401,
					"400": _400,
					"500": _500,
				},
			}
		}

		if !skipRoute(config, tableName(config, t.Name), "create") {
			p.Post = &Operation{
				Tags:        []string{tableName(config, t.Name)},
				Description: fmt.Sprintf("Insert %s", tableName(config, t.Name)),
				RequestBody: &RequestBody{
					Content: Content{
						"application/json": {
							Schema: &Schema{
								Type: "array",
								Items: &Schema{
									Ref: "#/components/schemas/" + t.Name + "Create",
								},
							},
						},
					},
				},
				Responses: map[string]Response{
					"200": _200(t.Name, "create"),
					"401": _401,
					"400": _400,
					"500": _500,
				},
			}
		}

		if !skipRoute(config, tableName(config, t.Name), "update") {
			p.Patch = &Operation{
				Tags:        []string{tableName(config, t.Name)},
				Description: fmt.Sprintf("Update %s", tableName(config, t.Name)),
				RequestBody: &RequestBody{
					Content: Content{
						"application/json": {
							Schema: &Schema{
								Type: "array",
								Items: &Schema{
									Ref: "#/components/schemas/" + t.Name + "Update",
								},
							},
						},
					},
				},
				Responses: map[string]Response{
					"200": _200(t.Name, "update"),
					"401": _401,
					"400": _400,
					"500": _500,
				},
			}
		}

		if !skipRoute(config, tableName(config, t.Name), "delete") {
			p.Delete = &Operation{
				Tags:        []string{tableName(config, t.Name)},
				Description: fmt.Sprintf("Delete %s", tableName(config, t.Name)),
				RequestBody: &RequestBody{
					Content: Content{
						"application/json": {
							Schema: &Schema{
								Ref: "#/components/schemas/" + t.Name + "Where",
							},
						},
					},
				},
				Responses: map[string]Response{
					"200": _200(t.Name, "delete"),
					"401": _401,
					"400": _400,
					"500": _500,
				},
			}
		}

		paths["/"+tableName(config, t.Name)] = p

		st := Schema{
			Type:       "object",
			Properties: map[string]Schema{},
		}

		sc := Schema{
			Type:       "object",
			Properties: map[string]Schema{},
		}

		su := Schema{
			Type:       "object",
			Properties: map[string]Schema{},
		}

		sq := Schema{
			Type: "object",
			Properties: map[string]Schema{
				"limit":  {Type: "number"},
				"offset": {Type: "number"},
				"where":  {Ref: "#/components/schemas/" + t.Name + "Where"},
			},
		}

		selects := []string{}
		omits := []string{}
		edges := []types.Column{}

		for _, c := range t.Columns {
			omits = append(omits, columnName(config, c))
			if c.Tags.Gorming.SwaggerType != "" {
				c.Type = c.Tags.Gorming.SwaggerType
			}

			if c.Tags.Json.Ignore {
				continue
			}
			if c.Edge != nil {
				edges = append(edges, c)
				if c.Edge.Unique {
					st.Properties[columnName(config, c)] = Schema{
						Ref: "#/components/schemas/" + c.RawType,
					}
					sc.Properties[columnName(config, c)] = Schema{
						Ref: "#/components/schemas/" + c.RawType + "Create",
					}
					su.Properties[columnName(config, c)] = Schema{
						Ref: "#/components/schemas/" + c.RawType + "Update",
					}
				} else {
					st.Properties[columnName(config, c)] = Schema{
						Type: "array",
						Items: &Schema{
							Ref: "#/components/schemas/" + c.RawType,
						},
					}
					sc.Properties[columnName(config, c)] = Schema{
						Type: "array",
						Items: &Schema{
							Ref: "#/components/schemas/" + c.RawType + "Create",
						},
					}
					su.Properties[columnName(config, c)] = Schema{
						Type: "array",
						Items: &Schema{
							Ref: "#/components/schemas/" + c.RawType + "Update",
						},
					}
				}
			} else {
				selects = append(selects, columnName(config, c))
				if c.Slice {
					if unicode.IsUpper(rune(c.Type[0])) {
						st.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Ref: "#/components/schemas/" + c.RawType,
							},
						}
						sc.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Ref: "#/components/schemas/" + c.RawType,
							},
						}
						su.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Ref: "#/components/schemas/" + c.RawType,
							},
						}
					} else {
						cType := strings.ReplaceAll(c.Type, "[]", "")
						st.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Type: columnType(cType),
							},
						}
						sc.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Type: columnType(cType),
							},
						}
						su.Properties[columnName(config, c)] = Schema{
							Type: "array",
							Items: &Schema{
								Type: columnType(cType),
							},
						}
					}
				} else {
					if unicode.IsUpper(rune(c.Type[0])) {
						st.Properties[columnName(config, c)] = Schema{
							Ref: "#/components/schemas/" + c.Type,
						}
						if columnName(config, c) != "id" {
							sc.Properties[columnName(config, c)] = Schema{
								Ref: "#/components/schemas/" + c.Type,
							}
						}
						su.Properties[columnName(config, c)] = Schema{
							Ref: "#/components/schemas/" + c.Type,
						}
					} else {
						if len(c.Tags.Gorming.Enum) > 0 {
							st.Properties[columnName(config, c)] = Schema{
								Type: columnType(c.RawType),
								Enum: c.Tags.Gorming.Enum,
							}
							su.Properties[columnName(config, c)] = Schema{
								Type: columnType(c.RawType),
								Enum: c.Tags.Gorming.Enum,
							}
							if columnName(config, c) != "id" {
								sc.Properties[columnName(config, c)] = Schema{
									Type: columnType(c.RawType),
									Enum: c.Tags.Gorming.Enum,
								}
							}
						} else {
							st.Properties[columnName(config, c)] = Schema{
								Type: columnType(c.Type),
							}
							if columnName(config, c) != "id" {
								sc.Properties[columnName(config, c)] = Schema{
									Type: columnType(c.Type),
								}
							}
							su.Properties[columnName(config, c)] = Schema{
								Type: columnType(c.Type),
							}
						}

					}
				}
			}
		}

		sw := Schema{
			Type: "object",
			Properties: map[string]Schema{
				"not": {
					Ref: "#/components/schemas/" + t.Name + "Where",
				},
				"and": {
					Type: "array",
					Items: &Schema{
						Ref: "#/components/schemas/" + t.Name + "Where",
					},
				},
				"or": {
					Type: "array",
					Items: &Schema{
						Ref: "#/components/schemas/" + t.Name + "Where",
					},
				},
				"field": {
					Type: "object",
					Properties: map[string]Schema{
						"name": {
							Type: "string",
							Enum: selects,
						},
						"predicate": {
							Type: "string",
							Enum: []string{"between", "like", "in", "not in", "=", "<>", ">", ">=", "<", "<=", "null", "not null"},
						},
						"value": {
							Type: "object",
						},
					},
					Required: []string{"name", "predicate"},
				},
			},
		}

		sq.Properties["select"] = Schema{
			Type: "array",
			Items: &Schema{
				Type: "string",
				Enum: selects,
			},
		}
		sq.Properties["omit"] = Schema{
			Type: "array",
			Items: &Schema{
				Type: "string",
				Enum: omits,
			},
		}

		sp := Schema{
			Type:       "object",
			Properties: map[string]Schema{},
		}

		for _, e := range edges {
			sp.Properties[columnName(config, e)] = Schema{
				Ref: "#/components/schemas/" + e.RawType + "Query",
			}
		}
		sq.Properties["preloads"] = sp
		sq.Properties["orders"] = Schema{
			Type: "object",
			Properties: map[string]Schema{
				"field": {
					Type: "string",
					Enum: selects,
				},
				"direction": {
					Type: "string",
					Enum: []string{"ASC", "DESC"},
				},
			},
		}
		components.Schemas[t.Name] = st
		components.Schemas[t.Name+"Query"] = sq
		components.Schemas[t.Name+"Where"] = sw
		components.Schemas[t.Name+"Create"] = sc
		components.Schemas[t.Name+"Update"] = su
	}

	for _, t := range Types {
		tc := Schema{
			Type:       "object",
			Properties: map[string]Schema{},
		}

		for _, c := range t.Columns {
			if c.Slice {
				tc.Properties[columnName(config, c)] = Schema{
					Type: "array",
					Items: &Schema{
						Type: c.Type,
					},
				}
			} else {
				tc.Properties[columnName(config, c)] = Schema{
					Type: c.Type,
				}
			}
		}
		components.Schemas[t.Name] = tc
	}

	if len(config.SwaggerConfig.PreservePaths) > 0 || len(config.SwaggerConfig.PreserveSchemas) > 0 {
		currentDocs := Swagger{}
		cwd, _ := os.Getwd()
		file, err := os.ReadFile(filepath.Join(cwd, config.SwaggerConfig.Output, config.SwaggerConfig.FileName+".json"))

		if err == nil {
			if err := json.Unmarshal(file, &currentDocs); err == nil {
				for _, v := range config.SwaggerConfig.PreservePaths {
					for ck, cp := range currentDocs.Paths {
						if v == ck {
							paths[v] = cp
						}
					}
				}

				for _, v := range config.SwaggerConfig.PreserveSchemas {
					for k, s := range currentDocs.Components.Schemas {
						if v == k {
							components.Schemas[v] = s
						}
					}
				}
			}
		}
	}

	docs.Paths = paths
	docs.Components = components

	return docs
}
