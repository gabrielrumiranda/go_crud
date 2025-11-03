package repository

import (
	"context"
	"errors"
	"fmt"
	"go_crud/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, name, price, stock_quantity, created_at, updated_at
		FROM products
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.StockQuantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}

		products = append(products, product)

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating products: %w", err)
		}
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, price, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("error fetching product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products(name, price, stock_quantity)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, stock_quantity, created_at, updated_at
	`

	var newProduct models.Product
	err := r.db.QueryRow(ctx, query, product.Name, product.Price, product.StockQuantity).Scan(
		&newProduct.ID,
		&newProduct.Name,
		&newProduct.Price,
		&newProduct.StockQuantity,
		&newProduct.CreatedAt,
		&newProduct.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return &newProduct, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, product *models.UpdateProductRequest) (*models.Product, error) {
	existing, err := r.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if product.Name != "" {
		existing.Name = product.Name
	}
	if product.Price > 0.0 {
		existing.Price = product.Price
	}
	if product.StockQuantity > 0 {
		existing.StockQuantity = product.StockQuantity
	}

	query := `
		UPDATE products
		SET name = $1, price = $2, stock_quantity = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, name, price, stock_quantity, created_at, updated_at
	`

	var updateProduct models.Product

	err = r.db.QueryRow(ctx, query, existing.Name, existing.Price, existing.StockQuantity, id).Scan(
		&updateProduct.ID,
		&updateProduct.Name,
		&updateProduct.Price,
		&updateProduct.StockQuantity,
		&updateProduct.CreatedAt,
		&updateProduct.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return &updateProduct, err
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products where id = $1`

	result, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
