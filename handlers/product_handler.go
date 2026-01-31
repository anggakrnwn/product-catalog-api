package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anggakrnwn/product-catalog-api/models"
	"github.com/anggakrnwn/product-catalog-api/services"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []models.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    products,
		"meta": map[string]interface{}{
			"count": len(products),
			"total": len(products),
		},
	})
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/products/")

	if idStr == "" {
		http.Error(w, "product ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetWithCategory(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    product,
	})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.Create(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "required") {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "product created successfully",
		"data":    product,
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	if idStr == "" {
		http.Error(w, "product ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.Update(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "required") {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "product updated successfully",
		"data":    product,
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	if idStr == "" {
		http.Error(w, "product ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "product deleted successfully",
		"data": map[string]string{
			"id": id.String(),
		},
	})
}

func (h *ProductHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		http.Error(w, "category_id query parameter is required", http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	products, err := h.service.GetByCategoryID(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []models.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    products,
		"meta": map[string]interface{}{
			"count": len(products),
		},
	})
}
