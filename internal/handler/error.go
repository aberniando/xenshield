package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Message []ValidationErrorItem `json:"message"`
}

func returnErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, ErrorResponse{Message: msg})
}

func returnValidationErrorResponse(c *gin.Context, items []ValidationErrorItem) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ValidationErrorResponse{Message: items})
}
