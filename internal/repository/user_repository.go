// Package repository provides storage implementations for domain models.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Release-Orchestrator/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepositoryInterface defines storage operations for users.
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
}

// UserRepository implements UserRepositoryInterface using Postgres.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository returns a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into storage.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetByID returns a user by id.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetAll returns all users.
func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.UpdatedAt, user.ID)
	return err
}

// Delete removes a user by id.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ExistsByEmail checks if a user exists with the given email.
// excludeID, if non-nil, excludes that ID from the check.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)`
		args = []interface{}{email, *excludeID}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
		args = []interface{}{email}
	}

	var exists bool
	err := r.db.QueryRow(ctx, query, args...).Scan(&exists)
	return exists, err
}
