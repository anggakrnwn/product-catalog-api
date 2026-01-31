package services

import (
	"errors"
	"strings"
	"time"

	"github.com/anggakrnwn/product-catalog-api/models"
	"github.com/anggakrnwn/product-catalog-api/repositories"
	"github.com/google/uuid"
)

type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id string) (*models.Category, error)
	Create(req *models.CreateCategoryRequest) (*models.Category, error)
	Update(id string, req *models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id string) error
	BulkCreate(req *models.BulkCreateRequest) ([]models.Category, []error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id string) (*models.Category, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("category ID is required")
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}
	return s.repo.GetByID(categoryID)
}

func (s *categoryService) Create(req *models.CreateCategoryRequest) (*models.Category, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if len(req.Name) < 3 {
		return nil, errors.New("name must be at least 3 characters")
	}

	if len(req.Name) > 100 {
		return nil, errors.New("name must not exceed 100 characters")
	}

	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Update(id string, req *models.UpdateCategoryRequest) (*models.Category, error) {

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("category ID is required")
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if len(req.Name) < 3 {
		return nil, errors.New("name must be at least 3 characters")
	}

	if len(req.Name) > 100 {
		return nil, errors.New("name must not exceed 100 characters")
	}

	existing, err := s.repo.GetByID(categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	existingByName, _ := s.repo.FindByName(req.Name)
	if existingByName != nil && existingByName.ID != existing.ID {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		ID:          existing.ID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Update(categoryID, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id string) error {

	if strings.TrimSpace(id) == "" {
		return errors.New("category ID is required")
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid category ID format")
	}

	_, err = s.repo.GetByID(categoryID)
	if err != nil {
		return errors.New("category not found")
	}

	return s.repo.Delete(categoryID)
}

func (s *categoryService) BulkCreate(req *models.BulkCreateRequest) ([]models.Category, []error) {

	var created []models.Category
	var errs []error

	for _, item := range req.Categories {
		category, err := s.Create(&item)
		if err != nil {
			errs = append(errs, err)
		} else {
			created = append(created, *category)
		}
	}

	return created, errs
}
