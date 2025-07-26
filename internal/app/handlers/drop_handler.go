package handlers

import (
	"RandomItems/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DropHandler struct {
	dropService *services.DropService
}

func NewDropHandler(dropService *services.DropService) *DropHandler {
	return &DropHandler{dropService: dropService}
}

func (h *DropHandler) GenerateDrop(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	item, err := h.dropService.GenerateDrop(c.Request.Context(), userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No drop this time, pity counter increased"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}
