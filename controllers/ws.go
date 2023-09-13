package controllers

import (
	"github.com/AYGA2K/GoFiberWebSockets/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type (
	Handler struct {
		hub *ws.Hub
	}
)

func NewHandler(h *ws.Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type ReqRoom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(ctx *fiber.Ctx) error {
	var room ReqRoom
	ctx.BodyParser(&room)
	h.hub.Rooms[room.ID] = &ws.Room{
		ID:      room.ID,
		Name:    room.Name,
		Clients: make(map[string]*ws.Client),
	}
	return ctx.Status(200).JSON("Successfully")
}

func (h *Handler) JoinRoom(c *websocket.Conn) {
	roomID := c.Params("roomId")
	clientID := c.Params("userId")
	username := c.Params("username")
	cl := &ws.Client{
		Conn:     c,
		Message:  make(chan *ws.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &ws.Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(ctx *fiber.Ctx) error {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}
	return ctx.Status(200).JSON(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(ctx *fiber.Ctx) error {
	var clients []ClientRes
	roomId := ctx.Params("roomId")
	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)

		return ctx.Status(200).JSON(clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}
	return ctx.Status(200).JSON(clients)
}
