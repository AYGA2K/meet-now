package server

import (
	"api/internal/database"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	db      database.Service
	Clients map[string]string
	Rooms   map[string][]string
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "api",
			AppName:      "api",
		}),
		Clients: make(map[string]string),
		Rooms:   make(map[string][]string),
		db:      database.New(),
	}

	return server
}
