package game

type Hub struct {
	rooms       map[string]map[*Client]bool
	broadcast   chan Message
	register    chan *Client
	unregister  chan *Client
	roomManager *RoomManager
}

type Client struct {
	hub      *Hub
	roomCode string
	playerID string
	send     chan Message
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]map[*Client]bool),
		broadcast:   make(chan Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		roomManager: NewRoomManager(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.rooms[client.roomCode]; !ok {
				h.rooms[client.roomCode] = make(map[*Client]bool)
			}
			h.rooms[client.roomCode][client] = true

		case client := <-h.unregister:
			if _, ok := h.rooms[client.roomCode]; ok {
				if _, ok := h.rooms[client.roomCode][client]; ok {
					delete(h.rooms[client.roomCode], client)
					close(client.send)
					if len(h.rooms[client.roomCode]) == 0 {
						delete(h.rooms, client.roomCode)
					}
				}
			}

		case message := <-h.broadcast:
			for roomCode := range h.rooms {
				for client := range h.rooms[roomCode] {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.rooms[roomCode], client)
					}
				}
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(roomCode string, message Message) {
	if clients, ok := h.rooms[roomCode]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

func (h *Hub) GetRoomManager() *RoomManager {
	return h.roomManager
}
