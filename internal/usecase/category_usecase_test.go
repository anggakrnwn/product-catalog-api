package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
)

// ini fake repository yang saya buat untuk unit test
type fakeCategoryRepository struct {
	data map[string]domain.Category
}

func newFakeCategoryRepository() *fakeCategoryRepository {
	return &fakeCategoryRepository{
		data: make(map[string]domain.Category),
	}
}

func (f *fakeCategoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	var result []domain.Category
	for _, v := range f.data {
		result = append(result, v)
	}
	return result, nil
}

func (f *fakeCategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	c, ok := f.data[id]
	if !ok {
		return nil, errors.New("category not found")
	}
	return &c, nil
}

func (f *fakeCategoryRepository) Save(ctx context.Context, c domain.Category) error {
	f.data[c.ID] = c
	return nil
}

func (f *fakeCategoryRepository) Update(ctx context.Context, id string, c domain.Category) error {
	if _, ok := f.data[id]; !ok {
		return errors.New("category not found")
	}
	f.data[id] = c
	return nil
}

func (f *fakeCategoryRepository) Delete(ctx context.Context, id string) error {
	if _, ok := f.data[id]; !ok {
		return errors.New("category not found")
	}
	delete(f.data, id)
	return nil
}

// implementasi
// test create
func TestCreateCategory_Success(t *testing.T) {
	repo := newFakeCategoryRepository()
	uc := NewCategoryUsecase(repo)

	err := uc.Create(context.Background(), domain.Category{
		ID:   "cat-1",
		Name: "perabotan-lucu",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// test create tapi nama kosong
func TestCreateCategory_Failed_NameRequired(t *testing.T) {
	repo := newFakeCategoryRepository()
	uc := NewCategoryUsecase(repo)

	err := uc.Create(context.Background(), domain.Category{
		ID: "cat-1",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// test getbyid
func TestGetByID_Success(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.Save(context.Background(), domain.Category{
		ID:   "cat-1",
		Name: "daster",
	})

	uc := NewCategoryUsecase(repo)

	cat, err := uc.GetByID(context.Background(), "cat-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cat.Name != "daster" {
		t.Fatalf("expected Fashion, got %s", cat.Name)
	}
}

// test delete
func TestDeleteCategory(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.Save(context.Background(), domain.Category{
		ID:   "cat-1",
		Name: "gorengan",
	})

	uc := NewCategoryUsecase(repo)

	err := uc.Delete(context.Background(), "cat-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
