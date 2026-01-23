package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
)

type categoryRepository struct {
	data map[string]domain.Category
	mu   sync.RWMutex
}

func NewCategoryRepository() *categoryRepository {
	return &categoryRepository{
		data: make(map[string]domain.Category),
	}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var categories []domain.Category
	for _, category := range r.data {
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, exists := r.data[id]
	if !exists {
		return nil, errors.New("not found")
	}

	return &category, nil

}

func (r *categoryRepository) Save(ctx context.Context, c domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.data {
		if existing.Name == c.Name {
			return errors.New("category with this name already exists")
		}
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.data[c.ID] = c
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, id string, c domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[id]
	if !exists {
		return errors.New("category not found")
	}

	for _, cat := range r.data {
		if cat.ID != id && cat.Name == c.Name {
			return errors.New("category with this name already exists")
		}
	}

	c.ID = id
	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now()
	r.data[id] = c
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("category not found")
	}

	delete(r.data, id)
	return nil
}

func (r *categoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = make(map[string]domain.Category)
}
