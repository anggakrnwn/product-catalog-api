package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
	"github.com/anggakrnwn/product-catalog-api/internal/dto"
	"github.com/anggakrnwn/product-catalog-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	uc usecase.CategoryUsecase
}

func NewCategoryHandler(uc usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"data":    []interface{}{},
		})
		return
	}

	var responses []dto.CategoryResponse
	for _, category := range categories {
		responses = append(responses, dto.ToCategoryResponse(category))
	}

	// jika kosong
	if responses == nil {
		responses = []dto.CategoryResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"meta": gin.H{
			"count": len(responses),
			"total": len(responses),
		},
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "category ID is required",
		})
		return
	}

	category, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dto.ToCategoryResponse(*category),
	})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	category := domain.Category{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.uc.Create(c.Request.Context(), category)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "category created successfully",
		"data":    dto.ToCategoryResponse(category),
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "category ID is required",
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	existing, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "category not found",
		})
		return
	}

	updatedCategory := domain.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err = h.uc.Update(c.Request.Context(), id, updatedCategory)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := h.uc.GetByID(c.Request.Context(), id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "category updated successfully",
		"data":    dto.ToCategoryResponse(*updated),
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "category ID is required",
		})
		return
	}

	existing, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "category not found",
		})
		return
	}

	err = h.uc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "category deleted successfully",
		"data": gin.H{
			"id": id,
		},
	})
}

func (h *CategoryHandler) BulkCreate(c *gin.Context) {
	var req dto.BulkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var created []dto.CategoryResponse
	var errors []string

	for _, item := range req.Categories {
		category := domain.Category{
			ID:          uuid.New().String(),
			Name:        strings.TrimSpace(item.Name),
			Description: strings.TrimSpace(item.Description),
		}

		err := h.uc.Create(c.Request.Context(), category)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			created = append(created, dto.ToCategoryResponse(category))
		}
	}

	response := gin.H{
		"success": true,
		"created": len(created),
		"failed":  len(errors),
		"data":    created,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusCreated, response)
}
