package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/anggakrnwn/product-catalog-api/config"
	"github.com/anggakrnwn/product-catalog-api/database"
	"github.com/anggakrnwn/product-catalog-api/handlers"
	"github.com/anggakrnwn/product-catalog-api/repositories"
	"github.com/anggakrnwn/product-catalog-api/services"
)

func main() {
	// load config
	cfg := config.Load()

	// setup database
	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// dependency injection
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	// setup router
	// categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAll(w, r)
		case http.MethodPost:
			categoryHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetByID(w, r)
		case http.MethodPut:
			categoryHandler.Update(w, r)
		case http.MethodDelete:
			categoryHandler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/categories/bulk", categoryHandler.BulkCreate)

	// products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("category_id") != "" {
				productHandler.GetByCategory(w, r)
			} else {
				productHandler.GetAll(w, r)
			}
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPut:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// home dan health
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/health", healthHandler)

	// start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"message":   "Product Catalog API",
		"version":   "v1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  "Supabase PostgreSQL",
		"challenge": []string{
			"Products use category_id as foreign key",
			"JOIN products and categories on product detail",
			"GET /api/products/{id} returns category name",
		},
		"endpoints": []map[string]string{
			{"method": "GET", "path": "/api/categories", "description": "List all categories"},
			{"method": "POST", "path": "/api/categories", "description": "Create category"},
			{"method": "GET", "path": "/api/categories/{id}", "description": "Get category by ID"},
			{"method": "PUT", "path": "/api/categories/{id}", "description": "Update category"},
			{"method": "DELETE", "path": "/api/categories/{id}", "description": "Delete category"},
			{"method": "POST", "path": "/api/categories/bulk", "description": "Bulk create categories"},

			{"method": "GET", "path": "/api/products", "description": "List all products (optional query: category_id=uuid)"},
			{"method": "POST", "path": "/api/products", "description": "Create product with category_id"},
			{"method": "GET", "path": "/api/products/{id}", "description": "Get product detail with category name (JOIN)"},
			{"method": "PUT", "path": "/api/products/{id}", "description": "Update product"},
			{"method": "DELETE", "path": "/api/products/{id}", "description": "Delete product"},

			{"method": "GET", "path": "/health", "description": "Health check"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "product-catalog-api",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  "connected",
		"tables":    []string{"categories", "products"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
