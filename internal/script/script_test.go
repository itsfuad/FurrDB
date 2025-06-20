package script

import (
	"furr/internal/db"
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

func TestScriptDSL_LetIfEnd(t *testing.T) {
	scriptStr := `LET x = GET foo; IF x == bar; SET foo baz; END; GET foo`
	_, _ = db.Commands["SET"]([]string{"foo", "bar"})
	res, err := EvalScript(scriptStr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "baz" {
		t.Errorf("expected baz, got %s", res)
	}
}

func TestScriptDSL_Sandbox(t *testing.T) {
	scriptStr := `FLUSHDB; GET foo`
	_, _ = db.Commands["SET"]([]string{"foo", "bar"})
	_, err := EvalScript(scriptStr, nil)
	if err == nil || err.Error() != "ERR command FLUSHDB not allowed in script on line 1" {
		t.Errorf("expected sandbox error, got %v", err)
	}
}

func TestScriptDSL_LetSyntaxError(t *testing.T) {
	scriptStr := `LET x GET foo`
	_, err := EvalScript(scriptStr, nil)
	if err == nil || err.Error() != "ERR invalid LET syntax on line 1" {
		t.Errorf("expected LET syntax error, got %v", err)
	}
}
