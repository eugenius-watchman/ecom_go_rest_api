package cart

import (
	"database/sql"
	"fmt"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	const query = `
		INSERT INTO orders (userId, total, status, address, createdAt)
		VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		order.UserID,
		order.Total,
		order.Status,
		order.Address,
		order.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get order ID: %w", err)
	}

	return int(orderID), nil

}

func (s *Store) CreateOrderItem(item types.OrderItem) error {
	const query = `
			INSERT INTO order_items (orderId, productId, quantity, price)
				VALUES (?, ?, ?, ?)`
	
	_, err := s.db.Exec(
		query,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
	)
	if err != nil {
		return fmt.Errorf("failed too create order item: %w", err)
	}

	return nil
}

func (s *Store) GetOrderByID(id int) (*types.Order, error) {
	const query = `
		SELECT id, userId, total, status, address, createdAt 
		FROM orders WHERE id = ?`

	row := s.db.QueryRow(query, id)

	var order types.Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil

}

func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	const query = `
		SELECT id, userId, total, status, address, createdAt
		FROM orders WHERE userId = ?
		ORDER BY createdAt DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query oders: %w", err)
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return orders, nil
}