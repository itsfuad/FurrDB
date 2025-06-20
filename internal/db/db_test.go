package db

import (
	"strings"
	"testing"
)

func TestStoreBasicOps(t *testing.T) {
	db := NewStore()
	db.mu.Lock()
	db.data["foo"] = "bar"
	db.mu.Unlock()

	db.mu.RLock()
	v := db.data["foo"]
	db.mu.RUnlock()
	if v != "bar" {
		t.Errorf("expected bar, got %s", v)
	}
}

func TestSetGetHandler(t *testing.T) {
	_, err := setHandler([]string{"foo", "bar"})
	if err != nil {
		t.Fatal(err)
	}
	val, err := getHandler([]string{"foo"})
	if err != nil {
		t.Fatal(err)
	}
	if val != "bar" {
		t.Errorf("expected bar, got %s", val)
	}
}

func TestDelExistsHandler(t *testing.T) {
	_, _ = setHandler([]string{"baz", "qux"})
	exists, _ := existsHandler([]string{"baz"})
	if exists != "1" {
		t.Errorf("expected exists=1, got %s", exists)
	}
	_, _ = delHandler([]string{"baz"})
	exists, _ = existsHandler([]string{"baz"})
	if exists != "0" {
		t.Errorf("expected exists=0, got %s", exists)
	}
}

func TestListCommands(t *testing.T) {
	DefaultStore = NewStore()
	_, _ = lpushHandler([]string{"mylist", "a", "b"}) // b, a
	_, _ = rpushHandler([]string{"mylist", "c"})      // b, a, c
	val, _ := lpopHandler([]string{"mylist"})         // b
	if val != "b" {
		t.Errorf("expected b, got %s", val)
	}
	val, _ = rpopHandler([]string{"mylist"}) // c
	if val != "c" {
		t.Errorf("expected c, got %s", val)
	}
	_, _ = lpushHandler([]string{"mylist", "x"}) // x, a
	out, _ := lrangeHandler([]string{"mylist", "0", "1"})
	if out != "x,a" {
		t.Errorf("expected x,a, got %s", out)
	}
}

func TestSetCommands(t *testing.T) {
	DefaultStore = NewStore()
	_, _ = saddHandler([]string{"myset", "a", "b", "c"})
	_, _ = sremHandler([]string{"myset", "b"})
	out, _ := smembersHandler([]string{"myset"})
	if !strings.Contains(out, "a") || !strings.Contains(out, "c") || strings.Contains(out, "b") {
		t.Errorf("expected a and c, not b; got %s", out)
	}
}

func TestMetaCommands(t *testing.T) {
	DefaultStore = NewStore()
	_, _ = setHandler([]string{"k1", "v1"})
	_, _ = setHandler([]string{"k2", "v2"})
	_, _ = saddHandler([]string{"s1", "x"})
	keys, _ := keysHandler(nil)
	if !strings.Contains(keys, "k1") || !strings.Contains(keys, "k2") || !strings.Contains(keys, "s1") {
		t.Errorf("expected all keys, got %s", keys)
	}
	info, _ := infoHandler(nil)
	if !strings.Contains(info, "keys:3") {
		t.Errorf("expected keys:3, got %s", info)
	}
	_, _ = flushdbHandler(nil)
	info, _ = infoHandler(nil)
	if !strings.Contains(info, "keys:0") {
		t.Errorf("expected keys:0 after flush, got %s", info)
	}
}
