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

		// Intercept commands
		if strings.HasPrefix(message, "/") {
			handleCommand(conn, message, nickname)
			continue
		}

		// Format and broadcast message
		broadcast <- fmt.Sprintf("[%s]: %s", nickname, message)
	}
}

func handleCommand(conn net.Conn, message, nickname string) {
	parts := strings.SplitN(message, " ", 3)
	command := parts[0]

	switch command {
	case "/list":
		mutex.Lock()
		var names []string
		for _, name := range clients {
			names = append(names, name)
		}
		mutex.Unlock()
		fmt.Fprintf(conn, "Usuarios conectados: %s\n", strings.Join(names, ", "))
	case "/msg":
		if len(parts) < 3 {
			fmt.Fprintln(conn, "Uso: /msg [NombreUsuario] [Mensaje]")
			return
		}
		targetNickname := parts[1]
		secretMessage := parts[2]

		var targetConn net.Conn
		mutex.Lock()
		for c, name := range clients {
			if name == targetNickname {
				targetConn = c
				break
			}
		}
		mutex.Unlock()

		if targetConn != nil {
			fmt.Fprintf(targetConn, "[Privado de %s]: %s\n", nickname, secretMessage)
		} else {
			fmt.Fprintln(conn, "El usuario no se encuentra en la sala")
		}
	default:
		fmt.Fprintln(conn, "Comando desconocido. Usa /list o /msg")
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
