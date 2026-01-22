package repository

import (
	"context"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	Save(ctx context.Context, c domain.Category) error
	Update(ctx context.Context, id string, c domain.Category) error
	Delete(ctx context.Context, id string) error
}
