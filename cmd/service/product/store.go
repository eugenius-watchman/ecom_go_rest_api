package product

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// method to implement ProductExists interface
func (s *Store) ProductExists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`
	err := s.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	const query = `
			SELECT id, name, description, image, price, quantity, createdAt
			FROM products
			ORDER BY createdAt DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		p, err := scanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	const query = `
		SELECT id, name, description, image, price, quantity, createdAt 
		FROM products WHERE id = ?`

	row := s.db.QueryRow(query, id)

	product, err := scanRowIntoProduct(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	const query = `
			INSERT INTO products (name, description, image, price, quantity)
				VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity,
	)
	return err
}


// update product quantity
func (s *Store) UpdateProductQuantity(id int, newQuantity int) error {
	query := `UPDATE products SET quantity = $1, updated_at = $2 WHERE id = $3`

	_, err := s.db.Exec(query, newQuantity, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating product quantity: %w", err)
	}
	return nil
}



func (s *Store) UpdateProduct(id int, product types.Product) error {
	const query = `
			UPDATE products
			SET name = ?, description = ?, image = ?, price = ?, quantity = ?
			WHERE id = ?`

	result, err := s.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// check if any row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}

func scanRowIntoProduct(row *sql.Row) (*types.Product, error) {
	product := new(types.Product)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func scanRowsIntoProducts(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
