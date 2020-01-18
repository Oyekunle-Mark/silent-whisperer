package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	socket *websocket.Conn
	send chan []byte
	room *room
}
