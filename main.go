package main

import (
	"log"
	"os"

	"invela-be/internal/config"
	"invela-be/internal/database"
	"invela-be/internal/handlers"
	"invela-be/internal/repositories"
	"invela-be/internal/router"
	"invela-be/internal/services"
)

func main() {
	cfg := config.Load()
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	roleHandler := handlers.NewRoleHandler(roleService)

	engine := router.New(cfg, db, roleHandler)

	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
