package handlers

import (
	"fmt"
	"strings"

	"furr/internal/db"
	"furr/internal/script"
)

func regscriptHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for REGSCRIPT")
	}
	scriptStr := strings.Join(args, " ")
	hash := script.RegisterScript(scriptStr)
	return hash, nil
}

func runscriptHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for RUNSCRIPT")
	}
	hash := args[0]
	return script.RunScript(hash, args[1:])
}

func evalHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("missing argument for EVAL")
	}
	scriptStr := strings.Join(args, " ")
	return script.EvalScript(scriptStr, nil)
}

func init() {
	db.Commands["REGSCRIPT"] = regscriptHandler
	db.Commands["RUNSCRIPT"] = runscriptHandler
	db.Commands["EVAL"] = evalHandler
}
