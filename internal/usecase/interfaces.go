package usecase

import (
	"context"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
)

type CategoryUsecase interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	Create(ctx context.Context, c domain.Category) error
	Update(ctx context.Context, id string, c domain.Category) error
	Delete(ctx context.Context, id string) error
}
