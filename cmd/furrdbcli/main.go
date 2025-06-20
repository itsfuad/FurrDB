package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: furrdbcli <scriptfile> [host:port]")
		os.Exit(1)
	}
	scriptFile := os.Args[1]
	host := "localhost:7070"
	if len(os.Args) > 2 {
		host = os.Args[2]
	}

	file, err := os.Open(scriptFile)
	if err != nil {
		fmt.Println("Error opening script file:", err)
		os.Exit(1)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(file)
	serverReader := bufio.NewReader(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fmt.Fprintf(conn, "%s\n", line)
		resp, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Server closed connection.")
			return
		}
		fmt.Printf("> %s\n%s", line, resp)
	}
}
