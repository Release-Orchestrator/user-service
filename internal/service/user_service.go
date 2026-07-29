// Package service implements business logic for users.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Release-Orchestrator/user-service/internal/model"
	"github.com/Release-Orchestrator/user-service/internal/repository"
	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned when a user cannot be found.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailExists is returned when an email is already in use.
	ErrEmailExists = errors.New("email already exists")
	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")
)

// UserServiceInterface defines user service operations.
type UserServiceInterface interface {
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserService implements UserServiceInterface.
type UserService struct {
	repo repository.UserRepositoryInterface
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	if req.Name == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	exists, err := s.repo.ExistsByEmail(ctx, req.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID returns a user by id.
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetAll returns all users.
func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}

// Update updates a user.
func (s *UserService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.Email != "" {
		exists, err := s.repo.ExistsByEmail(ctx, req.Email, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete deletes a user.
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}
