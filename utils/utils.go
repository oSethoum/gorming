package utils

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
)

func Camel(s string) string {
	return strcase.ToLowerCamel(s)
}

func Snake(s string) string {
	return strcase.ToSnake(s)
}

func Choice[T comparable](primary T, backups ...T) T {
	var empty T

	if primary != empty {
		return primary
	} else {
		for _, backup := range backups {
			if backup != empty {
				return backup
			}
		}
	}
	return empty
}

func In[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func CleanString(s string, parts ...string) string {
	for _, part := range parts {
		s = strings.ReplaceAll(s, part, "")
	}
	return s
}

func CurrentGoMod() (string, string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

label:
	if data, err := os.ReadFile(path.Join(cwd, "go.mod")); err == nil {
		return cwd, CleanString(strings.Split(strings.Split(string(data), "\n")[0], " ")[1])
	} else {
		cwd = path.Dir(cwd)
		if len(cwd) == 0 {
			log.Fatalln("cannot find go.mod")
		}
		goto label
	}
}
