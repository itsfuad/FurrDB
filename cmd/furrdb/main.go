package main

import (
	"fmt"
	"os"

	_ "furr/internal/handlers"
	"furr/internal/repl"
	"furr/internal/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--repl" {
		repl.Start()
		return
	}
	fmt.Println("🦊 FurrDB starting on localhost:7070...")
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
