package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/anggakrnwn/product-catalog-api/models"
	"github.com/google/uuid"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id uuid.UUID) (*models.Product, error)
	GetWithCategory(id uuid.UUID) (*models.ProductWithCategory, error)
	Create(product *models.Product) error
	Update(id uuid.UUID, product *models.Product) error
	Delete(id uuid.UUID) error
	GetByCategoryID(categoryID uuid.UUID) ([]models.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]models.Product, error) {

	query := `
        SELECT 
            p.id, 
            p.name, 
            p.price, 
            p.stock, 
            p.category_id, 
            p.created_at, 
            p.updated_at,
            json_build_object(
                'id', c.id,
                'name', c.name,
                'description', c.description,
                'created_at', c.created_at,
                'updated_at', c.updated_at
            ) as category
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        ORDER BY p.name
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var categoryData []byte

		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock,
			&p.CategoryID, &p.CreatedAt, &p.UpdatedAt,
			&categoryData,
		)
		if err != nil {
			return nil, err
		}

		if len(categoryData) > 0 && string(categoryData) != "null" {
			var category models.Category
			if err := json.Unmarshal(categoryData, &category); err == nil {
				p.Category = &category
			}
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) GetByID(id uuid.UUID) (*models.Product, error) {

	query := `
		SELECT id, name, price, stock, category_id, created_at, updated_at 
		FROM products 
		WHERE id = $1
	`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
		&p.CategoryID, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &p, nil
}

// join
func (r *productRepository) GetWithCategory(id uuid.UUID) (*models.ProductWithCategory, error) {

	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, 
			p.category_id, p.created_at, p.updated_at,
			c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var result models.ProductWithCategory
	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.Name, &result.Price, &result.Stock,
		&result.CategoryID, &result.CreatedAt, &result.UpdatedAt,
		&result.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &result, nil
}

func (r *productRepository) Create(product *models.Product) error {

	query := `
		INSERT INTO products (name, price, stock, category_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		strings.TrimSpace(product.Name),
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return errors.New("category not found")
		}
		return err
	}

	return nil
}

func (r *productRepository) Update(id uuid.UUID, product *models.Product) error {

	query := `
		UPDATE products 
		SET name = $1, price = $2, stock = $3, 
		    category_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	result, err := r.db.Exec(
		query,
		strings.TrimSpace(product.Name),
		product.Price,
		product.Stock,
		product.CategoryID,
		id,
	)

	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return errors.New("category not found")
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) Delete(id uuid.UUID) error {

	query := "DELETE FROM products WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) GetByCategoryID(categoryID uuid.UUID) ([]models.Product, error) {

	query := `
		SELECT id, name, price, stock, category_id, created_at, updated_at 
		FROM products 
		WHERE category_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock,
			&p.CategoryID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
