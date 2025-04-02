package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SignalMessage represents a signaling message for WebRTC
type SignalMessage struct {
	Type      string `json:"type"` // "offer", "answer", or "candidate"
	SDP       string `json:"sdp"`
	Candidate string `json:"candidate"`
	From      string `json:"from"`
	RoomID    string `json:"room_id"`
}

// RegisterFiberRoutes sets up middleware and routes
func (s *FiberServer) RegisterFiberRoutes() {
	// Enable CORS
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Middleware to allow WebSocket upgrades
	s.App.Use(func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

// RegisterSocketIo sets up WebSocket event handlers
func (s *FiberServer) RegisterSocketIo() {
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		fmt.Printf("User connected: %s\n", userID)
	})

	socketio.On("join-room", func(ep *socketio.EventPayload) {
		var msg SignalMessage
		if err := json.Unmarshal(ep.Data, &msg); err != nil {
			log.Println("Error parsing join-room message:", err)
			return
		}

		userID := ep.Kws.GetStringAttribute("user_id")
		roomID := msg.RoomID

		fmt.Printf("User %s joined room %s\n", userID, roomID)

		// Add user to the room
		s.mu.Lock()
		s.Rooms[roomID] = append(s.Rooms[roomID], ep.Kws.UUID)
		s.mu.Unlock()

		// Notify all users in the room about the new user
		broadcastToRoom(roomID, fmt.Sprintf("User %s joined room %s", userID, roomID), s)
	})

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		var signalMsg SignalMessage
		if err := json.Unmarshal(ep.Data, &signalMsg); err != nil {
			log.Println("Error parsing signaling message:", err)
			return
		}

		roomID := signalMsg.RoomID

		// Relay the message to all other users in the same room
		for _, memberUUID := range s.Rooms[roomID] {
			if memberUUID != ep.Kws.UUID { // Avoid sending to self
				err := ep.Kws.EmitTo(memberUUID, ep.Data, socketio.TextMessage)
				if err != nil {
					log.Printf("Failed to send message to user in room %s: %v\n", roomID, err)
				} else {
					log.Printf("Relayed message from %s to room %s\n", signalMsg.From, roomID)
				}
			}
		}
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		s.mu.Lock()
		delete(s.Clients, userID)

		// Remove user from all rooms
		for roomID, members := range s.Rooms {
			for i, memberUUID := range members {
				if memberUUID == ep.Kws.UUID {
					s.Rooms[roomID] = append(members[:i], members[i+1:]...)
					break
				}
			}
		}
		s.mu.Unlock()

		fmt.Printf("User disconnected: %s\n", userID)
	})

	// WebSocket route
	s.App.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {
		userID := kws.Params("id")

		// Add client to the map
		s.mu.Lock()
		s.Clients[userID] = kws.UUID
		s.mu.Unlock()

		fmt.Printf("New user connected: %s (UUID: %s)\n", userID, kws.UUID)

		// Set user ID as an attribute
		kws.SetAttribute("user_id", userID)

		// Broadcast welcome message
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s (UUID: %s)", userID, kws.UUID)), socketio.TextMessage)
	}))
}

func broadcastToRoom(roomID, message string, s *FiberServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, memberUUID := range s.Rooms[roomID] {
		// Use the `EmitTo` method to send the message to each user in the room
		err := socketio.EmitTo(memberUUID, []byte(message), socketio.TextMessage)
		if err != nil {
			log.Printf("Failed to broadcast to user in room %s: %v\n", roomID, err)
		}
	}
}
