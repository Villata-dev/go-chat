package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
	defer listener.Close()

	fmt.Println("Servidor Go-Chat iniciado en el puerto 8080...")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}

		// Handle connection in a new Goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Greet the user
	fmt.Fprintln(conn, "Bienvenido al servidor Go-Chat!")

	reader := bufio.NewReader(conn)
	for {
		// Read message from user
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Conexión cerrada por el usuario o error: %v\n", err)
			return
		}

		// Remove newline and trim space
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Echo the message back to the same user
		fmt.Fprintf(conn, "[Eco]: %s\n", message)
	}
}
