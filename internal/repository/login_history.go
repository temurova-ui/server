package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"myservice/internal/models"
)

type LoginHistoryRepository struct {
	db *pgxpool.Pool
}

func NewLoginHistoryRepository(db *pgxpool.Pool) *LoginHistoryRepository {
	return &LoginHistoryRepository{
		db: db,
	}
}

func (r *LoginHistoryRepository) AddLoginHistory(ctx context.Context, history models.LoginHistory) error {
	const query = `INSERT INTO login_history(user_id, ip_address, user_agent)VALUES ($1, $2, $3)`
	rows, err := r.db.Exec(ctx, query, history.UserID, history.IPAddress, history.UserAgent)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return errors.New("insert error")
	}
	return nil
}

func (r *LoginHistoryRepository) GetLoginHistory(ctx context.Context, userID int64) ([]models.LoginHistory, error) {
	const query = 	`SELECT
					ip_address, 
					user_agent,
					created_at
				FROM login_history 
				WHERE user_id = $1 
				ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil{
		return nil, fmt.Errorf("query login_history: %w", err)
	}
	defer rows.Close()

	var history []models.LoginHistory
	for rows.Next(){
		var historyItem models.LoginHistory
		err = rows.Scan(
			&historyItem.IPAddress,
			&historyItem.UserAgent,
			&historyItem.CreatedAt,
		)
		if err != nil{
			return nil, fmt.Errorf("scan login_history: %w", err)
		}
		history = append(history, historyItem)
	}
	return history, nil
}
