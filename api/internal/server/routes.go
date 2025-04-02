package server

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))
	s.App.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

func (s *FiberServer) RegisterSocketIo() {
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("User connected: %s\n", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On("join-room", func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		roomID := ep.Kws.GetStringAttribute("room_id")

		fmt.Printf("User %s joined room %s\n", userID, roomID)

		if _, exists := s.Rooms[roomID]; !exists {
			s.Rooms[roomID] = []string{}
		}
		s.Rooms[roomID] = append(s.Rooms[roomID], ep.Kws.UUID)
	})

	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		message := MessageObject{}
		if err := json.Unmarshal(ep.Data, &message); err != nil {
			fmt.Println("Error parsing message:", err)
			return
		}
		if message.Event != "" {
			ep.Kws.Fire(message.Event, []byte(message.Data))
		}
		if clientUUID, exists := s.Clients[message.To]; exists {
			err := ep.Kws.EmitTo(clientUUID, ep.Data, socketio.TextMessage)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		delete(s.Clients, userID)
		fmt.Printf("User disconnected: %s\n", userID)

		// Remove user from all rooms
		for roomID, members := range s.Rooms {
			for i, memberUUID := range members {
				if memberUUID == ep.Kws.UUID {
					s.Rooms[roomID] = slices.Delete(members, i, i+1)
					break
				}
			}
		}
	})

	s.App.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {
		userID := kws.Params("id")
		s.Clients[userID] = kws.UUID

		kws.Broadcast(fmt.Appendf(nil, "New user connected: %s (UUID: %s)", userID, kws.UUID), true, socketio.TextMessage)
		kws.Emit(fmt.Appendf(nil, "Hello user: %s (UUID: %s)", userID, kws.UUID), socketio.TextMessage)
	}))
}
