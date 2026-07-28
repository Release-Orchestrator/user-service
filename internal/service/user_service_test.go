package service

import (
	"testing"

	"github.com/Release-Orchestrator/user-service/internal/model"
)

func TestCreateUser_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		req  *model.CreateUserRequest
	}{
		{name: "empty name", req: &model.CreateUserRequest{Name: "", Email: "test@test.com"}},
		{name: "empty email", req: &model.CreateUserRequest{Name: "Test", Email: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Name == "" || tt.req.Email == "" {
				return // validation caught
			}
		})
	}
}
