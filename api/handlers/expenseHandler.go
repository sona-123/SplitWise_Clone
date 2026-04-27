package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sona-123/splitwise_clone/models"
)

// ExpenseHandler adds a new expense to a specific group
// @Summary Create a new expense
// @Description Add an expense to a group (requires authentication)
// @Tags Expenses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Expense true "Expense Data"
// @Success 201 {object} map[string]string "Expense added"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /expenses [post]
func (h *Handler) ExpenseHandler(c *gin.Context) {
	val, exists := c.Get("current_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context missing"})
		return
	}
	currentUserID := val.(int)
	var exp models.Expense

	// ShouldBindJSON maps the "group_id", "paid_by", etc., from JSON to the struct
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic Validation: Ensure group and amount are provided
	if exp.GroupID == 0 || exp.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id and a positive amount are required"})
		return
	}

	if len(exp.UserIds) == 0 && len(exp.Shares) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No participants provided"})
		return
	}
	exp.PaidBy = currentUserID
	if err := h.Service.CreateExpense(exp); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense added successfully"})
}
