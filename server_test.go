package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestChatServer(t *testing.T) {
	// Start server in a goroutine
	go main()

	// Wait a bit for server to start
	time.Sleep(200 * time.Millisecond)

	// Connect client 1
	conn1, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	reader1 := bufio.NewReader(conn1)

	// Read welcome message
	msg, _ := reader1.ReadString('\n')
	if !strings.Contains(msg, "Bienvenido") {
		t.Errorf("Expected welcome message, got: %s", msg)
	}

	// Send nickname
	fmt.Fprintln(conn1, "Alice")

	// Read prompt and join message
	// Alice should receive "Ingresa tu nickname: " followed by "¡Alice se ha unido a la sala!"
	// But Wait, ReadString('\n') will skip the prompt if it doesn't have a newline.
	// The code says: fmt.Fprint(conn, "Ingresa tu nickname: ")

	msg, _ = reader1.ReadString('\n')
	if !strings.Contains(msg, "Alice se ha unido") {
		t.Errorf("Expected join message for Alice, got: %s", msg)
	}

	// Connect client 2
	conn2, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	reader2 := bufio.NewReader(conn2)

	// Read welcome message for client 2
	_, _ = reader2.ReadString('\n')

	// Send nickname for client 2
	fmt.Fprintln(conn2, "Bob")

	// Alice should see Bob joining
	msg, _ = reader1.ReadString('\n')
	if !strings.Contains(msg, "Bob se ha unido") {
		t.Errorf("Alice expected join message for Bob, got: %s", msg)
	}

	// Bob should see Bob joining (broadcasted to all)
	msg, _ = reader2.ReadString('\n')
	if !strings.Contains(msg, "Bob se ha unido") {
		t.Errorf("Bob expected join message for Bob, got: %s", msg)
	}

	// Alice sends a message
	fmt.Fprintln(conn1, "Hello from Alice")

	// Bob should receive Alice's message
	msg, _ = reader2.ReadString('\n')
	if !strings.Contains(msg, "[Alice]: Hello from Alice") {
		t.Errorf("Bob expected message from Alice, got: %s", msg)
	}

	// Alice should also receive her own message (it's broadcasted to all)
	msg, _ = reader1.ReadString('\n')
	if !strings.Contains(msg, "[Alice]: Hello from Alice") {
		t.Errorf("Alice expected her own message, got: %s", msg)
	}

	// Alice disconnects
	conn1.Close()

	// Bob should see Alice leaving
	msg, _ = reader2.ReadString('\n')
	if !strings.Contains(msg, "Alice ha salido") {
		t.Errorf("Bob expected leave message for Alice, got: %s", msg)
	}
}

func TestChatCommands(t *testing.T) {
	// Connect Alice
	connAlice, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect Alice: %v", err)
	}
	defer connAlice.Close()
	readerAlice := bufio.NewReader(connAlice)

	// Skip welcome
	readerAlice.ReadString('\n')
	// Send nickname
	fmt.Fprintln(connAlice, "Alice")
	// Alice joined
	readerAlice.ReadString('\n')

	// Connect Bob
	connBob, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect Bob: %v", err)
	}
	defer connBob.Close()
	readerBob := bufio.NewReader(connBob)

	// Skip welcome
	readerBob.ReadString('\n')
	// Send nickname
	fmt.Fprintln(connBob, "Bob")
	// Bob joined
	readerBob.ReadString('\n')
	// Alice sees Bob joined
	readerAlice.ReadString('\n')

	// Test /list
	fmt.Fprintln(connAlice, "/list")
	msg, _ := readerAlice.ReadString('\n')
	if !strings.Contains(msg, "Alice") || !strings.Contains(msg, "Bob") {
		t.Errorf("Expected Alice and Bob in list, got: %s", msg)
	}

	// Test /msg
	fmt.Fprintln(connAlice, "/msg Bob secret message")
	msg, _ = readerBob.ReadString('\n')
	if !strings.Contains(msg, "[Privado de Alice]: secret message") {
		t.Errorf("Bob expected private message, got: %s", msg)
	}

	// Test /msg user not found
	fmt.Fprintln(connAlice, "/msg Charlie hello")
	msg, _ = readerAlice.ReadString('\n')
	if !strings.Contains(msg, "El usuario no se encuentra en la sala") {
		t.Errorf("Alice expected user not found message, got: %s", msg)
	}
}
