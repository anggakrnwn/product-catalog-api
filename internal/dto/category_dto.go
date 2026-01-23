package dto

import (
	"time"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
)

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description,omitempty"`
}

type BulkCreateRequest struct {
	Categories []CreateCategoryRequest `json:"categories"`
}

func ToCategoryResponse(c domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToDomainFromCreate(req CreateCategoryRequest, id string) domain.Category {
	now := time.Now()
	return domain.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
