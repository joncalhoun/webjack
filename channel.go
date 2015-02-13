package webjack

import (
	"log"
)

type Channel struct {
	clients map[int]*Client
}

func NewChannel() *Channel {
	return &Channel{
		make(map[int]*Client),
	}
}

func (self *Channel) RemoveClient(c *Client) {
	delete(self.clients, c.id)
}

func (self *Channel) AddClient(c *Client) {
	self.clients[c.id] = c
}

func (self *Channel) SendAll(msg interface{}) {
	log.Printf("Send all with %+v\n", msg)
	for _, c := range self.clients {
		err := c.Send(msg)
		if err != nil {
			log.Printf("Error sending to client #%d: %s\n", c.id, err)
		}
	}
}
