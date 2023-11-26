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
		gormTag := types.GormTag{}
		gormTagStringArray := strings.Split(gormTagString, ";")
		for _, value := range gormTagStringArray {
			if value == "unique" {
				gormTag.Unique = true
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

	gormingTagString := strings.TrimSpace(f.Tag.Get("gorming"))
	if len(gormingTagString) > 0 {
		gormingTag := types.GormingTag{}
		for _, value := range strings.Split(gormingTagString, ";") {
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "enum=") {
				gormingTag.Enum = strings.Split(utils.CleanString(value, "enum="), ",")
			}
			if strings.HasPrefix(value, "tsType=") {
				gormingTag.TsType = utils.CleanString(value, "tsType=")
			}
			if strings.HasPrefix(value, "swaggerType=") {
				gormingTag.SwaggerType = utils.CleanString(value, "swaggerType=")
			}
			if strings.HasPrefix(value, "dartType=") {
				gormingTag.DartType = utils.CleanString(value, "dartType=")
			}
			if strings.HasPrefix(value, "notEmpty") {
				gormingTag.Validation = append(gormingTag.Validation, types.ValidationTag{
					Rule: "notEmpty",
				})
			}
			if strings.HasPrefix(value, "minLen=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "minLen",
						Parameter: strings.TrimPrefix(value, "minLen="),
					},
				)
			}
			if strings.HasPrefix(value, "maxLen=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "maxLen",
						Parameter: strings.TrimPrefix(value, "maxLen="),
					},
				)
			}
			if strings.HasPrefix(value, "url") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "url",
					},
				)

			}
			if strings.HasPrefix(value, "alphaSpace") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "alphaSpace",
					},
				)
			}
			if strings.HasPrefix(value, "numeric") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "numeric",
					},
				)
			}

			if strings.HasPrefix(value, "alpha") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "alpha",
					},
				)
			}
			if strings.HasPrefix(value, "alphanumeric") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "alphanumeric",
					},
				)
			}
			if strings.HasPrefix(value, "cron") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "cron",
					},
				)
			}
			if strings.HasPrefix(value, "email") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule: "email",
					},
				)
			}
			if strings.HasPrefix(value, "match=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "match",
						Parameter: strings.TrimPrefix(value, "prefix="),
					},
				)
			}
			if strings.HasPrefix(value, "in=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "in",
						Parameter: strings.TrimPrefix(value, "in="),
					},
				)
			}
			if strings.HasPrefix(value, "out=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "out",
						Parameter: strings.TrimPrefix(value, "out="),
					},
				)
			}
			if strings.HasPrefix(value, "min=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "min",
						Parameter: strings.TrimPrefix(value, "min="),
					},
				)
			}
			if strings.HasPrefix(value, "max=") {
				gormingTag.Validation = append(gormingTag.Validation,
					types.ValidationTag{
						Rule:      "max",
						Parameter: strings.TrimPrefix(value, "max="),
					},
				)
			}

		}
		tags.Gorming = gormingTag
	}
	return tags
}
