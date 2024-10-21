package websocket

import (
    "log"
    "github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel

// HandleConnections handles incoming WebSocket connections
func HandleConnections(c *websocket.Conn) {
    defer func() {
        c.Close()
        delete(clients, c)
    }()

    // Register the new client
    clients[c] = true

    // Listen for messages (if needed)
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, c)
            break
        }

        // Optionally handle incoming messages from the client
        log.Printf("Received: %s", message)
    }
}

// HandleMessages broadcasts messages to all clients
func HandleMessages() {
    for {
        // Grab the next message from the broadcast channel
        message := <-broadcast

        // Send the message to all connected clients
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, message)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

// SendUpdate broadcasts schedule updates to all clients
func SendUpdate(updateMessage []byte) {
    broadcast <- updateMessage
}