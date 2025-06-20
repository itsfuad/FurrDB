package script

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"furr/internal/db"
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
	vars := make(map[string]string)
	var last string
	maxLines := 100
	if len(lines) > maxLines {
		return "", fmt.Errorf("ERR script too long (max %d lines)", maxLines)
	}
	whitelist := map[string]bool{
		"SET": true, "GET": true, "DEL": true, "EXISTS": true,
		"LPUSH": true, "RPUSH": true, "LPOP": true, "RPOP": true, "LRANGE": true,
		"SADD": true, "SREM": true, "SMEMBERS": true,
	}
	var skipToEnd int // 0=not skipping, >0=inside nested IFs
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if skipToEnd > 0 {
			if strings.HasPrefix(line, "IF ") {
				skipToEnd++
			} else if line == "END" {
				skipToEnd--
			}
			continue
		}
		if strings.HasPrefix(line, "LET ") {
			// LET x = CMD args
			parts := strings.Fields(line)
			if len(parts) < 5 || parts[2] != "=" {
				return "", fmt.Errorf("ERR invalid LET syntax on line %d", i+1)
			}
			varName := parts[1]
			cmd := strings.ToUpper(parts[3])
			cmdArgs := parts[4:]
			if !whitelist[cmd] {
				return "", fmt.Errorf("ERR command %s not allowed in LET on line %d", cmd, i+1)
			}
			handler, ok := db.Commands[cmd]
			if !ok {
				return "", fmt.Errorf("ERR unknown command '%s' in LET on line %d", cmd, i+1)
			}
			result, err := handler(cmdArgs)
			if err != nil {
				return "", fmt.Errorf("ERR %v on line %d", err, i+1)
			}
			vars[varName] = result
			last = result
			continue
		}
		if strings.HasPrefix(line, "IF ") {
			// IF var == value
			parts := strings.Fields(line)
			if len(parts) != 4 || parts[2] != "==" {
				return "", fmt.Errorf("ERR invalid IF syntax on line %d", i+1)
			}
			varName := parts[1]
			expected := parts[3]
			if vars[varName] != expected {
				skipToEnd = 1
			}
			continue
		}
		if line == "END" {
			continue
		}
		// Normal DB command
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		cmd := strings.ToUpper(tokens[0])
		params := tokens[1:]
		if !whitelist[cmd] {
			return "", fmt.Errorf("ERR command %s not allowed in script on line %d", cmd, i+1)
		}
		handler, ok := db.Commands[cmd]
		if !ok {
			return "", fmt.Errorf("ERR unknown command '%s' on line %d", cmd, i+1)
		}
		result, err := handler(params)
		if err != nil {
			return "", fmt.Errorf("ERR %v on line %d", err, i+1)
		}
		last = result
	}
	return last, nil
}
