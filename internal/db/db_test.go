package db

import (
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
