# Go-Chat: Servidor de mensajería TCP en tiempo real

Go-Chat es un servidor de mensajería TCP concurrente escrito en Go, diseñado para ofrecer una comunicación fluida y segura en tiempo real.

## 🚀 Características

*   **Concurrencia nativa con Goroutines**: Maneja múltiples conexiones simultáneas de forma eficiente.
*   **Broadcasting de mensajes**: Los mensajes enviados se distribuyen automáticamente a todos los usuarios conectados.
*   **Protección de memoria con Mutex**: Garantiza la integridad de los datos en entornos concurrentes.
*   **Soporte para comandos y mensajes privados**: Funcionalidades extendidas para una mejor interacción.

## 🛠️ Instalación y Uso

### Levantar el servidor

Para iniciar el servidor, asegúrate de tener Go instalado y ejecuta el siguiente comando en la raíz del proyecto:

```bash
go run main.go
```

Por defecto, el servidor escuchará en el puerto **8080**.

### Conectarse como cliente

Puedes conectarte al servidor utilizando herramientas estándar de terminal:

*   **Linux/Mac (usando netcat):**
    ```bash
    nc localhost 8080
    ```
*   **Windows (usando telnet):**
    ```bash
    telnet localhost 8080
    ```

## 💬 Comandos Disponibles

Una vez conectado, puedes interactuar con el chat a través de los siguientes comandos:

*   `/list`: Muestra una lista de todos los usuarios actualmente conectados.
*   `/msg [usuario] [mensaje]`: Envía un mensaje directo y privado a un usuario específico.

---
*Desarrollado con ❤️ en Go.*
