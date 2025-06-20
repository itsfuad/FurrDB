package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"furr/internal/db"
)

func Start() error {
	ln, err := net.Listen("tcp", "localhost:7070")
	if err != nil {
		return err
	}
	defer ln.Close()
	fmt.Println("[server] Listening on localhost:7070")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[server] Accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		
		tokens := parseInput(line)
		if len(tokens) == 0 {
			continue
		}
		
		cmd := strings.ToUpper(tokens[0])
		args := tokens[1:]
		
		if cmd == "EXIT" {
			w.WriteString("BYE\n")
			w.Flush()
			return
		}
		
		resp := processCommand(cmd, args)
		w.WriteString(resp + "\n")
		w.Flush()
	}
}

func parseInput(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	return strings.Fields(line)
}

func processCommand(cmd string, args []string) string {
	if cmd == "PING" {
		return "PONG"
	}
	
	handler, ok := db.Commands[cmd]
	if !ok {
		return "ERR unknown command"
	}
	
	result, err := handler(args)
	if err != nil {
		return "ERR " + err.Error()
	}
	return result
}
