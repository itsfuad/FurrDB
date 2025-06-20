package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func sendCommand(conn net.Conn, cmd string) string {
	fmt.Fprintf(conn, "%s\n", cmd)
	reader := bufio.NewReader(conn)
	resp, _ := reader.ReadString('\n')
	return strings.TrimSpace(resp)
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7070")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to FurrDB")

	// CREATE
	fmt.Println("[CREATE] Set user:1 name and email")
	fmt.Println(sendCommand(conn, "SET user:1:name Alice"))
	fmt.Println(sendCommand(conn, "SET user:1:email alice@example.com"))

	// READ
	fmt.Println("\n[READ] Get user:1 name and email")
	fmt.Println("Name:", sendCommand(conn, "GET user:1:name"))
	fmt.Println("Email:", sendCommand(conn, "GET user:1:email"))

	// UPDATE
	fmt.Println("\n[UPDATE] Update user:1 name")
	fmt.Println(sendCommand(conn, "SET user:1:name Alicia"))
	fmt.Println("Updated Name:", sendCommand(conn, "GET user:1:name"))

	// EXISTS
	fmt.Println("\n[EXISTS] Check if user:1:email exists")
	fmt.Println("Exists:", sendCommand(conn, "EXISTS user:1:email"))

	// DELETE
	fmt.Println("\n[DELETE] Delete user:1:email")
	fmt.Println(sendCommand(conn, "DEL user:1:email"))
	fmt.Println("Email after delete:", sendCommand(conn, "GET user:1:email"))

	// KEYS
	fmt.Println("\n[KEYS] List all keys")
	fmt.Println(sendCommand(conn, "KEYS"))

	// EXIT
	sendCommand(conn, "EXIT")
	fmt.Println("\nConnection closed")
}
