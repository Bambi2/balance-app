package delivery

import (
	"net/http"
	"strconv"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get Balance
// @Description Get balance of a specified user (in cents)
// @Tags balance
// @ID get-balance
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/balance/{userId} [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	amount, err := h.service.User.GetBalance(userId)
	if err != nil {
		if ce, ok := err.(*service.CustomError); ok {
			newErrorResponse(c, ce.StatusCode, ce.Message)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"amount": amount,
	})
}

// @Summary Add To Balance
// @Description Add money to the balance of a user or create a new balance if he/she doesn't exist yet
// @Tags balance
// @ID add-to-balance
// @Accept json
// @Produce json
// @Param user body domain.User true "money to add (in cents) and user id"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/balance [post]
func (h *Handler) AddToBalance(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	amount, err := h.service.User.AddToBalance(user)
	if err != nil {
		if ce, ok := err.(*service.CustomError); ok {
			newErrorResponse(c, ce.StatusCode, ce.Message)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"amount": amount,
	})
}

type sendBalanceInput struct {
	SenderId   int         `json:"senderId" binding:"required"`
	RecieverId int         `json:"recieverId" binding:"required"`
	Amount     domain.CENT `json:"amount" binding:"required"`
}

// @Summary Send Money From Another User
// @Description Make a transaction between users
// @Tags balance
// @ID send-money
// @Accept json
// @Produce json
// @Param input body sendBalanceInput true "information about the transaction between users"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/balance/send [post]
func (h *Handler) SendMoney(c *gin.Context) {
	var input sendBalanceInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.service.User.SendMoney(input.SenderId, input.RecieverId, input.Amount); err != nil {
		if ce, ok := err.(*service.CustomError); ok {
			newErrorResponse(c, ce.StatusCode, ce.Message)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
