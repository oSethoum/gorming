package parser

import (
	"reflect"
	"strings"

	"github.com/oSethoum/gorming/types"
)

func cleanString(s string, parts ...string) string {
	for _, part := range parts {
		s = strings.ReplaceAll(s, part, "")
	}
	return s
}

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
	jsonTagString := cleanString(f.Tag.Get("json"), " ")

	if len(jsonTagString) > 0 {
		jsonTag := types.JsonTag{
			OmitEmpty: strings.Contains(jsonTagString, "omitempty"),
			Ignore:    strings.Contains(jsonTagString, "-"),
		}
		jsonTag.Name = cleanString(jsonTagString, "omitempty", ",", "-")
		tags.Json = jsonTag
	}

	gormTagString := cleanString(f.Tag.Get("gorm"), " ")
	if len(gormTagString) > 0 {
		gormTag := types.GormTag{}
		gormTagStringArray := strings.Split(gormTagString, ";")
		for _, value := range gormTagStringArray {
			if value == "unique" {
				gormTag.Unique = true
			}
			if strings.HasPrefix(value, "foreignKey:") {
				gormTag.ForeignKey = cleanString(value, "foreignKey:")
			}
			if strings.HasPrefix(value, "references:") {
				gormTag.References = cleanString(value, "references:")
			}
			if strings.HasPrefix(value, "column:") {
				gormTag.Column = cleanString(value, "column:")
			}
		}
		tags.Gorm = gormTag
	}

	gormingTagString := cleanString(f.Tag.Get("gorming"), " ")
	if len(gormingTagString) > 0 {
		gormingTag := types.GormingTag{}
		for _, value := range strings.Split(gormingTagString, ":") {
			if strings.HasPrefix(value, "enum") {
				gormingTag.Enum = strings.Split(cleanString(value, "enum:"), ",")
			}
		}
		tags.Gorming = gormingTag
	}
	return tags
}

func choice(primary string, backups ...string) string {
	if primary != "" {
		return primary
	} else {
		for _, backup := range backups {
			if backup != "" {
				return backup
			}
		}
	}
	return ""
}
