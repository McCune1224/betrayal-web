package game

import "backend/internal/logging"

// Hub manages all WebSocket clients and coordinates message broadcasting.
// WHY: Centralizes all WebSocket state to prevent race conditions.
// The hub runs in its own goroutine and uses channels for all communication.
type Hub struct {
	// rooms maps room codes to sets of clients in that room.
	// WHY SET: O(1) add/remove/lookup. Using map[*Client]bool as a set.
	rooms map[string]map[*Client]bool

	// broadcast sends messages to all clients in all rooms.
	// WHY: Used for system-wide announcements or testing.
	broadcast chan WSMessage

	// register accepts new clients connecting to the hub.
	// CHANNEL PATTERN: Fan-in (many clients -> one hub).
	// Multiple client goroutines send to one hub goroutine.
	register chan *Client

	// unregister accepts clients disconnecting from the hub.
	// WHY SEPARATE CHANNEL: Allows clean cleanup sequence.
	unregister chan *Client

	// roomManager handles game logic and room state.
	roomManager *RoomManager
}

// NewHub creates a new Hub instance.
// Call Run() in a separate goroutine to start the event loop.
func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]map[*Client]bool),
		broadcast:   make(chan WSMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		roomManager: NewRoomManager(),
	}
}

// Run starts the hub's event loop.
// WHY GOROUTINE: This runs forever, handling all client registration/unregistration.
// It must be started in a separate goroutine to avoid blocking.
func (h *Hub) Run() {
	logger := logging.Logger()
	logger.Info("hub_started")

	for {
		// select statement is the event loop - waits for any channel operation.
		// This is the core of the hub's coordination logic.
		select {
		case client := <-h.register:
			// New client is joining a room
			logger := logging.WSLogger(client.RoomCode, client.PlayerID)

			// WHY CHECK: Lazily create room entry if it doesn't exist
			if _, ok := h.rooms[client.RoomCode]; !ok {
				h.rooms[client.RoomCode] = make(map[*Client]bool)
				logger.Info("room_created_by_first_client")
			}

			// Add client to the room's client set
			h.rooms[client.RoomCode][client] = true

			clientCount := len(h.rooms[client.RoomCode])
			logger.Info("client_registered",
				"room_client_count", clientCount,
				"total_rooms", len(h.rooms),
			)

		case client := <-h.unregister:
			logger := logging.WSLogger(client.RoomCode, client.PlayerID)

			// Client is disconnecting
			// WHY NESTED CHECKS: Client might have already been removed
			if _, ok := h.rooms[client.RoomCode]; ok {
				if _, ok := h.rooms[client.RoomCode][client]; ok {
					// Remove client from room
					delete(h.rooms[client.RoomCode], client)
					// Close client's send channel to signal writePump to exit
					close(client.Send)

					clientCount := len(h.rooms[client.RoomCode])
					logger.Info("client_unregistered",
						"room_client_count", clientCount,
					)

					// Clean up empty rooms to prevent memory leak
					if len(h.rooms[client.RoomCode]) == 0 {
						delete(h.rooms, client.RoomCode)
						logger.Info("room_deleted_empty")
					}
				}
			}

		case message := <-h.broadcast:
			// Broadcast to all clients in all rooms
			// WHY: System announcements, global events
			logger := logging.Logger()
			totalClients := 0

			for roomCode := range h.rooms {
				for client := range h.rooms[roomCode] {
					totalClients++
					select {
					case client.Send <- message:
						// Message queued successfully
					default:
						// NON-BLOCKING SEND: Client's buffer is full (too slow).
						// Close connection and clean up to prevent blocking the hub.
						wsLogger := logging.WSLogger(roomCode, client.PlayerID)
						wsLogger.Warn("client_too_slow_disconnecting")

						close(client.Send)
						delete(h.rooms[roomCode], client)
					}
				}
			}

			logger.Debug("global_broadcast_sent",
				"message_type", message.Type,
				"total_clients", totalClients,
			)
		}
	}
}

// Register returns the register channel for clients.
// WHY METHOD: Keeps hub's internal channel structure encapsulated.
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister returns the unregister channel for clients.
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}

// BroadcastToRoom sends a message to all clients in a specific room.
// CHANNEL PATTERN: Fan-out (one hub -> many clients).
// Iterates over room's clients and sends to each client's buffered channel.
func (h *Hub) BroadcastToRoom(roomCode string, message WSMessage) {
	logger := logging.RoomLogger(roomCode)

	if clients, ok := h.rooms[roomCode]; ok {
		sentCount := 0
		dropCount := 0

		for client := range clients {
			select {
			case client.Send <- message:
				// Message queued successfully - client will receive it
				sentCount++
			default:
				// Buffer full - client is too slow. Disconnect them.
				// WHY: Prevents one slow client from blocking broadcasts to others.
				wsLogger := logging.WSLogger(roomCode, client.PlayerID)
				wsLogger.Warn("broadcast_to_slow_client_dropping")

				close(client.Send)
				delete(clients, client)
				dropCount++
			}
		}

		logger.Debug("room_broadcast_sent",
			"message_type", message.Type,
			"sent_count", sentCount,
			"dropped_count", dropCount,
		)
	} else {
		logger.Warn("broadcast_to_nonexistent_room")
	}
}

// GetRoomManager returns the room manager for game logic access.
func (h *Hub) GetRoomManager() *RoomManager {
	return h.roomManager
}

// GetRoomClientCount returns the number of connected clients in a room.
// Useful for debugging and monitoring.
func (h *Hub) GetRoomClientCount(roomCode string) int {
	if clients, ok := h.rooms[roomCode]; ok {
		return len(clients)
	}
	return 0
}
