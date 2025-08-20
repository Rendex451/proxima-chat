package models

import "sync"

type Room struct {
	Mux     sync.RWMutex
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
