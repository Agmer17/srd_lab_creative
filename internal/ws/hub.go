package ws

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type WebsocketHub struct {
	mel   *melody.Melody
	mu    sync.RWMutex
	Rooms map[string]*Room

	SessionRooms map[*melody.Session]map[string]bool
}

func NewWebsocketHub(m *melody.Melody) *WebsocketHub {
	h := &WebsocketHub{
		mel:          m,
		Rooms:        make(map[string]*Room),
		SessionRooms: make(map[*melody.Session]map[string]bool),
	}

	h.hubInit()
	return h

}

func (hub *WebsocketHub) hubInit() {
	hub.mel.HandleConnect(hub.HandleConnectionRequest)
	hub.mel.HandleDisconnect(hub.RemoveFromAllRooms)
}

func (hub *WebsocketHub) GetOrCreateRoom(id string) *Room {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	r, ok := hub.Rooms[id]
	if !ok {
		newRoom := NewRoom(id)
		hub.Rooms[id] = newRoom
		return newRoom
	}
	return r
}

func (hub *WebsocketHub) RemoveFromAllRooms(s *melody.Session) {

	hub.mu.Lock()
	defer hub.mu.Unlock()

	joinedRooms := hub.SessionRooms[s]

	for id := range joinedRooms {
		if r, ok := hub.Rooms[id]; ok {
			r.UnregisterClients(s)
			if r.GetRoomClientsCount() == 0 {
				delete(hub.Rooms, r.Id)
				fmt.Println("room ", id, "deleted")
			}
		}
	}

	delete(hub.SessionRooms, s)
}

func (hub *WebsocketHub) HandleConnectionRequest(s *melody.Session) {

	idStr := s.Request.Header.Get("X-User-ID")

	userId, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("id str : ", idStr)
		fmt.Println("id str : ", idStr)

		s.Close()
		return
	}

	s.Set("user_id", userId)
	room := hub.GetOrCreateRoom(userId.String())
	room.RegisterConnection(s)

	if hub.SessionRooms[s] == nil {
		hub.SessionRooms[s] = make(map[string]bool)
	}

	fmt.Println("room that user join : ", room.Id)
	fmt.Println("user id : ", userId)

	hub.SessionRooms[s][userId.String()] = true
}
