package parser

import (
	"reflect"
	"strings"

	"github.com/oSethoum/gorming/types"
	"github.com/oSethoum/gorming/utils"
)

func fields(fieldsMap *types.FieldMap, s reflect.Type) {
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			fields(fieldsMap, f.Type)
		} else {
			(*fieldsMap)[f.Name] = f
		}
	}
}

func tags(f reflect.StructField) types.Tags {
	tags := types.Tags{}
	jsonTagString := utils.CleanString(f.Tag.Get("json"), " ")

	if len(jsonTagString) > 0 {
		jsonTag := types.JsonTag{
			OmitEmpty: strings.Contains(jsonTagString, "omitempty"),
			Ignore:    strings.Contains(jsonTagString, "-"),
		}
		jsonTag.Name = utils.CleanString(jsonTagString, "omitempty", ",", "-")
		tags.Json = jsonTag
	}

	gormTagString := utils.CleanString(f.Tag.Get("gorm"), " ")
	if len(gormTagString) > 0 {
		gormTag := types.GormTag{
			PrimaryKey: strings.Contains(gormTagString, "primarykey"),
		}
		gormTagStringArray := strings.Split(gormTagString, ";")
		for _, value := range gormTagStringArray {

			if value == "unique" {
				gormTag.Unique = true
			}

			if strings.Contains(value, "OnUpdate:CASCADE") {
				gormTag.OnUpdate = "CASCADE"
			}

			if strings.Contains(value, "OnUpdate:SET NLL") {
				gormTag.OnUpdate = "SET NLL"
			}

			if strings.Contains(value, "OnUpdate:RESTRICT") {
				gormTag.OnUpdate = "RESTRICT"
			}

			if strings.Contains(value, "OnUpdate:NO ACTION") {
				gormTag.OnUpdate = "NO ACTION"
			}

			if strings.Contains(value, "OnDelete:CASCADE") {
				gormTag.OnDelete = "CASCADE"
			}

			if strings.Contains(value, "OnDelete:SET NLL") {
				gormTag.OnDelete = "SET NLL"
			}

			if strings.Contains(value, "OnDelete:RESTRICT") {
				gormTag.OnDelete = "RESTRICT"
			}

			if strings.Contains(value, "OnDelete:NO ACTION") {
				gormTag.OnDelete = "NO ACTION"
			}

			if strings.HasPrefix(value, "foreignKey:") {
				gormTag.ForeignKey = utils.CleanString(value, "foreignKey:")
			}
			if strings.HasPrefix(value, "references:") {
				gormTag.References = utils.CleanString(value, "references:")
			}
			if strings.HasPrefix(value, "column:") {
				gormTag.Column = utils.CleanString(value, "column:")
			}
			if strings.HasPrefix(value, "default:") {
				gormTag.Default = utils.CleanString(value, "default:")
			}
			if strings.HasPrefix(value, "many2many:") {
				gormTag.Many2Many = utils.CleanString(value, "many2many:")
			}
		}
		tags.Gorm = gormTag
	}

	swaggerTagString := strings.TrimSpace(f.Tag.Get("swagger"))
	if len(swaggerTagString) > 0 {
		swaggerTag := types.SwaggerTag{}
		for _, value := range strings.Split(swaggerTagString, ";") {
			if strings.HasPrefix(value, "type=") {
				swaggerTag.Type = strings.TrimPrefix(value, "type=")
			}
			if strings.HasPrefix(value, "example=") {
				swaggerTag.Example = strings.TrimPrefix(value, "example=")
			}
		}
		tags.Swagger = swaggerTag
	}

	typescriptTagString := strings.TrimSpace(f.Tag.Get("typescript"))
	if len(typescriptTagString) > 0 {
		typescriptTag := types.TypescriptTag{}
		for _, value := range strings.Split(typescriptTagString, ";") {
			if strings.HasPrefix(value, "type=") {
				typescriptTag.Type = strings.TrimPrefix(value, "type=")
			}
			if strings.HasPrefix(value, "enum=") {
				typescriptTag.Enum = strings.Split(utils.CleanString(value, "enum="), ",")
			}
			if strings.Contains(value, "skipEdge") {
				typescriptTag.SkipEdge = true
			}

			typescriptTag.Optional = strings.Contains(value, "optional")
		}
		tags.Typescript = typescriptTag
	}

	validatorTagString := strings.TrimSpace(f.Tag.Get("validate"))
	if len(validatorTagString) > 0 {

		validatorTag := []types.ValidatorTag{}
		for _, value := range strings.Split(validatorTagString, ";") {
			if strings.HasPrefix(value, "notEmpty") {
				validatorTag = append(validatorTag, types.ValidatorTag{
					Rule: "notEmpty",
				})
			}
			if strings.HasPrefix(value, "minLen=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "minLen",
						Parameter: strings.TrimPrefix(value, "minLen="),
					},
				)
			}
			if strings.HasPrefix(value, "maxLen=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "maxLen",
						Parameter: strings.TrimPrefix(value, "maxLen="),
					},
				)
			}
			if strings.HasPrefix(value, "url") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "url",
					},
				)

			}
			if strings.HasPrefix(value, "alphaSpace") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "alphaSpace",
					},
				)
			}
			if strings.HasPrefix(value, "numeric") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "numeric",
					},
				)
			}

			if strings.HasPrefix(value, "alpha") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "alpha",
					},
				)
			}
			if strings.HasPrefix(value, "alphanumeric") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "alphanumeric",
					},
				)
			}
			if strings.HasPrefix(value, "cron") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "cron",
					},
				)
			}
			if strings.HasPrefix(value, "email") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule: "email",
					},
				)
			}
			if strings.HasPrefix(value, "match=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "match",
						Parameter: strings.TrimPrefix(value, "prefix="),
					},
				)
			}
			if strings.HasPrefix(value, "in=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "in",
						Parameter: strings.TrimPrefix(value, "in="),
					},
				)
			}
			if strings.HasPrefix(value, "out=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "out",
						Parameter: strings.TrimPrefix(value, "out="),
					},
				)
			}
			if strings.HasPrefix(value, "min=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "min",
						Parameter: strings.TrimPrefix(value, "min="),
					},
				)
			}
			if strings.HasPrefix(value, "max=") {
				validatorTag = append(validatorTag,
					types.ValidatorTag{
						Rule:      "max",
						Parameter: strings.TrimPrefix(value, "max="),
					},
				)
			}
		}
		tags.Validator = validatorTag
	}

	return tags
}
