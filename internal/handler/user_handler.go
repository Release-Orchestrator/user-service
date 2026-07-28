package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Release-Orchestrator/user-service/internal/model"
	"github.com/Release-Orchestrator/user-service/internal/service"
)

type UserHandler struct {
	svc service.UserServiceInterface
}

func NewUserHandler(svc service.UserServiceInterface) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.Get)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	user, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": user})
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_ID", "message": "invalid user ID"}})
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_ID", "message": "invalid user ID"}})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()}})
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_ID", "message": "invalid user ID"}})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "USER_NOT_FOUND", "message": "user not found"}})
	case service.ErrEmailExists:
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": gin.H{"code": "EMAIL_ALREADY_EXISTS", "message": "email already exists"}})
	case service.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "INVALID_INPUT", "message": "invalid input"}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
	}
}