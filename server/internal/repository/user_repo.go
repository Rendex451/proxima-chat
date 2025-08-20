package repository

import (
	"context"
	"server/internal/models"
)

type UserRepo struct {
	db DBTX
}

func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.ID)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
