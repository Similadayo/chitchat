package utils

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

// UpgradeConnection upgrades the HTTP connection to a WebSocket connection
func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// HandleMessages handles incoming messages from the WebSocket connection
func HandleMessages(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}

		client.Send <- msg
	}
}

// WriteMessage writes a message to the WebSocket connection
func WriteMessage(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		msg, ok := <-client.Send
		if !ok {
			return
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error writing message: ", err)
			break
		}
	}
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn, id string) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte),
	}
}
