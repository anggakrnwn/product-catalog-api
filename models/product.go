package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Price      int64     `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uuid.UUID `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name       string    `json:"name" binding:"required"`
	Price      int64     `json:"price" binding:"required,min=0"`
	Stock      int       `json:"stock" binding:"min=0"`
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name       *string    `json:"name,omitempty"`
	Price      *int64     `json:"price,omitempty" binding:"min=0"`
	Stock      *int       `json:"stock,omitempty" binding:"min=0"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `json:"category_name"`
}
