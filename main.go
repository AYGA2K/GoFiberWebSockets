package main

import (
	"log"

	"github.com/AYGA2K/GoFiberWebSockets/controllers"
	"github.com/AYGA2K/GoFiberWebSockets/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	hub := ws.NewHub()
	wsHandler := controllers.NewHandler(hub)

	app.Get("ws/:roomId/:userId/:username", websocket.New(wsHandler.JoinRoom))
	app.Post("ws/createRoom", wsHandler.CreateRoom)

	log.Fatal(app.Listen(":3000"))
}
