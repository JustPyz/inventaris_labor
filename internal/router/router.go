package router

import (
	"database/sql"
	"net/http"

	"invela-be/internal/config"
	"invela-be/internal/handlers"
	"invela-be/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(
	cfg config.Config,
	db *gorm.DB,
	authHandler *handlers.AuthHandler,
	kelasHandler *handlers.KelasHandler,
	jurusanHandler *handlers.JurusanHandler,
	userHandler *handlers.UserHandler,
	perangkatHandler *handlers.PerangkatHandler,
	kategoriHandler *handlers.KategoriHandler,
	laborHandler *handlers.LaborHandler,
	itemInstanceHandler *handlers.ItemInstanceHandler,
	peminjamanHandler *handlers.PeminjamanHandler,
	penggunaanHandler *handlers.PenggunaanHandler,
	kerusakanHandler *handlers.KerusakanHandler,
	perbaikanHandler *handlers.PerbaikanHandler,
	riwayatPerbaikanHandler *handlers.RiwayatPerbaikanHandler,
) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middlewares.CORS())

	router.GET("/health", func(c *gin.Context) {
		status := "ok"
		dbState := "unknown"

		if db != nil {
			if sqlDB, err := db.DB(); err == nil {
				dbState = databaseState(sqlDB)
				if dbState != "ok" {
					status = "degraded"
				}
			} else {
				dbState = "error"
				status = "degraded"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"app":    "invela-be",
			"port":   cfg.Port,
			"db":     dbState,
		})
	})

	api := router.Group("/api")
	{
		// ── Public: Login (tidak butuh token) ──────────────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// ── Protected: semua route di bawah wajib token valid ──────────────
		protected := api.Group("")
		protected.Use(middlewares.RequireAuth(cfg.JWTSecret))
		{
			// /kelas — admin: CRUD | guru: GET
			kelas := protected.Group("/kelas")
			{
				kelas.GET("", middlewares.RequireRole("admin", "guru", "kabeng"), kelasHandler.List)
				kelas.GET("/:id", middlewares.RequireRole("admin", "guru", "kabeng"), kelasHandler.Get)
				kelas.POST("", middlewares.RequireRole("admin"), kelasHandler.Create)
				kelas.PUT("/:id", middlewares.RequireRole("admin"), kelasHandler.Update)
				kelas.DELETE("/:id", middlewares.RequireRole("admin"), kelasHandler.Delete)
			}

			// /jurusan — admin: CRUD
			jurusan := protected.Group("/jurusan")
			{
				jurusan.GET("", middlewares.RequireRole("admin", "kabeng"), jurusanHandler.List)
				jurusan.POST("", middlewares.RequireRole("admin"), jurusanHandler.Create)
				jurusan.GET("/:id", middlewares.RequireRole("admin", "kabeng"), jurusanHandler.Get)
				jurusan.PUT("/:id", middlewares.RequireRole("admin"), jurusanHandler.Update)
				jurusan.DELETE("/:id", middlewares.RequireRole("admin"), jurusanHandler.Delete)
			}

			// /user — admin: CRUD
			users := protected.Group("/user")
			{
				users.GET("", middlewares.RequireRole("admin"), userHandler.List)
				users.POST("", middlewares.RequireRole("admin"), userHandler.Create)
				users.GET("/:id", middlewares.RequireRole("admin"), userHandler.Get)
				users.PATCH("/:id", middlewares.RequireRole("admin"), userHandler.Update)
				users.DELETE("/:id", middlewares.RequireRole("admin"), userHandler.Delete)
			}

			// /perangkat — kabeng & admin: CRUD | kaprog & sapras: GET
			perangkat := protected.Group("/perangkat")
			{
				perangkat.GET("", middlewares.RequireRole("kabeng", "admin", "guru", "kaprog", "sapras"), perangkatHandler.List)
				perangkat.GET("/:id", middlewares.RequireRole("kabeng", "admin", "guru", "kaprog", "sapras"), perangkatHandler.Get)
				perangkat.POST("", middlewares.RequireRole("kabeng", "admin"), perangkatHandler.Create)
				perangkat.PATCH("/:id", middlewares.RequireRole("kabeng", "admin"), perangkatHandler.Update)
				perangkat.DELETE("/:id", middlewares.RequireRole("kabeng", "admin"), perangkatHandler.Delete)
			}

			// /kategori — admin: CRUD
			kategori := protected.Group("/kategori")
			{
				kategori.GET("", middlewares.RequireRole("admin", "kabeng", "sapras", "kaprog"), kategoriHandler.List)
				kategori.POST("", middlewares.RequireRole("admin"), kategoriHandler.Create)
				kategori.GET("/:id", middlewares.RequireRole("admin", "kabeng", "sapras", "kaprog"), kategoriHandler.Get)
				kategori.PUT("/:id", middlewares.RequireRole("admin"), kategoriHandler.Update)
				kategori.DELETE("/:id", middlewares.RequireRole("admin"), kategoriHandler.Delete)
			}

			// /labor — admin: CRUD | guru: GET
			labor := protected.Group("/labor")
			{
				labor.GET("", middlewares.RequireRole("admin", "guru", "kabeng", "sapras", "kaprog"), laborHandler.List)
				labor.GET("/:id", middlewares.RequireRole("admin", "guru", "kabeng", "sapras", "kaprog"), laborHandler.Get)
				labor.POST("", middlewares.RequireRole("admin"), laborHandler.Create)
				labor.PUT("/:id", middlewares.RequireRole("admin"), laborHandler.Update)
				labor.DELETE("/:id", middlewares.RequireRole("admin"), laborHandler.Delete)
			}

			// /item-instance — kabeng & admin: CRUD | kaprog & sapras: GET
			itemInstance := protected.Group("/item-instance")
			{
				itemInstance.GET("", middlewares.RequireRole("kabeng", "admin", "guru", "kaprog", "sapras"), itemInstanceHandler.List)
				itemInstance.GET("/:id", middlewares.RequireRole("kabeng", "admin", "guru", "kaprog", "sapras"), itemInstanceHandler.Get)
				itemInstance.POST("", middlewares.RequireRole("kabeng", "admin"), itemInstanceHandler.Create)
				itemInstance.PUT("/:id", middlewares.RequireRole("kabeng", "admin"), itemInstanceHandler.Update)
				itemInstance.DELETE("/:id", middlewares.RequireRole("kabeng", "admin"), itemInstanceHandler.Delete)
			}

			// /peminjaman — kabeng: CRUD | kaprog: GET
			peminjaman := protected.Group("/peminjaman")
			{
				peminjaman.GET("", middlewares.RequireRole("kabeng", "kaprog"), peminjamanHandler.List)
				peminjaman.GET("/:id", middlewares.RequireRole("kabeng", "kaprog"), peminjamanHandler.Get)
				peminjaman.POST("", middlewares.RequireRole("kabeng"), peminjamanHandler.Create)
				peminjaman.PATCH("/:id", middlewares.RequireRole("kabeng"), peminjamanHandler.Update)
				peminjaman.DELETE("/:id", middlewares.RequireRole("kabeng"), peminjamanHandler.Delete)
			}

			// /penggunaan — kabeng: GET & DELETE | guru: POST
			penggunaan := protected.Group("/penggunaan")
			{
				penggunaan.GET("", middlewares.RequireRole("kabeng"), penggunaanHandler.List)
				penggunaan.GET("/:id", middlewares.RequireRole("kabeng"), penggunaanHandler.Get)
				penggunaan.POST("", middlewares.RequireRole("guru"), penggunaanHandler.Create)
				penggunaan.DELETE("/:id", middlewares.RequireRole("kabeng"), penggunaanHandler.Delete)
			}

			// /kerusakan — kabeng: GET, PUT, DELETE | guru: POST
			kerusakan := protected.Group("/kerusakan")
			{
				kerusakan.GET("", middlewares.RequireRole("kabeng", "kaprog", "sapras"), kerusakanHandler.List)
				kerusakan.GET("/:id", middlewares.RequireRole("kabeng"), kerusakanHandler.Get)
				kerusakan.POST("", middlewares.RequireRole("guru", "kabeng"), kerusakanHandler.Create)
				kerusakan.PUT("/:id", middlewares.RequireRole("kabeng"), kerusakanHandler.Update)
				kerusakan.DELETE("/:id", middlewares.RequireRole("kabeng"), kerusakanHandler.Delete)
				kerusakan.GET("/stats/:id_item_instance", middlewares.RequireRole("kabeng", "guru", "kaprog", "sapras", "admin"), kerusakanHandler.GetStats)
			}

			// /perbaikan — kabeng: CRUD
			perbaikan := protected.Group("/perbaikan")
			{
				perbaikan.GET("", middlewares.RequireRole("kabeng", "kaprog", "sapras"), perbaikanHandler.List)
				perbaikan.GET("/:id", middlewares.RequireRole("kabeng", "kaprog", "sapras"), perbaikanHandler.Get)
				perbaikan.POST("", middlewares.RequireRole("kabeng"), perbaikanHandler.Create)
				perbaikan.PUT("/:id", middlewares.RequireRole("kabeng"), perbaikanHandler.Update)
				perbaikan.DELETE("/:id", middlewares.RequireRole("kabeng"), perbaikanHandler.Delete)
			}

			// /riwayat-perbaikan — kabeng & admin: GET only (append-only audit log)
			riwayat := protected.Group("/riwayat-perbaikan")
			{
				riwayat.GET("", middlewares.RequireRole("kabeng", "admin", "sapras", "kaprog"), riwayatPerbaikanHandler.List)
				riwayat.GET("/:kode_asset", middlewares.RequireRole("kabeng", "admin", "sapras", "kaprog"), riwayatPerbaikanHandler.GetByKodeAsset)
			}
		}
	}

	return router  
}

func databaseState(db *sql.DB) string {
	if err := db.Ping(); err != nil {
		return "error"
	}

	return "ok"
}
