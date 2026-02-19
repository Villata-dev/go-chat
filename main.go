package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mutex     sync.Mutex
)

func main() {
	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
	defer listener.Close()

	fmt.Println("Servidor Go-Chat iniciado en el puerto 8080...")

	go broadcaster()

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

	// Ask for nickname
	fmt.Fprint(conn, "Ingresa tu nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error al leer el nickname: %v\n", err)
		return
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = "Anónimo"
	}

	// Register client
	mutex.Lock()
	clients[conn] = nickname
	mutex.Unlock()

	// Announce arrival
	broadcast <- fmt.Sprintf("¡%s se ha unido a la sala!", nickname)

	for {
		// Read message from user
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Conexión cerrada por el usuario o error: %v\n", err)

			// Unregister client
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()

			// Announce departure
			broadcast <- fmt.Sprintf("¡%s ha salido de la sala!", nickname)
			return
		}

		// Remove newline and trim space
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// Format and broadcast message
		broadcast <- fmt.Sprintf("[%s]: %s", nickname, message)
	}
}

func broadcaster() {
	for msg := range broadcast {
		mutex.Lock()
		for client := range clients {
			fmt.Fprintln(client, msg)
		}
		mutex.Unlock()
	}
}
