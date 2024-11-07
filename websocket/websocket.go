package websocket

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]string)// Connected clients
var broadcast = make(chan []byte)            // Broadcast channel

func HandleConnections(c *websocket.Conn) {
    defer func() {
        c.Close()
        delete(clients, c)
    }()

    // Wait for the client to send the initial message with the userId
    var userId string
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, c)
            break
        }

        // Assume the first message from the client is the userId
        var initMsg map[string]string
        if err := json.Unmarshal(message, &initMsg); err != nil {
            log.Printf("error unmarshaling init message: %v", err)
            continue
        }

        if id, ok := initMsg["userId"]; ok {
            userId = id
            clients[c] = userId
            log.Printf("Registered userId %s with connection", userId)
            break
        }
    }

    // Listen for further messages (optional)
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