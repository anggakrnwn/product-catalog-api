package services

import (
	"errors"
	"strings"

	"github.com/anggakrnwn/product-catalog-api/models"
	"github.com/anggakrnwn/product-catalog-api/repositories"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll() ([]models.Product, error)
	GetByID(id uuid.UUID) (*models.Product, error)
	GetWithCategory(id uuid.UUID) (*models.ProductWithCategory, error)
	Create(req *models.CreateProductRequest) (*models.Product, error)
	Update(id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error)
	Delete(id uuid.UUID) error
	GetByCategoryID(categoryID uuid.UUID) ([]models.Product, error)
}

type productService struct {
	repo         repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(repo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetByID(id uuid.UUID) (*models.Product, error) {

	if id == uuid.Nil {
		return nil, errors.New("product ID is required")
	}
	return s.repo.GetByID(id)
}

func (s *productService) GetWithCategory(id uuid.UUID) (*models.ProductWithCategory, error) {

	if id == uuid.Nil {
		return nil, errors.New("product ID is required")
	}
	return s.repo.GetWithCategory(id)
}

func (s *productService) Create(req *models.CreateProductRequest) (*models.Product, error) {

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if len(req.Name) < 3 {
		return nil, errors.New("name must be at least 3 characters")
	}

	if req.Price < 0 {
		return nil, errors.New("price must be positive")
	}

	if req.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	if req.CategoryID == uuid.Nil {
		return nil, errors.New("category ID is required")
	}

	_, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product := &models.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	err = s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Update(id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {

	if id == uuid.Nil {
		return nil, errors.New("product ID is required")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// untuk name
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		if len(name) < 3 {
			return nil, errors.New("name must be at least 3 characters")
		}
		existing.Name = name
	}

	// untuk price
	if req.Price != nil {
		if *req.Price < 0 {
			return nil, errors.New("price must be positive")
		}
		existing.Price = *req.Price
	}

	// untuk stock
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("stock cannot be negative")
		}
		existing.Stock = *req.Stock
	}

	// untuk category
	if req.CategoryID != nil {
		if *req.CategoryID == uuid.Nil {
			return nil, errors.New("category ID is invalid")
		}

		_, err := s.categoryRepo.GetByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		existing.CategoryID = *req.CategoryID
	}

	if err := s.repo.Update(id, existing); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *productService) Delete(id uuid.UUID) error {

	if id == uuid.Nil {
		return errors.New("product ID is required")
	}

	return s.repo.Delete(id)
}

func (s *productService) GetByCategoryID(categoryID uuid.UUID) ([]models.Product, error) {
	if categoryID == uuid.Nil {
		return nil, errors.New("category ID is required")
	}

	return s.repo.GetByCategoryID(categoryID)
}
