package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Release-Orchestrator/user-service/internal/model"
	"github.com/Release-Orchestrator/user-service/internal/repository"
	"github.com/google/uuid"
)

type mockRepo struct {
	repository.UserRepositoryInterface
	users map[uuid.UUID]*model.User

	createErr     error
	getByIDErr    error
	getAllErr     error
	updateErr     error
	deleteErr     error
	existsByEmail func(email string, excludeID *uuid.UUID) (bool, error)
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users: make(map[uuid.UUID]*model.User),
	}
}

func (m *mockRepo) Create(_ context.Context, user *model.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockRepo) GetAll(_ context.Context) ([]*model.User, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	var result []*model.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *mockRepo) Update(_ context.Context, user *model.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepo) Delete(_ context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepo) ExistsByEmail(_ context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	if m.existsByEmail != nil {
		return m.existsByEmail(email, excludeID)
	}
	for id, u := range m.users {
		if u.Email == email {
			if excludeID != nil && id == *excludeID {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

func TestCreateUser_Success(t *testing.T) {
	svc := NewUserService(newMockRepo())

	user, err := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "John Doe" {
		t.Fatalf("expected name John Doe, got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %s", user.Email)
	}
	if user.ID == uuid.Nil {
		t.Fatal("expected non-nil UUID")
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	svc := NewUserService(newMockRepo())

	_, err := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "",
		Email: "john@example.com",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateUser_EmptyEmail(t *testing.T) {
	svc := NewUserService(newMockRepo())

	_, err := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John Doe",
		Email: "",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "dup@example.com",
	})

	_, err := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "Jane",
		Email: "dup@example.com",
	})
	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}

func TestCreateUser_RepoError(t *testing.T) {
	mock := newMockRepo()
	mock.createErr = errors.New("db error")
	svc := NewUserService(mock)

	_, err := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected 'db error', got %v", err)
	}
}

func TestGetByID_Success(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	created, _ := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})

	user, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != created.ID {
		t.Fatalf("expected ID %v, got %v", created.ID, user.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := NewUserService(newMockRepo())

	_, err := svc.GetByID(context.Background(), uuid.New())
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetAll_Success(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	svc.Create(context.Background(), &model.CreateUserRequest{Name: "A", Email: "a@test.com"})
	svc.Create(context.Background(), &model.CreateUserRequest{Name: "B", Email: "b@test.com"})

	users, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestGetAll_Empty(t *testing.T) {
	svc := NewUserService(newMockRepo())

	users, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	created, _ := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})

	updated, err := svc.Update(context.Background(), created.ID, &model.UpdateUserRequest{
		Name: "John Updated",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "John Updated" {
		t.Fatalf("expected 'John Updated', got %s", updated.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	svc := NewUserService(newMockRepo())

	_, err := svc.Update(context.Background(), uuid.New(), &model.UpdateUserRequest{
		Name: "Nobody",
	})
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUser_DuplicateEmail(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})
	user2, _ := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "Jane",
		Email: "jane@example.com",
	})

	_, err := svc.Update(context.Background(), user2.ID, &model.UpdateUserRequest{
		Email: "john@example.com",
	})
	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}

func TestUpdateUser_SameEmail(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	created, _ := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})

	updated, err := svc.Update(context.Background(), created.ID, &model.UpdateUserRequest{
		Email: "john@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Email != "john@example.com" {
		t.Fatalf("expected same email, got %s", updated.Email)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mock := newMockRepo()
	svc := NewUserService(mock)

	created, _ := svc.Create(context.Background(), &model.CreateUserRequest{
		Name:  "John",
		Email: "john@example.com",
	})

	err := svc.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc := NewUserService(newMockRepo())

	err := svc.Delete(context.Background(), uuid.New())
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
