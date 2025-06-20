package engine

import (
	"furr/internal/db"
	"os"
	"testing"
)

func TestAppendLoadFlush(t *testing.T) {
	// Use a temp file for AOF
	tmp := "test_aof.log"
	aofPath = tmp
	defer os.Remove(tmp)
	aofFile = nil

	// Clear DB
	db.DefaultStore = db.NewStore()

	err := Append("SET testkey testval")
	if err != nil {
		t.Fatal(err)
	}
	err = Flush()
	if err != nil {
		t.Fatal(err)
	}

	// Clear DB again
	db.DefaultStore = db.NewStore()

	err = Load()
	if err != nil {
		t.Fatal(err)
	}
	val, _ := db.Commands["GET"]([]string{"testkey"})
	if val != "testval" {
		t.Errorf("expected testval, got %s", val)
	}
}
