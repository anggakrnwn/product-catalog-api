package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/anggakrnwn/product-catalog-api/models"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	GetByID(id uuid.UUID) (*models.Category, error)
	Create(category *models.Category) error
	Update(id uuid.UUID, category *models.Category) error
	Delete(id uuid.UUID) error
	FindByName(name string) (*models.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {

	query := "SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {

	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1"

	var c models.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &c, nil
}

func (r *categoryRepository) Create(category *models.Category) error {

	query := `
    INSERT INTO categories (id, name, description) 
    VALUES ($1, $2, $3)
    RETURNING created_at, updated_at
    `

	err := r.db.QueryRow(query,
		category.ID,
		strings.TrimSpace(category.Name),
		strings.TrimSpace(category.Description),
	).Scan(&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("category with this name already exists")
		}
		return err
	}

	return nil
}

func (r *categoryRepository) Update(id uuid.UUID, category *models.Category) error {

	query := `
    UPDATE categories 
    SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP 
    WHERE id = $3 
    RETURNING updated_at
    `

	result, err := r.db.Exec(query,
		strings.TrimSpace(category.Name),
		strings.TrimSpace(category.Description),
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("category with this name already exists")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return r.db.QueryRow(
		"SELECT updated_at FROM categories WHERE id = $1",
		id,
	).Scan(&category.UpdatedAt)
}

func (r *categoryRepository) Delete(id uuid.UUID) error {

	query := "DELETE FROM categories WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {

	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE name = $1"

	var c models.Category
	err := r.db.QueryRow(query, name).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
