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
	log.Println(" Connected to Database")

	// Auth
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// Kelas
	kelasRepository := repositories.NewKelasRepository(db)
	kelasService := services.NewKelasService(kelasRepository)
	kelasHandler := handlers.NewKelasHandler(kelasService)

	// Jurusan
	jurusanRepository := repositories.NewJurusanRepository(db)
	jurusanService := services.NewJurusanService(jurusanRepository)
	jurusanHandler := handlers.NewJurusanHandler(jurusanService)

	// User
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	// Perangkat
	perangkatRepository := repositories.NewPerangkatRepository(db)
	perangkatService := services.NewPerangkatService(perangkatRepository)
	perangkatHandler := handlers.NewPerangkatHandler(perangkatService)

	// Kategori
	kategoriRepository := repositories.NewKategoriRepository(db)
	kategoriService := services.NewKategoriService(kategoriRepository)
	kategoriHandler := handlers.NewKategoriHandler(kategoriService)

	// Labor
	laborRepository := repositories.NewLaborRepository(db)
	laborService := services.NewLaborService(laborRepository)
	laborHandler := handlers.NewLaborHandler(laborService)

	// Item Instance
	itemInstanceRepository := repositories.NewItemInstanceRepository(db)
	itemInstanceService := services.NewItemInstanceService(itemInstanceRepository)
	itemInstanceHandler := handlers.NewItemInstanceHandler(itemInstanceService)

	// Peminjaman
	peminjamanRepository := repositories.NewPeminjamanRepository(db)
	peminjamanService := services.NewPeminjamanService(peminjamanRepository)
	peminjamanHandler := handlers.NewPeminjamanHandler(peminjamanService)

	// Penggunaan
	penggunaanRepository := repositories.NewPenggunaanRepository(db)
	penggunaanService := services.NewPenggunaanService(penggunaanRepository)
	penggunaanHandler := handlers.NewPenggunaanHandler(penggunaanService)

	// Kerusakan
	kerusakanRepository := repositories.NewKerusakanRepository(db)
	kerusakanService := services.NewKerusakanService(kerusakanRepository)
	kerusakanHandler := handlers.NewKerusakanHandler(kerusakanService)

	engine := router.New(cfg, db, authHandler, kelasHandler, jurusanHandler, userHandler, perangkatHandler, kategoriHandler, laborHandler, itemInstanceHandler, peminjamanHandler, penggunaanHandler, kerusakanHandler)
	log.Println(" Server Started")
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
