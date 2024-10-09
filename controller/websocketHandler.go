package controller

import (
	"fmt"
	"github/similadayo/chitchat/utils"
	"net/http"
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := utils.UpgradeConnection(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &utils.Client{
		ID:   "1",
		Conn: conn,
		Send: make(chan []byte),
	}

	go utils.HandleMessages(client)
	go utils.WriteMessage(client)
}
