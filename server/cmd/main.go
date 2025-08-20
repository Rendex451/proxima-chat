package main

import (
	"log"

	"server/internal/db"
	"server/internal/handlers"
	"server/internal/models"
	"server/internal/repository"
	"server/internal/router"
	"server/internal/service"

	_ "github.com/lib/pq"
)

var (
	addr = "0.0.0.0:8080"
)

func main() {
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRepo := repository.NewUserRepo(dbConn.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	hub := models.NewHub()
	hubHandler := handlers.NewHubHandler(hub)
	go hub.Run()

	r := router.InitRouter(hubHandler, userHandler)
	router.Start(addr, r)
}
