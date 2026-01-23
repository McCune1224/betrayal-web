package game

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	mu          sync.RWMutex
	rooms       map[string]*RoomConnections
	register    chan *Client
	unregister  chan *Client
	roomManager *RoomManager
}

type RoomConnections struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

type Client struct {
	RoomCode string
	PlayerID string
	Send     chan Message
	conn     interface{}
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]*RoomConnections),
		register:    make(chan *Client, 256),
		unregister:  make(chan *Client, 256),
		roomManager: NewRoomManager(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.rooms[client.RoomCode]; !exists {
		h.rooms[client.RoomCode] = &RoomConnections{
			clients: make(map[*Client]bool),
		}
	}
	h.rooms[client.RoomCode].clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomConns := h.rooms[client.RoomCode]
	if roomConns != nil {
		delete(roomConns.clients, client)
		close(client.Send)

		if len(roomConns.clients) == 0 {
			delete(h.rooms, client.RoomCode)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomCode string, msg Message) {
	h.mu.RLock()
	roomConns := h.rooms[roomCode]
	h.mu.RUnlock()

	if roomConns == nil {
		return
	}

	roomConns.mu.RLock()
	defer roomConns.mu.RUnlock()

	for client := range roomConns.clients {
		select {
		case client.Send <- msg:
		default:
			// Client send queue full, skip
		}
	}
}

func (h *Hub) GetRoomManager() *RoomManager {
	return h.roomManager
}
