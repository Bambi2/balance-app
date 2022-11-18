package delivery

import (
	"net/http"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary Reserve Money
// @Description Reserve user's money for a service
// @Tags reservation
// @ID reserve
// @Accept json
// @Produce json
// @Param invoice body domain.Invoice true "user id and invoice info"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/reserve [post]
func (h *Handler) Reserve(c *gin.Context) {
	var invoice domain.Invoice

	if err := c.BindJSON(&invoice); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Invoice.Reserve(invoice); err != nil {
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

type unreserveInput struct {
	UserId    int   `json:"userId" binding:"required"`
	ServiceId int   `json:"serviceId" binding:"required"`
	OrderId   int64 `json:"orderId" binding:"required"`
}

// @Summary Unreserve Money
// @Description Unreserve user's money for a service if the service isn't aplied yet
// @Tags reservation
// @ID unreserve
// @Accept json
// @Produce json
// @Param input body unreserveInput true "user id and invoice info"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/unreserve [post]
func (h *Handler) Unreserve(c *gin.Context) {
	var input unreserveInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Invoice.Unreserve(input.UserId, input.ServiceId, input.OrderId); err != nil {
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
