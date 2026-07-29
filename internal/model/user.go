// Package model contains domain models used by the user-service.
package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest is the payload to create a user.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest is the payload to update a user.
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}
