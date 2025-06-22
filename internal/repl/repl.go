package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"furr/internal/db"
)

func Start() {
	fmt.Println("🦊 FurrDB REPL (type HELP for commands, EXIT to quit)")
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		cmd := strings.ToUpper(tokens[0])
		args := tokens[1:]
		switch cmd {
		case "EXIT":
			fmt.Println("bye!")
			return
		case "HELP":
			printHelp()
		case "CLEAR":
			clearScreen()
		default:
			handler, ok := db.Commands[cmd]
			if !ok {
				fmt.Println("ERR unknown command")
				continue
			}
			result, err := handler(args)
			if err != nil {
				fmt.Println("ERR", err)
				continue
			}
			fmt.Println(result)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[2J\033[H") // ANSI escape code
}

func printHelp() {
	fmt.Println(`Available commands:
	SET key value      - Set key to value
	GET key            - Get value of key
	DEL key            - Delete key
	EXISTS key         - Check if key exists
	LPUSH k v [v..]    - Push value(s) to head of list
	RPUSH k v [v..]    - Push value(s) to tail of list
	LPOP k             - Pop value from head of list
	RPOP k             - Pop value from tail of list
	LRANGE k s e       - Get list elements from s to e
	SADD k v [v..]     - Add value(s) to set
	SREM k v [v..]     - Remove value(s) from set
	SMEMBERS k         - List all set members
	KEYS               - List all keys
	FLUSHDB            - Clear the database
	INFO               - Show server info
	PING               - Responds with PONG
	REGSCRIPT script   - Register script, returns hash
	RUNSCRIPT hash     - Run registered script by hash
	EVAL script        - Evaluate script string
	SAVE               - Force persistence flush
	CLEAR              - Clear the screen
	EXIT               - Exit the REPL
	HELP               - Show this help menu
	`)
}
