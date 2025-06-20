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
		var resp string
		switch cmd {
		case "PING":
			resp = "PONG"
		case "EXIT":
			w.WriteString("BYE\n")
			w.Flush()
			return
		default:
			handler, ok := db.Commands[cmd]
			if !ok {
				resp = "ERR unknown command"
			} else {
				result, err := handler(args)
				if err != nil {
					resp = "ERR " + err.Error()
				} else {
					resp = result
				}
			}
		}
		w.WriteString(resp + "\n")
		w.Flush()
	}
}
