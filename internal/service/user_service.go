package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/release-platform/user-service/internal/model"
	"github.com/release-platform/user-service/internal/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidInput    = errors.New("invalid input")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

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

func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}

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