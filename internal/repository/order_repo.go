package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"myservice/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *models.Order) error {
	query := `INSERT INTO orders (user_id, item_name, quantity, status, created_at) 
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	
	return r.db.QueryRow(ctx, query, o.UserID, o.ItemName, o.Quantity, o.Status).Scan(&o.ID, &o.CreatedAt)
}