package db

import (
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Command handler type
type HandlerFunc func(args []string) (string, error)

// Command dispatch map
var Commands = map[string]HandlerFunc{
	"SET":    setHandler,
	"GET":    getHandler,
	"DEL":    delHandler,
	"EXISTS": existsHandler,
}

var DefaultStore = NewStore()

func setHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}
	key, value := args[0], args[1]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	DefaultStore.data[key] = value
	return "OK", nil
}

func getHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", nil
	}
	key := args[0]
	DefaultStore.mu.RLock()
	defer DefaultStore.mu.RUnlock()
	val, ok := DefaultStore.data[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func delHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "0", nil
	}
	key := args[0]
	DefaultStore.mu.Lock()
	defer DefaultStore.mu.Unlock()
	_, existed := DefaultStore.data[key]
	delete(DefaultStore.data, key)
	if existed {
		return "1", nil
	}
	return "0", nil
}

func existsHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "0", nil
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
