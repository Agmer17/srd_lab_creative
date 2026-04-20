package ws

import (
	"log"

	"github.com/olahol/melody"
)

type Room struct {
	Id      string
	Clients map[*melody.Session]bool
}

func NewRoom(id string) *Room {
	return &Room{
		Id:      id,
		Clients: make(map[*melody.Session]bool),
	}
}

func (room *Room) RegisterConnection(c *melody.Session) {
	room.Clients[c] = true
}

func (room *Room) UnregisterClients(c *melody.Session) {
	if userID, ok := c.Get("userID"); ok {
		log.Printf("User %s disconnected", userID)
	}
	delete(room.Clients, c)
}

func (room *Room) GetRoomClientsCount() int {
	return len(room.Clients)
}
