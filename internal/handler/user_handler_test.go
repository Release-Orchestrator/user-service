package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Release-Orchestrator/user-service/internal/model"
	"github.com/Release-Orchestrator/user-service/internal/service"
)

type mockService struct {
	createFunc  func(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.User, error)
	getAllFunc  func(ctx context.Context) ([]*model.User, error)
	updateFunc  func(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error)
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	return m.createFunc(ctx, req)
}

func (m *mockService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockService) GetAll(ctx context.Context) ([]*model.User, error) {
	return m.getAllFunc(ctx)
}

func (m *mockService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error) {
	return m.updateFunc(ctx, id, req)
}

func (m *mockService) Delete(ctx context.Context, id uuid.UUID) error {
	return m.deleteFunc(ctx, id)
}

func setupRouter(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

func now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func TestHealthEndpoint(t *testing.T) {
	h := NewUserHandler(&mockService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "UP" {
		t.Fatalf("expected status UP, got %s", resp["status"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	svc := &mockService{
		createFunc: func(_ context.Context, req *model.CreateUserRequest) (*model.User, error) {
			return &model.User{
				ID:        uuid.New(),
				Name:      req.Name,
				Email:     req.Email,
				CreatedAt: now(),
				UpdatedAt: now(),
			}, nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]string{"name": "John", "email": "john@test.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success true")
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	svc := &mockService{}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewReader([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_EmailExists(t *testing.T) {
	svc := &mockService{
		createFunc: func(_ context.Context, _ *model.CreateUserRequest) (*model.User, error) {
			return nil, service.ErrEmailExists
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]string{"name": "John", "email": "dup@test.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestGetUser_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockService{
		getByIDFunc: func(_ context.Context, id uuid.UUID) (*model.User, error) {
			return &model.User{
				ID:    id,
				Name:  "John",
				Email: "john@test.com",
			}, nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/"+userID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc := &mockService{
		getByIDFunc: func(_ context.Context, _ uuid.UUID) (*model.User, error) {
			return nil, service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	h := NewUserHandler(&mockService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListUsers_Success(t *testing.T) {
	svc := &mockService{
		getAllFunc: func(_ context.Context) ([]*model.User, error) {
			return []*model.User{
				{ID: uuid.New(), Name: "John", Email: "john@test.com"},
				{ID: uuid.New(), Name: "Jane", Email: "jane@test.com"},
			}, nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Fatalf("expected 2 users, got %d", len(data))
	}
}

func TestListUsers_Empty(t *testing.T) {
	svc := &mockService{
		getAllFunc: func(_ context.Context) ([]*model.User, error) {
			return []*model.User{}, nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockService{
		updateFunc: func(_ context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error) {
			return &model.User{
				ID:    id,
				Name:  req.Name,
				Email: "john@test.com",
			}, nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]string{"name": "Updated"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	svc := &mockService{
		updateFunc: func(_ context.Context, _ uuid.UUID, _ *model.UpdateUserRequest) (*model.User, error) {
			return nil, service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]string{"name": "Nobody"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+uuid.New().String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	h := NewUserHandler(&mockService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+uuid.New().String(), bytes.NewReader([]byte(`{bad}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	svc := &mockService{
		deleteFunc: func(_ context.Context, _ uuid.UUID) error {
			return nil
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc := &mockService{
		deleteFunc: func(_ context.Context, _ uuid.UUID) error {
			return service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteUser_InvalidID(t *testing.T) {
	h := NewUserHandler(&mockService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/bad-id", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleError_InternalError(t *testing.T) {
	svc := &mockService{
		getByIDFunc: func(_ context.Context, _ uuid.UUID) (*model.User, error) {
			return nil, errors.New("unexpected db failure")
		},
	}
	h := NewUserHandler(svc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
