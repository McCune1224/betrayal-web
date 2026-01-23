package game

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
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

// Client represents a single WebSocket connection from a browser.
// Each client belongs to exactly one room and has its own goroutines
// for reading from and writing to the WebSocket connection.
type Client struct {
	// RoomCode identifies which game room this client is connected to.
	// Used by the Hub to route broadcasts to the correct clients.
	RoomCode string

	// PlayerID is the unique identifier for this player in the game.
	// This links the WebSocket connection to the player's game state.
	PlayerID string

	// Send is a buffered channel for outgoing messages.
	// The writePump goroutine reads from this channel and writes to the WebSocket.
	// Using a channel decouples message production from network I/O.
	Send chan Message

	// Conn is the actual WebSocket connection to the browser.
	// Only the readPump and writePump goroutines should access this directly.
	// It's exported so the handlers package can set it during upgrade.
	Conn *websocket.Conn

	// Hub is a reference back to the hub for unregistration.
	// When the client disconnects, it sends itself to hub.unregister.
	Hub *Hub
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

// Register adds a client to the hub's register channel.
// This is the public API for the handlers package to register new connections.
// The actual registration happens asynchronously in the Run() goroutine.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister adds a client to the hub's unregister channel.
// This is the public API for disconnecting clients.
// The actual unregistration happens asynchronously in the Run() goroutine.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
