package repository

import (
	"context"
	"myservice/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo interface {
	Create(ctx context.Context, o *models.Order) error
	CancelOrder(ctx context.Context, id int, status string) error
	GetByID(ctx context.Context, id int) (models.Order, error)
	GetMyOrders(ctx context.Context, userID int) ([]models.Order, error)
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepo {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *models.Order) error {
	query := `INSERT INTO orders (user_id, item_name, quantity, status, created_at) 
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	
	return r.db.QueryRow(ctx, query, o.UserID, o.ItemName, o.Quantity, o.Status).Scan(&o.ID, &o.CreatedAt)
}

func (r *OrderRepository) CancelOrder(ctx context.Context, id int, status string) error {
	query := `UPDATE orders SET status = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (models.Order, error) {
	query := `
	SELECT 
		id,
		user_id, 
		item_name, 
		quantity, 
		status, 
		created_at
	FROM orders 
	WHERE id = $1`
	var order models.Order

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID, 
		&order.UserID, 
		&order.ItemName,
		&order.Quantity,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) GetMyOrders(ctx context.Context, userID int) ([]models.Order, error) {
	query := `
	SELECT 
		id,
		user_id,
		item_name,
		quantity,
		status,
		created_at
	FROM orders 
	WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ItemName, &o.Quantity, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}