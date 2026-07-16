package router

import (
	"database/sql"
	"net/http"

	"invela-be/internal/config"
	"invela-be/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB, roleHandler *handlers.RoleHandler) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

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

	api := router.Group("/api/v1")
	{
		roles := api.Group("/roles")
		{
			roles.GET("", roleHandler.List)
			roles.POST("", roleHandler.Create)
			roles.GET("/:id", roleHandler.Get)
			roles.PUT("/:id", roleHandler.Update)
			roles.DELETE("/:id", roleHandler.Delete)
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
