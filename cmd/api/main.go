package main

import (
	"context"
	"log"

	"github.com/anggakrnwn/product-catalog-api/internal/config"
	"github.com/anggakrnwn/product-catalog-api/internal/delivery/http"
	"github.com/anggakrnwn/product-catalog-api/internal/middleware"
	"github.com/anggakrnwn/product-catalog-api/internal/repository/inmemory"
	"github.com/anggakrnwn/product-catalog-api/internal/seed"
	"github.com/anggakrnwn/product-catalog-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	// load config
	cfg := config.Load()

	// set gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// setup dependen
	categoryRepo := inmemory.NewCategoryRepository()

	// seed data
	seed.SeedCategories(context.Background(), categoryRepo)

	// business layer
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := http.NewCategoryHandler(categoryUsecase)

	// setup router
	router := gin.Default()

	// middleware
	router.Use(middleware.CORS())

	// routes
	setupRoutes(router, categoryHandler, cfg)

	// start server
	addr := ":" + cfg.Port
	log.Printf("server starting on %s", addr)
	log.Printf("environment: %s", cfg.Environment)

	if err := router.Run(addr); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, handler *http.CategoryHandler, cfg config.Config) {
	// root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "product catalog API",
			"version":     "v1.0.0",
			"environment": cfg.Environment,
			"endpoints": gin.H{
				"GET    /api/v1/categories":      "list all categories",
				"POST   /api/v1/categories":      "create category",
				"GET    /api/v1/categories/:id":  "get category",
				"PUT    /api/v1/categories/:id":  "update category",
				"DELETE /api/v1/categories/:id":  "delete category",
				"POST   /api/v1/categories/bulk": "Bulk create",
				"GET    /health":                 "health check",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.GET("", handler.GetAll)
			categories.POST("", handler.Create)
			categories.GET("/:id", handler.GetByID)
			categories.PUT("/:id", handler.Update)
			categories.DELETE("/:id", handler.Delete)
			categories.POST("/bulk", handler.BulkCreate)
		}

		// health
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":      "healthy",
				"storage":     "in-memory",
				"environment": cfg.Environment,
			})
		})
	}

	// legacy health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"storage": "in-memory",
		})
	})

	// 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "endpoint not found",
			"path":    c.Request.URL.Path,
		})
	})
}
