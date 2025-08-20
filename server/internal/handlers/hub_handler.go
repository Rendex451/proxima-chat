package handlers

import (
	"encoding/json"
	"net/http"

	"server/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type HubHandler struct {
	hub *models.Hub
}

func NewHubHandler(hub *models.Hub) *HubHandler {
	return &HubHandler{hub: hub}
}

func (h *HubHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.hub.Mux.Lock()
	h.hub.Rooms[req.ID] = &models.Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*models.Client),
	}
	h.hub.Mux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *HubHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomID := chi.URLParam(r, "roomId")
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	cl := &models.Client{
		Conn:     conn,
		Message:  make(chan *models.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &models.Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	go cl.ReadMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *HubHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	h.hub.Mux.RLock()
	defer h.hub.Mux.RUnlock()

	rooms := make([]RoomRes, 0)
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *HubHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")

	h.hub.Mux.RLock()
	room, ok := h.hub.Rooms[roomId]
	h.hub.Mux.RUnlock()

	clients := make([]ClientRes, 0)
	if ok {
		for _, c := range room.Clients {
			clients = append(clients, ClientRes{
				ID:       c.ID,
				Username: c.Username,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
