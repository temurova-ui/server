package repository

import(
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"myservice/internal/models"
)

type UserRepo interface{
	Create (ctx context.Context, u *models.User)error
	GetByEmail(ctx context.Context, email string)(*models.User, error)
	GetByID(ctx context.Context, id int)(*models.User, error)
	Update (ctx context.Context, id int, name, email string)error
	UpdatePassword(ctx context.Context, id int, password string)error
	Delete (ctx context.Context, id int)error
}
type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool)UserRepo{
	return &UserRepository{db: db}
}

func (r *UserRepository)Create (ctx context.Context, u *models.User)error{
	query := `INSERT INTO users(name, email, password_hash, role, created_at)
			 VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, u.Name, u.Email, u.PasswordHash, u.Role).Scan(&u.ID, &u.CreatedAt)			
}

func (r *UserRepository)GetByEmail(ctx context.Context, email string)(*models.User, error){
	query := `
	SELECT 
		id, 
		name, 
		email, 
		password_hash, 
		role, 
		created_at 
	FROM users 
	WHERE email = $1 `
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, 
		&u.Name, 
		&u.Email, 
		&u.PasswordHash, 
		&u.Role, 
		&u.CreatedAt)
	if err != nil{
		return nil, err
	}
	return u, nil
}

func (r *UserRepository)GetByID(ctx context.Context, id int)(*models.User, error){
	query := `
	SELECT 
		id, 
		name, 
		email, 
		password_hash,
		role, 
		created_at 
	FROM users 
	WHERE id = $1`
	u := &models.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, 
		&u.Name, 
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt)
	if err != nil{
		return nil, err
	}
	return u, nil
}

func (r *UserRepository)Update (ctx context.Context, id int, name, email string)error{
	query := `UPDATE users SET name = $1 ,email = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, name, email, id)
	return err
}

func (r *UserRepository)UpdatePassword(ctx context.Context, id int, password string)error{
	query := `UPDATE users SET password_hash = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, password)
	return err
}

func (r *UserRepository)Delete (ctx context.Context, id int)error{
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}