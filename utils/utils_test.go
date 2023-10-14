package utils_test

import (
	"fmt"
	"testing"

	"github.com/oSethoum/gorming/utils"
)

func TestChoice(t *testing.T) {
	var result any
	result = utils.Choice(0, 1)
	if result == 0 {
		t.Errorf("should output 1 but got 0")
	}
	fmt.Printf("%v \n", result)

	result = utils.Choice("", "Hello")
	if result == "" {
		t.Errorf("should output Hello but got ''")
	}
	fmt.Printf("%v \n", result)

}

func TestCurrentGoMod(t *testing.T) {
	path, pkg := utils.CurrentGoMod()
	println(path)
	println(pkg)
}

func TestCamel(t *testing.T) {
	if v := utils.Camel("ID"); v != "id" {
		t.Errorf("expected id got %s", v)
	}
	if v := utils.Camel("userID"); v != "userId" {
		t.Errorf("expected userId got %s", v)
	}
	if v := utils.Camel("CreatedAt"); v != "createdAt" {
		t.Errorf("expected createdAt got %s", v)
	}
}

func TestSnake(t *testing.T) {
	if v := utils.Snake("ID"); v != "id" {
		t.Errorf("expected id got %s", v)
	}
	if v := utils.Snake("UserID"); v != "user_id" {
		t.Errorf("expected user_id got %s", v)
	}
	if v := utils.Snake("CreatedAt"); v != "created_at" {
		t.Errorf("expected created_at got %s", v)
	}
}
