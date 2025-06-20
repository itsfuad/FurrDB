package main

import (
	"fmt"
	"os"

	"furr/internal/db"
	_ "furr/internal/handlers"
	"furr/internal/repl"
	"furr/internal/server"
)

func main() {
	if _, err := os.Stat("dump.rdb"); err == nil {
		db.LoadSnapshot("dump.rdb")
	}
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
