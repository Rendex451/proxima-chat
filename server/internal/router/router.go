package router

import (
	"log"
	"net/http"
	"server/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(hubHandler *handlers.HubHandler, userHandler *handlers.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/ws/rooms", hubHandler.CreateRoom)
	r.Get("/ws/rooms", hubHandler.GetRooms)
	r.Get("/ws/rooms/{roomId}/clients", hubHandler.GetClients)
	r.Get("/ws/{roomId}", hubHandler.JoinRoom)

	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)
	r.Get("/logout", userHandler.Logout)

	return r
}

func Start(addr string, router *chi.Mux) error {
	log.Println("Server is running on", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("could not start server: %s", err)
		return err
	}
	return nil
}
