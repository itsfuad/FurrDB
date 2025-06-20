package script

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"furr/internal/db"
	"strings"
)

var scripts = make(map[string]string) // hash -> script

// RegisterScript stores a script and returns its hash
func RegisterScript(script string) string {
	h := sha256.Sum256([]byte(script))
	hash := hex.EncodeToString(h[:])
	scripts[hash] = script
	return hash
}

// RunScript executes a registered script by hash (stub)
func RunScript(hash string, args []string) (string, error) {
	script, ok := scripts[hash]
	if !ok {
		return "", nil
	}
	return evalScriptLines(script, args)
}

// EvalScript evaluates a script string without storing (stub)
func EvalScript(script string, args []string) (string, error) {
	return evalScriptLines(script, args)
}

func evalScriptLines(script string, args []string) (string, error) {
	lines := strings.Split(script, ";")
	var last string
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		cmd := strings.ToUpper(tokens[0])
		params := tokens[1:]
		handler, ok := db.Commands[cmd]
		if !ok {
			return "", fmt.Errorf("ERR unknown command '%s' on line %d", cmd, i+1)
		}
		// Check for required arguments (for SET, GET, DEL, EXISTS)
		switch cmd {
		case "SET":
			if len(params) < 2 {
				return "", fmt.Errorf("ERR missing argument for SET on line %d", i+1)
			}
		case "GET", "DEL", "EXISTS":
			if len(params) < 1 {
				return "", fmt.Errorf("ERR missing argument for %s on line %d", cmd, i+1)
			}
		}
		result, err := handler(params)
		if err != nil {
			return "", fmt.Errorf("ERR %v on line %d", err, i+1)
		}
		last = result
	}
	return last, nil
}
