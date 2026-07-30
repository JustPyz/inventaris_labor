package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"invela-be/internal/config"
	"invela-be/internal/database"
	"invela-be/internal/handlers"
	"invela-be/internal/jobs"
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

	// Jurusan
	jurusanRepository := repositories.NewJurusanRepository(db)
	jurusanService := services.NewJurusanService(jurusanRepository)
	jurusanHandler := handlers.NewJurusanHandler(jurusanService)

	// Kelas
	kelasRepository := repositories.NewKelasRepository(db)
	kelasService := services.NewKelasService(kelasRepository, jurusanRepository)
	kelasHandler := handlers.NewKelasHandler(kelasService)

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
	peminjamanService := services.NewPeminjamanService(db, peminjamanRepository, itemInstanceRepository)
	peminjamanHandler := handlers.NewPeminjamanHandler(peminjamanService)

	// Penggunaan
	penggunaanRepository := repositories.NewPenggunaanRepository(db)
	penggunaanService := services.NewPenggunaanService(penggunaanRepository)
	penggunaanHandler := handlers.NewPenggunaanHandler(penggunaanService)

	// Kerusakan
	kerusakanRepository := repositories.NewKerusakanRepository(db)
	kerusakanService := services.NewKerusakanService(db, kerusakanRepository, itemInstanceRepository)
	kerusakanHandler := handlers.NewKerusakanHandler(kerusakanService)

	// Perbaikan
	perbaikanRepository := repositories.NewPerbaikanRepository(db)
	riwayatPerbaikanRepository := repositories.NewRiwayatPerbaikanRepository(db)
	perbaikanService := services.NewPerbaikanService(db, perbaikanRepository, kerusakanRepository, itemInstanceRepository, userRepository, riwayatPerbaikanRepository)
	perbaikanHandler := handlers.NewPerbaikanHandler(perbaikanService)

	// Riwayat Perbaikan
	riwayatPerbaikanService := services.NewRiwayatPerbaikanService(riwayatPerbaikanRepository)
	riwayatPerbaikanHandler := handlers.NewRiwayatPerbaikanHandler(riwayatPerbaikanService)

	engine := router.New(cfg, db, authHandler, kelasHandler, jurusanHandler, userHandler, perangkatHandler, kategoriHandler, laborHandler, itemInstanceHandler, peminjamanHandler, penggunaanHandler, kerusakanHandler, perbaikanHandler, riwayatPerbaikanHandler)

	// Background jobs
	overdueJob := jobs.NewPeminjamanOverdueJob(peminjamanRepository, 1*time.Hour)
	overdueJob.Start()

	// Graceful shutdown: tangkap sinyal SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println(" Shutting down...")
		overdueJob.Stop()
		os.Exit(0)
	}()

	log.Println(" Server Started")
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
