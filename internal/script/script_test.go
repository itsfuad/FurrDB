package script

import (
	"testing"
)

func TestRegisterAndRunScript(t *testing.T) {
	hash := RegisterScript("SET foo bar; GET foo")
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	res, err := RunScript(hash, nil)
	if err != nil {
		t.Fatal(err)
	}
	if res != "bar" {
		t.Errorf("expected bar, got %s", res)
	}
}

func TestEvalScript(t *testing.T) {
	res, err := EvalScript("SET baz qux; GET baz", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res != "qux" {
		t.Errorf("expected qux, got %s", res)
	}
}
