package websocket

import (
    "log"
    "github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]string)// Connected clients
var broadcast = make(chan []byte)            // Broadcast channel

// HandleConnections handles incoming WebSocket connections
func HandleConnections(c *websocket.Conn, userId string) {
    defer func() {
        c.Close()
        delete(clients, c)
    }()

    // Register the new client
    clients[c] = userId

    // Listen for messages (if needed)
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, c)
            break
        }
        log.Printf("Received from %s: %s", userId, message)
    }
}

// SendUpdate broadcasts schedule updates to all clients
func SendUpdate(updateMessage []byte, targetUserId string) {
    for client, userId := range clients {
        if userId == targetUserId { // Only send to the specific user
            err := client.WriteMessage(websocket.TextMessage, updateMessage)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}