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

	kelasRepository := repositories.NewKelasRepository(db)
	kelasService := services.NewKelasService(kelasRepository)
	kelasHandler := handlers.NewKelasHandler(kelasService)

	jurusanRepository := repositories.NewJurusanRepository(db)
	jurusanService := services.NewJurusanService(jurusanRepository)
	jurusanHandler := handlers.NewJurusanHandler(jurusanService)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	engine := router.New(cfg, db, kelasHandler, jurusanHandler, userHandler)

	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
