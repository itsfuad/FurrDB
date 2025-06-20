package main

import (
	"fmt"
	"os"

	"furr/internal/server"
)

func main() {
	fmt.Println("🦊 FurrDB starting on localhost:7070...")
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
