package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

func Camel(s string) string {
	return strcase.ToLowerCamel(s)
}

func Camels(s string) string {
	return strcase.ToLowerCamel(inflection.Plural(s))
}

func Snake(s string) string {
	return strcase.ToSnake(s)
}

func Snakes(s string) string {
	return strcase.ToSnake(inflection.Plural(s))
}

func ArrayChoice(primary []string, backups ...[]string) []string {

	if len(primary) > 0 {
		return primary
	} else if len(backups) > 0 {
		for _, backup := range backups {
			if len(backup) > 0 {
				return backup
			}
		}
	}
	return primary
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

	for {
		if data, err := os.ReadFile(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd, CleanString(strings.Split(strings.Split(string(data), "\n")[0], " ")[1])
		}

		if cwd == "." {
			log.Fatalln("gorming: cannot find go.mod")
		}

		cwd = filepath.Dir(cwd)
	}

}
