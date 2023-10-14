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
		}
		tags.Gorm = gormTag
	}

	gormingTagString := utils.CleanString(f.Tag.Get("gorming"), " ")
	if len(gormingTagString) > 0 {
		gormingTag := types.GormingTag{}
		for _, value := range strings.Split(gormingTagString, ":") {
			if strings.HasPrefix(value, "enum") {
				gormingTag.Enum = strings.Split(utils.CleanString(value, "enum:"), ",")
			}
		}
		tags.Gorming = gormingTag
	}
	return tags
}
