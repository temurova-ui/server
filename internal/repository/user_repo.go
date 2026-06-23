package repository

import(
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"myservice/internal/models"
)

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool)*UserRepository{
	return &UserRepository{db: db}
}

func (r *UserRepository)Create (ctx context.Context, u *models.User)error{
	query := `INSERT INTO users(name, email, password_hash, role, created_at)
			 VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, u.Name, u.Email, u.PasswordHash, u.Role).Scan(&u.ID, &u.CreatedAt)			
}

func (r *UserRepository)GetByEmail(ctx context.Context, email string)(*models.User, error){
	query := `SELECT id, name, email, password_hash, role, role, created_at FROM users WHERE email = $1 `
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil{
		return nil, err
	}
	return u, nil
}