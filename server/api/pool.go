package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Pool struct {
	Server     *Server
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	Username string `json:"username"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewPool(server *Server) *Pool {
	return &Pool{
		Server:     server,
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (p *Pool) router(server *Server) {
	poolRouter := chi.NewRouter()
	server.ChiRouter.Mount("/pool", poolRouter)
	poolRouter.Post("/createRoom", p.CreateRoom)
	poolRouter.Get("/joinRoom/{roomId}", p.JoinRoom)
	poolRouter.Get("/getRooms", p.GetRooms)
	poolRouter.Get("/getClients/{roomId}", p.GetClients)
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			if room, ok := p.Rooms[client.RoomID]; ok {
				if _, ok := room.Clients[client.Username]; !ok {
					room.Clients[client.Username] = client
					p.Broadcast <- &Message{
						Content:  fmt.Sprintf("User %s has joined the room", client.Username),
						RoomID:   client.RoomID,
						Username: client.Username,
					}
				}
			}
		case client := <-p.Unregister:
			if room, ok := p.Rooms[client.RoomID]; ok {
				if _, ok := room.Clients[client.Username]; ok {
					if len(room.Clients) != 0 {
						p.Broadcast <- &Message{
							Content:  fmt.Sprintf("User %s has left the chat room", client.Username),
							RoomID:   client.RoomID,
							Username: client.Username,
						}
					}
					delete(room.Clients, client.Username)
					close(client.Message)
				}
			}
		case message := <-p.Broadcast:
			if room, ok := p.Rooms[message.RoomID]; ok {
				for _, client := range room.Clients {
					client.Message <- message
				}
			}
		}
	}
}

func (p *Pool) CreateRoom(w http.ResponseWriter, r *http.Request) {
	req := CreateRoomReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler(w, http.StatusBadRequest, fmt.Sprintf("Error decoding %v", err))
		return
	}

	p.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	jsonHandler(w, http.StatusOK, req)
}

func (p *Pool) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, err.Error())
		return
	}

	roomID := chi.URLParam(r, "roomId")
	username := r.URL.Query().Get("username")

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		RoomID:   roomID,
		Username: username,
	}

	p.Register <- client
	go client.writeMessage()
	client.readMessage(p)
}

func (p *Pool) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)
	for _, r := range p.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}
	jsonHandler(w, http.StatusOK, rooms)
}

func (p *Pool) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []ClientRes
	roomId := chi.URLParam(r, "roomId")

	if _, ok := p.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		jsonHandler(w, http.StatusOK, clients)
		return
	}

	for _, c := range p.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			Username: c.Username,
		})
	}

	if clients == nil {
		jsonHandler(w, http.StatusOK, struct{}{})
		return
	}
	jsonHandler(w, http.StatusOK, clients)
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(pool *Pool) {
	defer func() {
		pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		pool.Broadcast <- msg
		fmt.Printf("Message Received '%s' from: %s in Room %v\n", msg.Content, msg.Username, msg.RoomID)
	}
}
