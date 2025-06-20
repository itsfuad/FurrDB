package engine

import (
	"bufio"
	"furr/internal/db"
	"os"
	"strings"
)

var aofPath = "aof.log"
var aofFile *os.File

// Append appends a command to the AOF log (stub)
func Append(cmd string) error {
	if aofFile == nil {
		f, err := os.OpenFile(aofPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		aofFile = f
	}
	_, err := aofFile.WriteString(cmd + "\n")
	return err
}

// Load loads the AOF log and replays commands (stub)
func Load() error {
	f, err := os.Open(aofPath)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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
		if ok {
			_, _ = handler(params)
		}
	}
	return scanner.Err()
}

// Flush forces a flush to disk (stub)
func Flush() error {
	if aofFile != nil {
		return aofFile.Sync()
	}
	return nil
}
