package usecase

import (
	"context"
	"errors"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
	"github.com/anggakrnwn/product-catalog-api/internal/repository"
)

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(r repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: r}
}

func (u *categoryUsecase) GetAll(ctx context.Context) ([]domain.Category, error) {
	return u.repo.FindAll(ctx)
}

func (u *categoryUsecase) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return u.repo.FindByID(ctx, id)
}

func (u *categoryUsecase) Create(ctx context.Context, c domain.Category) error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return u.repo.Save(ctx, c)
}

func (u *categoryUsecase) Update(ctx context.Context, id string, c domain.Category) error {
	if id == "" {
		return errors.New("id is required")
	}
	return u.repo.Update(ctx, id, c)
}

func (u *categoryUsecase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return u.repo.Delete(ctx, id)
}
