package delivery

import (
	"net/http"
	"os"

	"github.com/bambi2/balance-app/internal/domain"
	"github.com/bambi2/balance-app/internal/service"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

// @Summary Check Invoice
// @Description Check an invoice
// @Tags check
// @ID check
// @Accept json
// @Produce json
// @Param check body domain.Check true "invoice info and amount (in cents) to check"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/check [post]
func (h *Handler) Check(c *gin.Context) {
	var check domain.Check

	if err := c.BindJSON(&check); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Check.Check(check); err != nil {
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

// @Summary Get Checks for A Month As CSV File
// @Description Returns url (using cloudinary) to a csv file with checks for the month specified in request body
// @Tags check
// @ID get-checks
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/checks/{year}/{month} [get]
func (h *Handler) GetChecks(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	csvFile, err := h.service.Check.GetCSVChecks(year, month)
	if err != nil {
		if ce, ok := err.(*service.CustomError); ok {
			newErrorResponse(c, ce.StatusCode, ce.Message)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	defer os.Remove(csvFile.Name())

	cld, err := cloudinary.New()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := cld.Upload.Upload(c, csvFile.Name(), uploader.UploadParams{
		PublicID: "checks",
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"url": result.URL,
	})
}
