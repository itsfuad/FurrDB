package db

import (
	"os"
	"strings"
	"testing"
	"time"
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

func TestTTLCommands(t *testing.T) {
	DefaultStore = NewStore()
	_, _ = setHandler([]string{"tk", "tv"})
	resp, _ := expireHandler([]string{"tk", "1"})
	if resp != "1" {
		t.Errorf("expected 1 from EXPIRE, got %s", resp)
	}

	ttl, _ := ttlHandler([]string{"tk"})
	if ttl != "1" && ttl != "0" { // allow for race
		t.Errorf("expected TTL 1 or 0, got %s", ttl)
	}
	time.Sleep(2 * time.Second)

	ttl, _ = ttlHandler([]string{"tk"})
	if ttl != "-2" {
		t.Errorf("expected TTL -2 after expiration, got %s", ttl)
	}
	val, _ := getHandler([]string{"tk"})
	if val != "" {
		t.Errorf("expected empty after expiration, got %s", val)
	}
}

func TestSnapshotSaveLoad(t *testing.T) {
	DefaultStore = NewStore()
	_, _ = setHandler([]string{"snapkey", "snapval"})
	_, _ = saddHandler([]string{"snapset", "a", "b"})
	_, _ = expireHandler([]string{"snapkey", "10"})

	err := SaveSnapshot("test_dump.rdb")
	if err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	// Clear DB
	DefaultStore = NewStore()
	val, _ := getHandler([]string{"snapkey"})
	if val != "" {
		t.Errorf("expected empty after clear, got %s", val)
	}

	err = LoadSnapshot("test_dump.rdb")
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}
	val, _ = getHandler([]string{"snapkey"})
	if val != "snapval" {
		t.Errorf("expected snapval after load, got %s", val)
	}
	members, _ := smembersHandler([]string{"snapset"})
	if !(strings.Contains(members, "a") && strings.Contains(members, "b")) {
		t.Errorf("expected set members a and b, got %s", members)
	}
	_ = os.Remove("test_dump.rdb")
}
