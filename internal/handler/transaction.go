package handler

import (
	"errors"
	"github.com/aberniando/xenshield/internal/enum"
	"github.com/aberniando/xenshield/pkg/logger"
	"github.com/go-playground/validator/v10"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aberniando/xenshield/internal/entity"
	"github.com/aberniando/xenshield/internal/usecases/transaction"
)

type TransactionHandler struct {
	svc    transaction.Service
	logger *logger.Logger
}

func NewTransactionHandler(svc transaction.Service, logger *logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		svc:    svc,
		logger: logger,
	}
}

func (r *TransactionHandler) Insert(c *gin.Context) {
	var request entity.InsertTransactionRequest
	var errorItems []ValidationErrorItem

	if err := c.ShouldBindJSON(&request); err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fe := range validationErrors {
				errorItems = append(errorItems, ValidationErrorItem{
					Field:   fe.Field(),
					Message: getValidationMessage(fe),
				})
			}
			returnValidationErrorResponse(c, errorItems)
			return
		}

		returnErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	errorItems = validateInsertTransactionRequest(&request)

	if len(errorItems) > 0 {
		returnValidationErrorResponse(c, errorItems)
		return
	}

	resp, err := r.svc.InsertTransaction(c.Request.Context(), &request)
	if err != nil {
		r.logger.Error(err, "error inserting new transaction")
		returnErrorResponse(c, http.StatusInternalServerError, "Error inserting new transaction")

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *TransactionHandler) GetIPAddressStolenCardHistory(c *gin.Context) {
	ipAddress := c.Param("ip_address")
	var err error

	if ipAddress == "" {
		err = errors.New("empty ip address path param")
		r.logger.Error(err)
		returnErrorResponse(c, http.StatusBadRequest, "ip_address on path param is required")
		return
	}

	if net.ParseIP(ipAddress) == nil {
		err = errors.New("invalid ip address path param")
		r.logger.Error(err)
		returnErrorResponse(c, http.StatusBadRequest, "invalid ip_address value")
		return
	}

	resp, err := r.svc.GetIPAddressStolenCardHistory(c.Request.Context(), ipAddress)
	if err != nil {
		r.logger.Error(err, "error getting ip address stolen card history")
		returnErrorResponse(c, http.StatusInternalServerError, "error getting ip address stolen card history")

		return
	}

	if !resp.HasTransaction {
		returnErrorResponse(c, http.StatusNotFound, "given ip address has not made any transaction yet")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	default:
		return fe.Error()
	}
}

func validateInsertTransactionRequest(req *entity.InsertTransactionRequest) []ValidationErrorItem {
	errorItems := make([]ValidationErrorItem, 0)
	if net.ParseIP(req.IPAddress) == nil {
		errorItems = append(errorItems, ValidationErrorItem{
			Field:   "IPAddress",
			Message: "invalid value format",
		})
	}

	if len(req.MaskedCardNumber) != 19 {
		errorItems = append(errorItems, ValidationErrorItem{
			Field:   "MaskedCardNumber",
			Message: "invalid length",
		})
	}

	if req.Status == enum.TransactionStatusSuccess && req.FailureReason != "" {
		errorItems = append(errorItems, ValidationErrorItem{
			Field:   "FailureReason",
			Message: "failure reason must be null on success transaction",
		})
	}

	if req.Status == enum.TransactionStatusFailed && req.FailureReason == "" {
		errorItems = append(errorItems, ValidationErrorItem{
			Field:   "FailureReason",
			Message: "failure reason must not be null on failed transaction",
		})
	}

	return errorItems
}
