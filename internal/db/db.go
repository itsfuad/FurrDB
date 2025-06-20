package db

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type valueType int

const (
	StringType valueType = iota
	ListType
	SetType
)

type Store struct {
	mu    sync.RWMutex
	data  map[string]any
	types map[string]valueType
}

func NewStore() *Store {
	return &Store{
		data:  make(map[string]any),
		types: make(map[string]valueType),
	}
}

var DefaultStore = NewStore()

type HandlerFunc func(args []string) (string, error)

var Commands = map[string]HandlerFunc{
	"SET":      setHandler,
	"GET":      getHandler,
	"DEL":      delHandler,
	"EXISTS":   existsHandler,
	"LPUSH":    lpushHandler,
	"RPUSH":    rpushHandler,
	"LPOP":     lpopHandler,
	"RPOP":     rpopHandler,
	"LRANGE":   lrangeHandler,
	"SADD":     saddHandler,
	"SREM":     sremHandler,
	"SMEMBERS": smembersHandler,
	"KEYS":     keysHandler,
	"FLUSHDB":  flushdbHandler,
	"INFO":     infoHandler,
}

// String commands
func setHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("missing argument for SET")
	}
	key, value := args[0], args[1]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	DefaultStore.data[key] = value
	DefaultStore.types[key] = StringType
	return "OK", nil
}

func getHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for GET")
	}
	key := args[0]
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	if DefaultStore.types[key] != StringType {
		return "", nil
	}
	val, ok := DefaultStore.data[key].(string)
	if !ok {
		return "", nil
	}
	return val, nil
}

// List commands
func lpushHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("missing argument for LPUSH")
	}
	key := args[0]
	vals := args[1:]
	// Reverse vals for Redis-like LPUSH
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != ListType {
		DefaultStore.data[key] = []string{}
		DefaultStore.types[key] = ListType
	}
	lst := DefaultStore.data[key].([]string)
	lst = append(vals, lst...)
	DefaultStore.data[key] = lst
	return fmt.Sprintf("%d", len(lst)), nil
}

func rpushHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("missing argument for RPUSH")
	}
	key := args[0]
	vals := args[1:]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != ListType {
		DefaultStore.data[key] = []string{}
		DefaultStore.types[key] = ListType
	}
	lst := DefaultStore.data[key].([]string)
	lst = append(lst, vals...)
	DefaultStore.data[key] = lst
	return fmt.Sprintf("%d", len(lst)), nil
}

func lpopHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for LPOP")
	}
	key := args[0]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != ListType {
		return "", nil
	}
	lst := DefaultStore.data[key].([]string)
	if len(lst) == 0 {
		return "", nil
	}
	val := lst[0]
	DefaultStore.data[key] = lst[1:]
	return val, nil
}

func rpopHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for RPOP")
	}
	key := args[0]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != ListType {
		return "", nil
	}
	lst := DefaultStore.data[key].([]string)
	if len(lst) == 0 {
		return "", nil
	}
	val := lst[len(lst)-1]
	DefaultStore.data[key] = lst[:len(lst)-1]
	return val, nil
}

func lrangeHandler(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("missing argument for LRANGE")
	}
	key := args[0]
	start := parseInt(args[1])
	end := parseInt(args[2])
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	if DefaultStore.types[key] != ListType {
		return "", nil
	}
	lst := DefaultStore.data[key].([]string)
	if start < 0 {
		start = 0
	}
	if end >= len(lst) {
		end = len(lst) - 1
	}
	if start > end || start >= len(lst) {
		return "", nil
	}
	return strings.Join(lst[start:end+1], ","), nil
}

// Set commands
func saddHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("missing argument for SADD")
	}
	key := args[0]
	vals := args[1:]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != SetType {
		DefaultStore.data[key] = map[string]struct{}{}
		DefaultStore.types[key] = SetType
	}
	set := DefaultStore.data[key].(map[string]struct{})
	added := 0
	for _, v := range vals {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			added++
		}
	}
	return fmt.Sprintf("%d", added), nil
}

func sremHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("missing argument for SREM")
	}
	key := args[0]
	vals := args[1:]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	if DefaultStore.types[key] != SetType {
		return "0", nil
	}
	set := DefaultStore.data[key].(map[string]struct{})
	removed := 0
	for _, v := range vals {
		if _, exists := set[v]; exists {
			delete(set, v)
			removed++
		}
	}
	return fmt.Sprintf("%d", removed), nil
}

func smembersHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for SMEMBERS")
	}
	key := args[0]
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	if DefaultStore.types[key] != SetType {
		return "", nil
	}
	set := DefaultStore.data[key].(map[string]struct{})
	members := make([]string, 0, len(set))
	for v := range set {
		members = append(members, v)
	}
	sort.Strings(members)
	return strings.Join(members, ","), nil
}

// Meta commands
func keysHandler(args []string) (string, error) {
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	keys := make([]string, 0, len(DefaultStore.data))
	for k := range DefaultStore.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ","), nil
}

func flushdbHandler(args []string) (string, error) {
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	DefaultStore.data = make(map[string]any)
	DefaultStore.types = make(map[string]valueType)
	return "OK", nil
}

func infoHandler(args []string) (string, error) {
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	return fmt.Sprintf("keys:%d", len(DefaultStore.data)), nil
}

// Utility
func parseInt(s string) int {
	n, _ := fmt.Sscanf(s, "%d", new(int))
	if n == 0 {
		return 0
	}
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// Existing DEL, EXISTS
func delHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "0", fmt.Errorf("missing argument for DEL")
	}
	key := args[0]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	_, existed := DefaultStore.data[key]
	delete(DefaultStore.data, key)
	delete(DefaultStore.types, key)
	if existed {
		return "1", nil
	}
	return "0", nil
}

func existsHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "0", fmt.Errorf("missing argument for EXISTS")
	}
	key := args[0]
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	_, ok := DefaultStore.data[key]
	if ok {
		return "1", nil
	}
	return "0", nil
}
