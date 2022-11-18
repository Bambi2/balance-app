package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get Transactions
// @Description Get user's transactions ordered by date and amount (pagination included)
// @Tags transactions
// @ID get-transactions
// @Accept */*
// @Produce json
// @Success 200 {object} []domain.Transaction
// @Failure 400,404 {object} errorResponse
// @Router /api/reserve/{userId}/{limit}/{offset} [get]
func (h *Handler) GetTransactions(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid 'limit' parameter")
		return
	}

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid 'limit' parameter")
		return
	}

	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid 'offset' parameter")
		return
	}

	transactions, err := h.service.Transaction.GetTransactions(userId, limit, offset)
	//might add more types of errors and delete userId field from transaction struct(domain)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}
