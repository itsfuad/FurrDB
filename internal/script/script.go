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
	return evalScriptLines(script)
}

// EvalScript evaluates a script string without storing (stub)
func EvalScript(script string) (string, error) {
	return evalScriptLines(script)
}

// getWhitelist returns the allowed commands
func getWhitelist() map[string]bool {
	return map[string]bool{
		"SET": true, "GET": true, "DEL": true, "EXISTS": true,
		"LPUSH": true, "RPUSH": true, "LPOP": true, "RPOP": true, "LRANGE": true,
		"SADD": true, "SREM": true, "SMEMBERS": true,
	}
}

// handleSkipping processes lines when in a skipping state
func handleSkipping(line string, skipToEnd int) int {
	if strings.HasPrefix(line, "IF ") {
		return skipToEnd + 1
	} else if line == "END" {
		return skipToEnd - 1
	}
	return skipToEnd
}

// handleLetStatement processes a LET statement
func handleLetStatement(line string, lineNum int, whitelist map[string]bool, vars map[string]string) (string, error) {
	parts := strings.Fields(line)
	if len(parts) < 5 || parts[2] != "=" {
		return "", fmt.Errorf("ERR invalid LET syntax on line %d", lineNum)
	}

	varName := parts[1]
	cmd := strings.ToUpper(parts[3])
	cmdArgs := parts[4:]

	if !whitelist[cmd] {
		return "", fmt.Errorf("ERR command %s not allowed in LET on line %d", cmd, lineNum)
	}

	handler, ok := db.Commands[cmd]
	if !ok {
		return "", fmt.Errorf("ERR unknown command '%s' in LET on line %d", cmd, lineNum)
	}

	result, err := handler(cmdArgs)
	if err != nil {
		return "", fmt.Errorf("ERR %v on line %d", err, lineNum)
	}

	vars[varName] = result
	return result, nil
}

// handleIfStatement processes an IF statement
func handleIfStatement(line string, lineNum int, vars map[string]string) (int, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 || parts[2] != "==" {
		return 0, fmt.Errorf("ERR invalid IF syntax on line %d", lineNum)
	}

	varName := parts[1]
	expected := parts[3]

	if vars[varName] != expected {
		return 1, nil
	}
	return 0, nil
}

// executeCommand executes a normal DB command
func executeCommand(line string, lineNum int, whitelist map[string]bool) (string, error) {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return "", nil
	}

	cmd := strings.ToUpper(tokens[0])
	params := tokens[1:]

	if !whitelist[cmd] {
		return "", fmt.Errorf("ERR command %s not allowed in script on line %d", cmd, lineNum)
	}

	handler, ok := db.Commands[cmd]
	if !ok {
		return "", fmt.Errorf("ERR unknown command '%s' on line %d", cmd, lineNum)
	}

	result, err := handler(params)
	if err != nil {
		return "", fmt.Errorf("ERR %v on line %d", err, lineNum)
	}

	return result, nil
}

func evalScriptLines(script string) (string, error) {
	lines := strings.Split(script, ";")
	vars := make(map[string]string)
	var last string

	const maxLines = 100
	if len(lines) > maxLines {
		return "", fmt.Errorf("ERR script too long (max %d lines)", maxLines)
	}

	whitelist := getWhitelist()
	var skipToEnd int

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if skipToEnd > 0 {
			skipToEnd = handleSkipping(line, skipToEnd)
			continue
		}
		result, skip, err := evalScriptLine(line, i+1, whitelist, vars)
		if err != nil {
			return "", err
		}
		if skip > 0 {
			skipToEnd = skip
			continue
		}
		if result != "" {
			last = result
		}
	}
	return last, nil
}

func evalScriptLine(line string, lineNum int, whitelist map[string]bool, vars map[string]string) (result string, skip int, err error) {
	switch {
	case strings.HasPrefix(line, "LET "):
		result, err = handleLetStatement(line, lineNum, whitelist, vars)
	case strings.HasPrefix(line, "IF "):
		skip, err = handleIfStatement(line, lineNum, vars)
	case line == "END":
		// nothing to do
	default:
		result, err = executeCommand(line, lineNum, whitelist)
	}
	return
}
