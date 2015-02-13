package webjack

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

var maxClientId int = 0

type Client struct {
	id int
	ws *websocket.Conn
}

func NewClient(ws *websocket.Conn) *Client {
	if ws == nil {
		log.Fatal("websocket.Conn cannot be nil")
	}

	maxClientId++
	return &Client{maxClientId, ws}
}

func (self *Client) Send(msg interface{}) error {
	return websocket.JSON.Send(self.ws, &msg)
}

func (self *Client) Receive(msg interface{}) error {
	return websocket.JSON.Receive(self.ws, msg)
}
