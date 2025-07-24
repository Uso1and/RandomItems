package handlers

import (
	"RandomItems/internal/domain/models"
	"RandomItems/internal/domain/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserHandler(userRepo repositories.UserRepositoryInterface) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	if err := h.userRepo.CreateUserRep(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
